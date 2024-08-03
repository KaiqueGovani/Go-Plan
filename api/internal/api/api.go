package api

import (
	"context"
	"encoding/json"
	"errors"
	"journey/internal/api/spec"
	"journey/internal/pgstore"
	"net/http"
	"slices"
	"time"

	"github.com/discord-gophers/goapi-gen/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type store interface {
	CreateTrip(ctx context.Context, pool *pgxpool.Pool, params spec.CreateTripRequest) (uuid.UUID, error)
	GetTrip(ctx context.Context, id uuid.UUID) (pgstore.Trip, error)
	GetAllTrips(ctx context.Context) ([]pgstore.Trip, error)
	UpdateTrip(ctx context.Context, arg pgstore.UpdateTripParams) error
	GetParticipant(ctx context.Context, participantID uuid.UUID) (pgstore.Participant, error)
	GetParticipants(ctx context.Context, tripID uuid.UUID) ([]pgstore.Participant, error)
	ConfirmParticipant(ctx context.Context, participantID uuid.UUID) error
	GetTripActivities(ctx context.Context, tripID uuid.UUID) ([]pgstore.Activity, error)
	CreateActivity(ctx context.Context, arg pgstore.CreateActivityParams) (uuid.UUID, error)
}

type mailer interface {
	SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error
	SendConfirmTripEmailToTripParticipants(tripID uuid.UUID) error
}

type API struct{
	store store
	logger *zap.Logger
	validator *validator.Validate
	pool *pgxpool.Pool
	mailer mailer
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger, mailer mailer) API {
	validator := validator.New(validator.WithRequiredStructEnabled())

	return API{pgstore.New(pool), logger, validator, pool, mailer}
}

// Confirms a participant on a trip.
// (PATCH /participants/{participantId}/confirm)
func (api API) PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request, participantID string) *spec.Response {
	id, err := uuid.Parse(participantID)
	if err != nil {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "Invalid participant ID"})
	}

	particiapant, err := api.store.GetParticipant(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "Participant not found"})
		}
		api.logger.Error("Failed to get participant", zap.Error(err), zap.String("participant_id", participantID))
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "Something went wrong, try again"}) 
	}

	if particiapant.IsConfirmed {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "Participant already confirmed"})
	}

	if err := api.store.ConfirmParticipant(r.Context(), id); err != nil {
		api.logger.Error("Failed to confirm participant", zap.Error(err), zap.String("participant_id", participantID))
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "Something went wrong, try again"}) 
	}

	return spec.PatchParticipantsParticipantIDConfirmJSON204Response(nil)
}

// Create a new trip
// (POST /trips)
func (api API) PostTrips(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.CreateTripRequest;
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "Invalid JSON: " + err.Error()}) 
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "Invalid request body: "+ err.Error()})
	}

	tripID, err := api.store.CreateTrip(r.Context(), api.pool, body)
	if err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	go func() {
		if err := api.mailer.SendConfirmTripEmailToTripOwner(tripID); err != nil {
			api.logger.Error("Failed to send email on PostTrips", zap.Error(err), zap.String("trip_id", tripID.String()), zap.String("owner_email", string(body.OwnerEmail)))
		}
	}()

	return spec.PostTripsJSON201Response(spec.CreateTripResponse{TripID: tripID.String()})
}

// Get all trips.
// (GET /trips)
func (api API) GetTrips(w http.ResponseWriter, r *http.Request) *spec.Response {
	trips, err := api.store.GetAllTrips(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return spec.GetTripsJSON400Response(spec.Error{Message: "No trips found"})	
		} 

		api.logger.Error("Failed to get trips", zap.Error(err))
		return spec.GetTripsJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	tripsResponse := make([]spec.GetTripDetailsResponseTripObj, len(trips))
	for i, trip := range trips {
		tripsResponse[i] = spec.GetTripDetailsResponseTripObj{
			ID: trip.ID.String(),
			Destination: trip.Destination,
			EndsAt: trip.EndsAt.Time,
			StartsAt: trip.StartsAt.Time,
			IsConfirmed: trip.IsConfirmed,
		}
	}

	return spec.GetTripsJSON200Response(spec.GetTripsResponse{
		Trips: tripsResponse,
	})
}

// Get a trip details.
// (GET /trips/{tripId})
func (api API) GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDJSON400Response(spec.Error{Message: "Invalid trip ID"})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return spec.GetTripsTripIDJSON400Response(spec.Error{Message: "Trip not found"})	
		} 
		api.logger.Error("Failed to get trip", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	return spec.GetTripsTripIDJSON200Response(spec.GetTripDetailsResponse{
		Trip: spec.GetTripDetailsResponseTripObj{
			ID: trip.ID.String(),
			Destination: trip.Destination,
			EndsAt: trip.EndsAt.Time,
			StartsAt: trip.StartsAt.Time,
			IsConfirmed: trip.IsConfirmed,
		},
	})
}

// Update a trip.
// (PUT /trips/{tripId})
func (api API) PutTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "Invalid trip ID"})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "Trip not found"})	
		} 
		api.logger.Error("Failed to get trip", zap.Error(err), zap.String("trip_id", tripID))
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	var body spec.PutTripsTripIDJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "Invalid JSON: " + err.Error()})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "Invalid request body: " + err.Error()})
	}

	if err := api.store.UpdateTrip(r.Context(), pgstore.UpdateTripParams{
		ID: id,
		Destination: body.Destination,
		EndsAt: pgtype.Timestamp{Valid: true, Time: body.EndsAt},
		StartsAt: pgtype.Timestamp{Valid: true, Time: body.StartsAt},
		IsConfirmed: trip.IsConfirmed,
	}); err != nil {
		api.logger.Error("Failed to update trip", zap.Error(err), zap.String("trip_id", tripID))
		return spec.PutTripsTripIDJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	return spec.PutTripsTripIDJSON204Response(nil)
}

// Get a trip activities.
// (GET /trips/{tripId}/activities)
func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Invalid trip ID"})
	}

	activities, err := api.store.GetTripActivities(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Activities not found"})
		}

		api.logger.Error("Failed to get activities", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	type Activity struct {
		Time time.Time
		Amount *int
	}


	var qtd int
	differentDates := make([]Activity, 0, len(activities))
	for _, activity := range activities {
		// Check if the date isnt already in the slice
		if !slices.ContainsFunc(
			differentDates, 
			func(item Activity) bool { 
				itemYear, itemMonth, itemDay := item.Time.Date()
				activityYear, activityMonth, activityDay := activity.OccursAt.Time.Date()
				return itemYear == activityYear && itemMonth == activityMonth && itemDay == activityDay
			}) {
			var amount int = 1;
			itemYear, itemMonth, itemDay := activity.OccursAt.Time.Date()
			differentDates = append(differentDates, Activity{
				Time: time.Date(itemYear, itemMonth, itemDay, 0, 0, 0, 0, time.UTC),
				Amount: &amount,
			})	
			qtd++
		} else {
			for i, date := range differentDates {
				dateYear, dateMonth, dateDay := date.Time.Date()
				activityYear, activityMonth, activityDay := activity.OccursAt.Time.Date()
				if dateYear == activityYear && dateMonth == activityMonth && dateDay == activityDay {
        			*date.Amount++
					differentDates[i] = date
				}
			}
		}
	}

	activitiesResponse := make([]spec.GetTripActivitiesResponseOuterArray, qtd)
	for i, item := range differentDates {
		activitiesInnerResponse := make([]spec.GetTripActivitiesResponseInnerArray, 0, *differentDates[i].Amount)

		for _, activity := range activities {
			dateYear, dateMonth, dateDay := item.Time.Date()
			activityYear, activityMonth, activityDay := activity.OccursAt.Time.Date()
			if dateYear == activityYear && dateMonth == activityMonth && dateDay == activityDay {
        			activitiesInnerResponse = append(activitiesInnerResponse, spec.GetTripActivitiesResponseInnerArray{
					ID: activity.ID.String(),
					Title: activity.Title,
					OccursAt: activity.OccursAt.Time,
				})
			}
		}
		
		activitiesResponse[i] = spec.GetTripActivitiesResponseOuterArray{
			Date: item.Time,
			Activities: activitiesInnerResponse,
		}
	}

	
	return spec.GetTripsTripIDActivitiesJSON200Response(spec.GetTripActivitiesResponse{
		Activities: activitiesResponse,
	})
}

// Create a trip activity.
// (POST /trips/{tripId}/activities)
func (api API) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Invalid trip ID"})
	}

	var body spec.CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Invalid JSON: " + err.Error()})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Invalid request body: " + err.Error()})
	}

	activityID, err := api.store.CreateActivity(r.Context(), pgstore.CreateActivityParams{
		TripID: id,
		Title: body.Title,
		OccursAt: pgtype.Timestamp{Valid: true, Time: body.OccursAt},
	})
	if err != nil {
		api.logger.Error("Failed to create activity", zap.Error(err), zap.String("trip_id", tripID))
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	return spec.PostTripsTripIDActivitiesJSON201Response(spec.CreateActivityResponse{ActivityID: activityID.String()})
}

// Confirm a trip and send e-mail invitations.
// (GET /trips/{tripId}/confirm)
func (api API) GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: "Invalid trip ID"})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: "Trip not found"})	
		} 
		api.logger.Error("Failed to get trip", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	if trip.IsConfirmed {
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: "Trip already confirmed"})
	}

	// Update trip to confirm
	if err := api.store.UpdateTrip(r.Context(), pgstore.UpdateTripParams{
		ID: id,
		Destination: trip.Destination,
		EndsAt: trip.EndsAt,
		StartsAt: trip.StartsAt,
		IsConfirmed: true,
	}); err != nil {
		api.logger.Error("Failed to confirm trip", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	// Send e-mail invitations to participants
	go func() {
		if err := api.mailer.SendConfirmTripEmailToTripParticipants(id); err != nil {
			api.logger.Error("Failed to send email on GetTripsTripIDConfirm", zap.Error(err), zap.String("trip_id", tripID))
		}
	}()

	return spec.GetTripsTripIDConfirmJSON204Response(nil)
}

// Invite someone to the trip.
// (POST /trips/{tripId}/invites)
func (api API) PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}

// Get a trip links.
// (GET /trips/{tripId}/links)
func (api API) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}

// Create a trip link.
// (POST /trips/{tripId}/links)
func (api API) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") // TODO: Implement
}

// Get a trip participants.
// (GET /trips/{tripId}/participants)
func (api API) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: "Invalid trip ID"})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: "Trip not found"})	
		} 
		api.logger.Error("Failed to get trip", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	participants, err := api.store.GetParticipants(r.Context(), trip.ID)
	if err != nil {
		api.logger.Error("Failed to get participants", zap.Error(err), zap.String("trip_id", tripID))
		return spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: "Something went wrong, try again"})
	}

	participantsResponse := make([]spec.GetTripParticipantsResponseArray , len(participants))
	for i, participant := range participants {

		participantsResponse[i] = spec.GetTripParticipantsResponseArray {
			ID: participant.ID.String(),
			Email: types.Email(participant.Email),
			IsConfirmed: participant.IsConfirmed,
		}
	}

	return spec.GetTripsTripIDParticipantsJSON200Response(spec.GetTripParticipantsResponse{
		Participants: participantsResponse,
	})
}
// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.3.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/discord-gophers/goapi-gen/runtime"
	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// CreateActivityRequest defines model for CreateActivityRequest.
type CreateActivityRequest struct {
	OccursAt time.Time `json:"occurs_at" validate:"required"`
	Title    string    `json:"title" validate:"required"`
}

// CreateActivityResponse defines model for CreateActivityResponse.
type CreateActivityResponse struct {
	ActivityID string `json:"activityId"`
}

// CreateLinkRequest defines model for CreateLinkRequest.
type CreateLinkRequest struct {
	Title string `json:"title" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}

// CreateLinkResponse defines model for CreateLinkResponse.
type CreateLinkResponse struct {
	LinkID string `json:"linkId"`
}

// CreateTripRequest defines model for CreateTripRequest.
type CreateTripRequest struct {
	Destination    string                `json:"destination" validate:"required,min=4"`
	EmailsToInvite []openapi_types.Email `json:"emails_to_invite" validate:"required,dive,email"`
	EndsAt         time.Time             `json:"ends_at" validate:"required"`
	OwnerEmail     openapi_types.Email   `json:"owner_email" validate:"required,email"`
	OwnerName      string                `json:"owner_name" validate:"required"`
	StartsAt       time.Time             `json:"starts_at" validate:"required"`
}

// CreateTripResponse defines model for CreateTripResponse.
type CreateTripResponse struct {
	TripID string `json:"tripId"`
}

// Bad request
type Error struct {
	Message string `json:"message"`
}

// GetLinksResponse defines model for GetLinksResponse.
type GetLinksResponse struct {
	Links []GetLinksResponseArray `json:"links"`
}

// GetLinksResponseArray defines model for GetLinksResponseArray.
type GetLinksResponseArray struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GetTripActivitiesResponse defines model for GetTripActivitiesResponse.
type GetTripActivitiesResponse struct {
	Activities []GetTripActivitiesResponseOuterArray `json:"activities"`
}

// GetTripActivitiesResponseInnerArray defines model for GetTripActivitiesResponseInnerArray.
type GetTripActivitiesResponseInnerArray struct {
	ID       string    `json:"id"`
	OccursAt time.Time `json:"occurs_at"`
	Title    string    `json:"title"`
}

// GetTripActivitiesResponseOuterArray defines model for GetTripActivitiesResponseOuterArray.
type GetTripActivitiesResponseOuterArray struct {
	Activities []GetTripActivitiesResponseInnerArray `json:"activities"`
	Date       time.Time                             `json:"date"`
}

// GetTripDetailsResponse defines model for GetTripDetailsResponse.
type GetTripDetailsResponse struct {
	Trip GetTripDetailsResponseTripObj `json:"trip"`
}

// GetTripDetailsResponseTripObj defines model for GetTripDetailsResponseTripObj.
type GetTripDetailsResponseTripObj struct {
	Destination string    `json:"destination"`
	EndsAt      time.Time `json:"ends_at"`
	ID          string    `json:"id"`
	IsConfirmed bool      `json:"is_confirmed"`
	StartsAt    time.Time `json:"starts_at"`
}

// GetTripParticipantsResponse defines model for GetTripParticipantsResponse.
type GetTripParticipantsResponse struct {
	Participants []GetTripParticipantsResponseArray `json:"participants"`
}

// GetTripParticipantsResponseArray defines model for GetTripParticipantsResponseArray.
type GetTripParticipantsResponseArray struct {
	Email       openapi_types.Email `json:"email"`
	ID          string              `json:"id"`
	IsConfirmed bool                `json:"is_confirmed"`
	Name        *string             `json:"name"`
}

// GetTripsResponse defines model for GetTripsResponse.
type GetTripsResponse struct {
	Trips []GetTripDetailsResponseTripObj `json:"trips"`
}

// InviteParticipantRequest defines model for InviteParticipantRequest.
type InviteParticipantRequest struct {
	Email openapi_types.Email `json:"email" validate:"required,email"`
}

// UpdateTripRequest defines model for UpdateTripRequest.
type UpdateTripRequest struct {
	Destination string    `json:"destination" validate:"required,min=4"`
	EndsAt      time.Time `json:"ends_at" validate:"required"`
	StartsAt    time.Time `json:"starts_at" validate:"required"`
}

// PostTripsJSONBody defines parameters for PostTrips.
type PostTripsJSONBody CreateTripRequest

// PutTripsTripIDJSONBody defines parameters for PutTripsTripID.
type PutTripsTripIDJSONBody UpdateTripRequest

// PostTripsTripIDActivitiesJSONBody defines parameters for PostTripsTripIDActivities.
type PostTripsTripIDActivitiesJSONBody CreateActivityRequest

// PostTripsTripIDInvitesJSONBody defines parameters for PostTripsTripIDInvites.
type PostTripsTripIDInvitesJSONBody InviteParticipantRequest

// PostTripsTripIDLinksJSONBody defines parameters for PostTripsTripIDLinks.
type PostTripsTripIDLinksJSONBody CreateLinkRequest

// PostTripsJSONRequestBody defines body for PostTrips for application/json ContentType.
type PostTripsJSONRequestBody PostTripsJSONBody

// Bind implements render.Binder.
func (PostTripsJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PutTripsTripIDJSONRequestBody defines body for PutTripsTripID for application/json ContentType.
type PutTripsTripIDJSONRequestBody PutTripsTripIDJSONBody

// Bind implements render.Binder.
func (PutTripsTripIDJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostTripsTripIDActivitiesJSONRequestBody defines body for PostTripsTripIDActivities for application/json ContentType.
type PostTripsTripIDActivitiesJSONRequestBody PostTripsTripIDActivitiesJSONBody

// Bind implements render.Binder.
func (PostTripsTripIDActivitiesJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostTripsTripIDInvitesJSONRequestBody defines body for PostTripsTripIDInvites for application/json ContentType.
type PostTripsTripIDInvitesJSONRequestBody PostTripsTripIDInvitesJSONBody

// Bind implements render.Binder.
func (PostTripsTripIDInvitesJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostTripsTripIDLinksJSONRequestBody defines body for PostTripsTripIDLinks for application/json ContentType.
type PostTripsTripIDLinksJSONRequestBody PostTripsTripIDLinksJSONBody

// Bind implements render.Binder.
func (PostTripsTripIDLinksJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// PatchParticipantsParticipantIDConfirmJSON204Response is a constructor method for a PatchParticipantsParticipantIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func PatchParticipantsParticipantIDConfirmJSON204Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        204,
		contentType: "application/json",
	}
}

// PatchParticipantsParticipantIDConfirmJSON400Response is a constructor method for a PatchParticipantsParticipantIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func PatchParticipantsParticipantIDConfirmJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsJSON200Response is a constructor method for a GetTrips response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsJSON200Response(body GetTripsResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsJSON400Response is a constructor method for a GetTrips response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsJSON201Response is a constructor method for a PostTrips response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsJSON201Response(body CreateTripResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTripsJSON400Response is a constructor method for a PostTrips response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDJSON200Response is a constructor method for a GetTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDJSON200Response(body GetTripDetailsResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDJSON400Response is a constructor method for a GetTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PutTripsTripIDJSON204Response is a constructor method for a PutTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func PutTripsTripIDJSON204Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        204,
		contentType: "application/json",
	}
}

// PutTripsTripIDJSON400Response is a constructor method for a PutTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func PutTripsTripIDJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDActivitiesJSON200Response is a constructor method for a GetTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDActivitiesJSON200Response(body GetTripActivitiesResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDActivitiesJSON400Response is a constructor method for a GetTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDActivitiesJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsTripIDActivitiesJSON201Response is a constructor method for a PostTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDActivitiesJSON201Response(body CreateActivityResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTripsTripIDActivitiesJSON400Response is a constructor method for a PostTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDActivitiesJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDConfirmJSON204Response is a constructor method for a GetTripsTripIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDConfirmJSON204Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        204,
		contentType: "application/json",
	}
}

// GetTripsTripIDConfirmJSON400Response is a constructor method for a GetTripsTripIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDConfirmJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsTripIDInvitesJSON201Response is a constructor method for a PostTripsTripIDInvites response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDInvitesJSON201Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTripsTripIDInvitesJSON400Response is a constructor method for a PostTripsTripIDInvites response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDInvitesJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDLinksJSON200Response is a constructor method for a GetTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDLinksJSON200Response(body GetLinksResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDLinksJSON400Response is a constructor method for a GetTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDLinksJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsTripIDLinksJSON201Response is a constructor method for a PostTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDLinksJSON201Response(body CreateLinkResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTripsTripIDLinksJSON400Response is a constructor method for a PostTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDLinksJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDParticipantsJSON200Response is a constructor method for a GetTripsTripIDParticipants response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDParticipantsJSON200Response(body GetTripParticipantsResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDParticipantsJSON400Response is a constructor method for a GetTripsTripIDParticipants response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDParticipantsJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Confirms a participant on a trip.
	// (PATCH /participants/{participantId}/confirm)
	PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request, participantID string) *Response
	// Lists all trips
	// (GET /trips)
	GetTrips(w http.ResponseWriter, r *http.Request) *Response
	// Create a new trip
	// (POST /trips)
	PostTrips(w http.ResponseWriter, r *http.Request) *Response
	// Get a trip details.
	// (GET /trips/{tripId})
	GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Update a trip.
	// (PUT /trips/{tripId})
	PutTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Get a trip activities.
	// (GET /trips/{tripId}/activities)
	GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Create a trip activity.
	// (POST /trips/{tripId}/activities)
	PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Confirm a trip and send e-mail invitations.
	// (GET /trips/{tripId}/confirm)
	GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Invite someone to the trip.
	// (POST /trips/{tripId}/invites)
	PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Get a trip links.
	// (GET /trips/{tripId}/links)
	GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Create a trip link.
	// (POST /trips/{tripId}/links)
	PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Get a trip participants.
	// (GET /trips/{tripId}/participants)
	GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// PatchParticipantsParticipantIDConfirm operation middleware
func (siw *ServerInterfaceWrapper) PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "participantId" -------------
	var participantID string

	if err := runtime.BindStyledParameter("simple", false, "participantId", chi.URLParam(r, "participantId"), &participantID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "participantId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PatchParticipantsParticipantIDConfirm(w, r, participantID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTrips operation middleware
func (siw *ServerInterfaceWrapper) GetTrips(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTrips(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTrips operation middleware
func (siw *ServerInterfaceWrapper) PostTrips(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTrips(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripID operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripID(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PutTripsTripID operation middleware
func (siw *ServerInterfaceWrapper) PutTripsTripID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PutTripsTripID(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDActivities operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDActivities(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTripsTripIDActivities operation middleware
func (siw *ServerInterfaceWrapper) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTripsTripIDActivities(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDConfirm operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDConfirm(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTripsTripIDInvites operation middleware
func (siw *ServerInterfaceWrapper) PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTripsTripIDInvites(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDLinks operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDLinks(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTripsTripIDLinks operation middleware
func (siw *ServerInterfaceWrapper) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTripsTripIDLinks(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDParticipants operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDParticipants(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:    "/",
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Patch("/participants/{participantId}/confirm", wrapper.PatchParticipantsParticipantIDConfirm)
		r.Get("/trips", wrapper.GetTrips)
		r.Post("/trips", wrapper.PostTrips)
		r.Get("/trips/{tripId}", wrapper.GetTripsTripID)
		r.Put("/trips/{tripId}", wrapper.PutTripsTripID)
		r.Get("/trips/{tripId}/activities", wrapper.GetTripsTripIDActivities)
		r.Post("/trips/{tripId}/activities", wrapper.PostTripsTripIDActivities)
		r.Get("/trips/{tripId}/confirm", wrapper.GetTripsTripIDConfirm)
		r.Post("/trips/{tripId}/invites", wrapper.PostTripsTripIDInvites)
		r.Get("/trips/{tripId}/links", wrapper.GetTripsTripIDLinks)
		r.Post("/trips/{tripId}/links", wrapper.PostTripsTripIDLinks)
		r.Get("/trips/{tripId}/participants", wrapper.GetTripsTripIDParticipants)
	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"5x4cjfh4PlQ3lst/AwAA//8vxohsEy0AAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

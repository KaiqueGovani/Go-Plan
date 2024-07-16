package mailpit

import (
	"context"
	"fmt"
	"journey/internal/pgstore"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
)

type store interface {
	GetTrip(context.Context, uuid.UUID) (pgstore.Trip, error)
	GetParticipants(context.Context, uuid.UUID) ([]pgstore.Participant, error)
}

type Mailpit struct {
	store store
}

func NewMailpit(pool *pgxpool.Pool) Mailpit {
	return Mailpit{pgstore.New(pool)}
}

func (mp Mailpit) SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error {
	ctx := context.Background()
	trip, err := mp.store.GetTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("mailpit: failed to get trip for SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg := mail.NewMsg()
	if err := msg.From("mailpit@journey.com"); err != nil {
		return fmt.Errorf("mailpit: failed to set From in email for SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := msg.To(trip.OwnerEmail); err != nil {
		return fmt.Errorf("mailpit: failed to set To in email for SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg.Subject("Confirm your trip");

	msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
		Olá, %s!
		
		A sua Viagem para %s que começa em %s precisa ser confirmada.
		Clique no botão abaixo para confirmar.
	`, trip.OwnerName, trip.Destination, trip.StartsAt.Time.Format(time.DateOnly),
	))

	client, err := mail.NewClient("mailpit", mail.WithTLSPortPolicy(mail.NoTLS), mail.WithPort(1025))
	if err != nil {
		return fmt.Errorf("mailpit: failed to create email client for SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("mailpit: failed to send email for SendConfirmTripEmailToTripOwner: %w", err)
	}

	return nil
}

func (mp Mailpit) SendConfirmTripEmailToTripParticipants(tripID uuid.UUID) error {
	ctx := context.Background()
	trip, err := mp.store.GetTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("mailpit: failed to get trip for SendConfirmTripEmailToTripParticipants: %w", err)
	}

	participants, err := mp.store.GetParticipants(ctx, tripID)
	if err != nil {
		return fmt.Errorf("mailpit: failed to get participants for SendConfirmTripEmailToTripParticipants: %w", err)
	}

	for _, participant := range participants {
		msg := mail.NewMsg()
		if err := msg.From("mailpit@journey.com"); err != nil {
			return fmt.Errorf("mailpit: failed to set From in email for SendConfirmTripEmailToTripParticipants: %w", err)
		}

		if err := msg.To(participant.Email); err != nil {
			return fmt.Errorf("mailpit: failed to set To in email for SendConfirmTripEmailToTripParticipants: %w", err)
		}

		msg.Subject("Confirm your trip");

		msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
			Olá!
			
			A sua Viagem com %s para %s que começa em %s precisa de sua confirmação.
			Clique no botão abaixo e confirme sua presença.
		`, trip.OwnerName, trip.Destination, trip.StartsAt.Time.Format(time.DateOnly),
		))

		client, err := mail.NewClient("mailpit", mail.WithTLSPortPolicy(mail.NoTLS), mail.WithPort(1025))
		if err != nil {
			return fmt.Errorf("mailpit: failed to create email client for SendConfirmTripEmailToTripParticipants: %w", err)
		}

		if err := client.DialAndSend(msg); err != nil {
			return fmt.Errorf("mailpit: failed to send email for SendConfirmTripEmailToTripParticipants: %w", err)
		}
	}

	return nil
}
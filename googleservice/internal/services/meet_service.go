package services

import (
	"context"
	"fmt"
	"google-service/internal/config"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type MeetService struct {
	calendarService *calendar.Service
}

func NewMeetService(config config.ServiceAccountConfig) (*MeetService, error) {
	ctx := context.Background()

	// Create JWT config from service account credentials
	jwtConfig, err := google.JWTConfigFromJSON(config.CredentialsJSON,
		calendar.CalendarEventsScope,
		calendar.CalendarScope,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT config: %v", err)
	}

	// Create calendar service with the JWT client
	calendarService, err := calendar.NewService(ctx,
		option.WithHTTPClient(jwtConfig.Client(ctx)))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %v", err)
	}

	return &MeetService{
		calendarService: calendarService,
	}, nil
}

// CreateMeetLink generates a new Google Meet link
func (s *MeetService) CreateMeetLink(meetingTitle string) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Create calendar event with Meet conferencing
	event := &calendar.Event{
		Summary: meetingTitle,
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(time.Hour).Format(time.RFC3339),
			TimeZone: "UTC",
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: fmt.Sprintf("%d", time.Now().UnixNano()),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
	}

	// Insert the event and create Meet link
	event, err := s.calendarService.Events.Insert("primary", event).
		Context(ctx).
		ConferenceDataVersion(1).
		Do()
	if err != nil {
		fmt.Printf("Error creating meet link: %v", err)
		return "", fmt.Errorf("failed to create event: %v", err)
	}

	// Extract Meet link from conference data
	if event.ConferenceData != nil && len(event.ConferenceData.EntryPoints) > 0 {
		for _, entryPoint := range event.ConferenceData.EntryPoints {
			if entryPoint.EntryPointType == "video" {
				return entryPoint.Uri, nil
			}
		}
	}

	return "", fmt.Errorf("no meet link found in created event")
}

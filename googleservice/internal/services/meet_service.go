package services

import (
	"context"
	"fmt"
	"google-service/internal/config"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type MeetService struct {
	config *config.GoogleOAuthConfig
}

func NewMeetService(config *config.GoogleOAuthConfig) (*MeetService, error) {
	return &MeetService{
		config: config,
	}, nil
}

// CreateMeetLink generates a new Google Meet link using user's OAuth token
func (s *MeetService) CreateMeetLink(accessToken string, meetingTitle string) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Create OAuth2 token
	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	// Create calendar service with OAuth token
	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(
		oauth2.StaticTokenSource(token),
	))
	if err != nil {
		return "", fmt.Errorf("failed to create calendar service: %v", err)
	}

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
	event, err = calendarService.Events.Insert("primary", event).
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

// GetTokenFromCode exchanges auth code for token
func (s *MeetService) GetTokenFromCode(code string) (*oauth2.Token, error) {
	ctx := context.Background()

	// Create OAuth2 config
	config := &oauth2.Config{
		ClientID:     s.config.ClientID,
		ClientSecret: s.config.ClientSecret,
		RedirectURL:  s.config.RedirectURL,
		Scopes:       append(s.config.Scopes, calendar.CalendarEventsScope),
		Endpoint:     s.config.Endpoint,
	}

	// Exchange code for token
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	return token, nil
}

// CreateMeetLinkWithCode creates a meet link using an auth code
func (s *MeetService) CreateMeetLinkWithCode(code string, meetingTitle string) (string, error) {
	token, err := s.GetTokenFromCode(code)
	if err != nil {
		return "", err
	}

	return s.CreateMeetLink(token.AccessToken, meetingTitle)
}

func (s *MeetService) CreateMeetLinkWithOAUTH(token *oauth2.Token, meetingTitle string) (string, error) {
	return s.CreateMeetLink(token.AccessToken, meetingTitle)
}

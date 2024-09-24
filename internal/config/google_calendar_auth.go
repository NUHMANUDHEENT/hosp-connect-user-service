package config

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)



var (
    googleOauthConfig *oauth2.Config
    oauthStateString  = "random" // For security purposes
)

func init() {
    googleOauthConfig = &oauth2.Config{
        ClientID:     "YOUR_CLIENT_ID",
        ClientSecret: "YOUR_CLIENT_SECRET",
        RedirectURL:  "http://localhost:8080/callback", // Your redirect URI
        Scopes:       []string{calendar.CalendarScope}, // Google Calendar scope
        Endpoint:     google.Endpoint,
    }
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
    url := googleOauthConfig.AuthCodeURL(oauthStateString)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
    if state != oauthStateString {
        http.Error(w, "Invalid state", http.StatusBadRequest)
        return
    }

    code := r.FormValue("code")
    token, err := googleOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
        return
    }

    client := googleOauthConfig.Client(context.Background(), token)
    calendarService, err := calendar.New(client)
    if err != nil {
        http.Error(w, "Failed to create calendar service", http.StatusInternalServerError)
        return
    }

    // Here, you can list or create calendar events
    listEvents(calendarService, w)
}

func listEvents(srv *calendar.Service, w http.ResponseWriter) {
    events, err := srv.Events.List("primary").Do()
    if err != nil {
        fmt.Fprintf(w, "Unable to retrieve calendar events: %v", err)
        return
    }

    for _, item := range events.Items {
        fmt.Fprintf(w, "Event: %s (%s)\n", item.Summary, item.Start.DateTime)
    }
}

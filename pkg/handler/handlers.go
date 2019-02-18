package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"

	"github.com/seregayoga/twitchstream/pkg/twitch"
)

const (
	stateCallbackKey = "oauth-state-callback"
	oauthSessionName = "oauth-session"
	oauthTokenKey    = "oauth-token"
	streamerNameKey  = "streamer-name"
)

// Handlers handlers
type Handlers struct {
	oauth2Config *oauth2.Config
	cookieStore  *sessions.CookieStore
	twitchAPI    *twitch.API
}

// NewHandlers creates Handlers
func NewHandlers(oauth2Config *oauth2.Config, cookieStore *sessions.CookieStore, twitchAPI *twitch.API) *Handlers {
	return &Handlers{
		oauth2Config: oauth2Config,
		cookieStore:  cookieStore,
		twitchAPI:    twitchAPI,
	}
}

// HumanReadableError represents error information
// that can be fed back to a human user.
//
// This prevents internal state that might be sensitive
// being leaked to the outside world.
type HumanReadableError interface {
	HumanError() string
	HTTPCode() int
}

// HumanReadableWrapper implements HumanReadableError
type HumanReadableWrapper struct {
	ToHuman string
	Code    int
	error
}

// HumanError returns human error
func (h HumanReadableWrapper) HumanError() string {
	return h.ToHuman
}

// HTTPCode returns http code
func (h HumanReadableWrapper) HTTPCode() int {
	return h.Code
}

// AnnotateError wraps an error with a message that is intended for a human end-user to read,
// plus an associated HTTP error code.
func AnnotateError(err error, annotation string, code int) error {
	if err == nil {
		return nil
	}
	return HumanReadableWrapper{ToHuman: annotation, error: err}
}

// Handler handler type
type Handler func(http.ResponseWriter, *http.Request) error

// HandleIndex is a Handler that shows a login button. In production, if the frontend is served / generated
// by Go, it should use html/template to prevent XSS attacks.
func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, indexHTML)

	return nil
}

// HandleLogin is a Handler that redirects the user to Twitch for login, and provides the 'state'
// parameter which protects against login CSRF.
func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) (err error) {
	session, err := h.cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	streamerName := r.URL.Query().Get("name")
	if streamerName == "" {
		return AnnotateError(err, "Empty streamer name!", http.StatusBadRequest)
	}
	session.Values[streamerNameKey] = streamerName

	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		return AnnotateError(err, "Couldn't generate a session!", http.StatusInternalServerError)
	}

	state := hex.EncodeToString(tokenBytes[:])

	session.AddFlash(state, stateCallbackKey)

	if err = session.Save(r, w); err != nil {
		return
	}

	http.Redirect(w, r, h.oauth2Config.AuthCodeURL(state), http.StatusTemporaryRedirect)

	return
}

// HandleOAuth2Callback is a Handler for oauth's 'redirect_uri' endpoint;
// it validates the state token and retrieves an OAuth token from the request parameters.
func (h *Handlers) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (err error) {
	session, err := h.cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	switch stateChallenge, state := session.Flashes(stateCallbackKey), r.FormValue("state"); {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge[0]:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'", state, stateChallenge[0])
	}

	if err != nil {
		session.Save(r, w)

		return AnnotateError(
			err,
			"Couldn't verify your confirmation, please try again.",
			http.StatusBadRequest,
		)
	}

	token, err := h.oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		session.Save(r, w)

		return
	}

	// add the oauth token to session
	session.Values[oauthTokenKey] = token

	if err := session.Save(r, w); err != nil {
		log.Printf("error saving session: %s", err)
	}

	http.Redirect(w, r, "/stream", http.StatusTemporaryRedirect)

	return
}

// HandleStream handles strem events page
func (h *Handlers) HandleStream(w http.ResponseWriter, r *http.Request) error {
	session, err := h.cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	streamerName, ok := session.Values[streamerNameKey].(string)
	if !ok {
		return errors.New("empty streamer name")
	}

	_, err = fmt.Fprintf(w, streamHTML, streamerName)

	return err
}

// HandleEvents sends stream events by websockets
func (h *Handlers) HandleEvents(w http.ResponseWriter, r *http.Request) error {
	session, err := h.cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	token, ok := session.Values[oauthTokenKey].(*oauth2.Token)
	if !ok {
		return errors.New("empty token")
	}

	streamerName, ok := session.Values[streamerNameKey].(string)
	if !ok {
		return errors.New("empty streamer name")
	}

	user, err := h.twitchAPI.GetUser(token.AccessToken)
	if err != nil {
		return err
	}

	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ircToWSStreamer := twitch.NewIRCToWSStreamer(ws, user.Name, streamerName, token.AccessToken)

	return ircToWSStreamer.Stream()
}

// HandleFunc wraps handler with helpful middleware
func HandleFunc(path string, handler Handler) {
	http.Handle(path, errorHandling(middleware(handler)))
}

func middleware(h Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// parse POST body, limit request size
		if err = r.ParseForm(); err != nil {
			return AnnotateError(err, "Something went wrong! Please try again.", http.StatusBadRequest)
		}

		return h(w, r)
	}
}

// errorHandling is a middleware that centralises error handling.
// this prevents a lot of duplication and prevents issues where a missing
// return causes an error to be printed, but functionality to otherwise continue
// see https://blog.golang.org/error-handling-and-go
func errorHandling(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			errorString := "Something went wrong! Please try again."
			errorCode := 500

			if v, ok := err.(HumanReadableError); ok {
				errorString, errorCode = v.HumanError(), v.HTTPCode()
			}

			log.Println(err)
			w.Write([]byte(errorString))
			w.WriteHeader(errorCode)
			return
		}
	})
}

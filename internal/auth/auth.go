package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type session struct {
	Id         string
	UserId     int
	Expiration time.Time
}

type sessionsContainer struct {
	s map[string]session
	mu       sync.Mutex
}

var (
	ErrLogin = errors.New("error logging in user")
	// the map must not be used directly for add and delete operations. helper functions below must be used instead
	sessions = sessionsContainer{
		s: make(map[string]session),
	}
)

func (s *sessionsContainer) addSession(newSession session) {
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	sessions.s[newSession.Id] = newSession
}

func (s *sessionsContainer) deleteSession(sessionId string) {
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	delete(sessions.s, sessionId)
}

func (s *sessionsContainer) getSession(sessionId string) (session, bool) {
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	session, ok := sessions.s[sessionId]
	return session, ok
}

func init() {
	go func() {
		for {
			time.Sleep(time.Hour * 12)
			sessions.mu.Lock()
			for sessionId, sessionValue := range sessions.s {
				if sessionValue.Expiration.Before(time.Now().UTC()) {
					sessions.deleteSession(sessionId)
				}
			}
			sessions.mu.Unlock()
		}
	}()
}

const (
	sessionCookieName         = "session_id"
	unauthenticatedErrMessage = "Сіз жүйеге кіруіңіз керек"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !Authenticated(r) {
			redirectLogin(w, r, unauthenticatedErrMessage)
			return
		}
		next(w, r)
	}
}

func Authenticated(r *http.Request) bool {
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return false
	}
	sessionId := sessionCookie.Value
	session, ok := sessions.getSession(sessionId)
	if ok && session.UserId != 0 && time.Now().UTC().Before(session.Expiration) {
		return true
	}
	return false
}

func redirectLogin(w http.ResponseWriter, r *http.Request, message string) {
	redirectUrl := fmt.Sprintf("/login?error=%s", url.PathEscape(message))
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func redirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

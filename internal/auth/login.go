package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/nurtai325/kaspi-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionExpirationDuration = time.Hour * 24 * 14
)

func Login(w http.ResponseWriter, r *http.Request, phone, password string) error {
	userRepo := repositories.NewUser()
	user, err := userRepo.OneByPhone(r.Context(), phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			redirectLogin(w, r, "Бұл номермен қолданушы табылмады")
			return ErrLogin
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		redirectLogin(w, r, "Құпиясөз қате")
		return ErrLogin
	}
	b := make([]byte, 32)
	rand.Read(b)
	sessionId := base64.URLEncoding.EncodeToString(b)
	sessions.addSession(session{
		Id:         sessionId,
		UserId:     user.Id,
		Expiration: time.Now().UTC().Add(sessionExpirationDuration),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionId,
		MaxAge:   3600 * 24 * 14,
		// TODO: uncomment this
		// Secure:   true,
		Secure: false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})
	redirectRoot(w, r)
	return nil
}

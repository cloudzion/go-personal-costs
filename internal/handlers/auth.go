package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Генерация случайного токена для сессии
func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Middleware авторизации: проверяет cookie и кладет userID в контекст запроса
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := Repo.GetUserIDBySession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Передаем ID авторизованного пользователя через context.Context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// Регистрация (/register)
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("web/templates/layout.html", "web/templates/register.html")
		tmpl.ExecuteTemplate(w, "layout", nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Хэшируем пароль через bcrypt
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка шифрования", http.StatusInternalServerError)
			return
		}

		err = Repo.CreateUser(email, string(hash))
		if err != nil {
			tmpl, _ := template.ParseFiles("web/templates/layout.html", "web/templates/register.html")
			tmpl.ExecuteTemplate(w, "layout", map[string]string{"Error": "Пользователь с таким Email уже существует"})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Вход (/login)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("web/templates/layout.html", "web/templates/login.html")
		tmpl.ExecuteTemplate(w, "layout", nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := Repo.GetUserByEmail(email)
		if err != nil {
			tmpl, _ := template.ParseFiles("web/templates/layout.html", "web/templates/login.html")
			tmpl.ExecuteTemplate(w, "layout", map[string]string{"Error": "Неверный Email или пароль"})
			return
		}

		// Проверяем пароль через bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			tmpl, _ := template.ParseFiles("web/templates/layout.html", "web/templates/login.html")
			tmpl.ExecuteTemplate(w, "layout", map[string]string{"Error": "Неверный Email или пароль"})
			return
		}

		// Создаем сессию в БД и ставим Cookie
		sessionID := generateSessionID()
		expires := time.Now().Add(24 * time.Hour)
		Repo.CreateSession(sessionID, user.ID, expires)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  expires,
			HttpOnly: true,
			Path:     "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Выход (/logout)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		Repo.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

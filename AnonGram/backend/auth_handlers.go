package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		Bio         string `json:"bio"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		jsonWrite(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	in.Username = strings.ToLower(strings.TrimSpace(in.Username))
	if len(in.Username) < 3 || len(in.Password) < 8 {
		jsonWrite(w, 400, map[string]string{"error": "username >=3 and password >=8"})
		return
	}
	h, e := hashPassword(in.Password)
	if e != nil {
		jsonWrite(w, 500, map[string]string{"error": "hash failed"})
		return
	}
	id := uuid.NewString()
	if _, e = a.db.Exec(r.Context(), `insert into users(id,username,display_name,password_hash,bio) values($1,$2,$3,$4,$5)`, id, in.Username, in.DisplayName, h, in.Bio); e != nil {
		jsonWrite(w, 409, map[string]string{"error": "username already exists"})
		return
	}
	access, _ := makeToken(id, "user", 30*time.Minute)
	refresh, _ := makeToken(id, "refresh", 30*24*time.Hour)
	jsonWrite(w, 201, map[string]string{"access_token": access, "refresh_token": refresh, "user_id": id})
}
func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		jsonWrite(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	var id, h string
	e := a.db.QueryRow(r.Context(), `select id,password_hash from users where username=$1`, strings.ToLower(strings.TrimSpace(in.Username))).Scan(&id, &h)
	if e != nil || !checkPassword(h, in.Password) {
		jsonWrite(w, 401, map[string]string{"error": "invalid credentials"})
		return
	}
	access, _ := makeToken(id, "user", 30*time.Minute)
	refresh, _ := makeToken(id, "refresh", 30*24*time.Hour)
	jsonWrite(w, 200, map[string]string{"access_token": access, "refresh_token": refresh, "user_id": id})
}
func (a *App) refresh(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		jsonWrite(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	c, e := ParseAccessToken(in.RefreshToken)
	if e != nil || c.Role != "refresh" {
		jsonWrite(w, 401, map[string]string{"error": "invalid refresh token"})
		return
	}
	access, _ := makeToken(c.Subject, "user", 30*time.Minute)
	jsonWrite(w, 200, map[string]string{"access_token": access})
}

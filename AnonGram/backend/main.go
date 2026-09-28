package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type App struct {
	db  *pgxpool.Pool
	rdb *redis.Client
	hub *Hub
}
type userKey struct{}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func jsonWrite(w http.ResponseWriter, s int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(v)
}
func bearer(r *http.Request) string {
	p := r.Header.Get("Authorization")
	if strings.HasPrefix(p, "Bearer ") {
		return strings.TrimPrefix(p, "Bearer ")
	}
	return ""
}
func main() {
	ctx := context.Background()
	db, e := pgxpool.New(ctx, env("DATABASE_URL", "postgres://anongram:anongram@localhost:5432/anongram?sslmode=disable"))
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	ro, e := redis.ParseURL(env("REDIS_URL", "redis://localhost:6379/0"))
	if e != nil {
		log.Fatal(e)
	}
	rdb := redis.NewClient(ro)
	defer rdb.Close()
	a := &App{db: db, rdb: rdb, hub: NewHub()}
	go a.hub.Run(ctx, rdb)
	m := http.NewServeMux()
	m.HandleFunc("/api/health", a.health)
	m.HandleFunc("/api/v1/auth/register", a.register)
	m.HandleFunc("/api/v1/auth/login", a.login)
	m.HandleFunc("/api/v1/auth/refresh", a.refresh)
	m.HandleFunc("/api/v1/me", a.auth(a.me))
	m.HandleFunc("/api/v1/chats", a.auth(a.chats))
	m.HandleFunc("/api/v1/chats/", a.auth(a.chatRoute))
	m.HandleFunc("/api/v1/search", a.auth(a.search))
	m.HandleFunc("/api/v1/profile/", a.auth(a.profile))
	m.HandleFunc("/ws", a.ws)
	m.Handle("/metrics", promhttp.Handler())
	log.Printf("AnonGram API :%s", env("HTTP_PORT", "8080"))
	log.Fatal(http.ListenAndServe(":"+env("HTTP_PORT", "8080"), cors(m)))
}
func cors(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		n.ServeHTTP(w, r)
	})
}
func (a *App) health(w http.ResponseWriter, r *http.Request) {
	c, x := context.WithTimeout(r.Context(), 2*time.Second)
	defer x()
	d := a.db.Ping(c) == nil
	q := a.rdb.Ping(c).Err() == nil
	s := 200
	if !d || !q {
		s = 503
	}
	jsonWrite(w, s, map[string]any{"ok": d && q, "db": d, "redis": q})
}
func (a *App) auth(n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, e := ParseAccessToken(bearer(r))
		if e != nil || c.Role != "user" {
			jsonWrite(w, 401, map[string]string{"error": "invalid token"})
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), userKey{}, c.Subject))
		n(w, r)
	}
}
func uid(r *http.Request) string { v, _ := r.Context().Value(userKey{}).(string); return v }
func (a *App) me(w http.ResponseWriter, r *http.Request) {
	var id, u, d string
	e := a.db.QueryRow(r.Context(), `select id,username,display_name from users where id=$1`, uid(r)).Scan(&id, &u, &d)
	if e != nil {
		jsonWrite(w, 404, map[string]string{"error": "user not found"})
		return
	}
	jsonWrite(w, 200, map[string]string{"id": id, "username": u, "display_name": d})
}
func (a *App) chats(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.createChat(w, r)
		return
	}
	if r.Method != http.MethodGet {
		jsonWrite(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	rows, e := a.db.Query(r.Context(), `select c.id,c.title,c.type,coalesce(c.avatar_url,'') from chats c join chat_members cm on cm.chat_id=c.id where cm.user_id=$1 order by c.updated_at desc limit 100`, uid(r))
	if e != nil {
		jsonWrite(w, 500, map[string]string{"error": "database error"})
		return
	}
	defer rows.Close()
	out := []map[string]string{}
	for rows.Next() {
		var id, t, ty, av string
		if rows.Scan(&id, &t, &ty, &av) == nil {
			out = append(out, map[string]string{"id": id, "title": t, "type": ty, "avatar_url": av})
		}
	}
	jsonWrite(w, 200, out)
}
func (a *App) chatRoute(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/chats/"), "/")
	if len(p) >= 2 && p[1] == "messages" {
		if r.Method == "GET" {
			a.messages(w, r, p[0])
			return
		}
		if r.Method == "POST" {
			a.send(w, r, p[0])
			return
		}
	}
	jsonWrite(w, 404, map[string]string{"error": "not found"})
}
func (a *App) messages(w http.ResponseWriter, r *http.Request, id string) {
	if !a.isMember(r.Context(), id, uid(r)) {
		jsonWrite(w, http.StatusForbidden, map[string]string{"error": "not a chat member"})
		return
	}
	rows, e := a.db.Query(r.Context(), `select m.id,m.sender_id,u.username,m.body,m.created_at from messages m join users u on u.id=m.sender_id where m.chat_id=$1 order by m.created_at desc limit 100`, id)
	if e != nil {
		jsonWrite(w, 500, map[string]string{"error": "database error"})
		return
	}
	defer rows.Close()
	out := []map[string]any{}
	for rows.Next() {
		var i, s, u, b string
		var t time.Time
		if rows.Scan(&i, &s, &u, &b, &t) == nil {
			out = append(out, map[string]any{"id": i, "sender_id": s, "username": u, "body": b, "created_at": t})
		}
	}
	jsonWrite(w, 200, out)
}
func (a *App) send(w http.ResponseWriter, r *http.Request, chat string) {
	if !a.isMember(r.Context(), chat, uid(r)) {
		jsonWrite(w, http.StatusForbidden, map[string]string{"error": "not a chat member"})
		return
	}
	var in struct {
		Body string `json:"body"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil || strings.TrimSpace(in.Body) == "" {
		jsonWrite(w, 400, map[string]string{"error": "body required"})
		return
	}
	var id string
	e := a.db.QueryRow(r.Context(), `insert into messages(chat_id,sender_id,body) values($1,$2,$3) returning id`, chat, uid(r), strings.TrimSpace(in.Body)).Scan(&id)
	if e != nil {
		jsonWrite(w, 500, map[string]string{"error": "database error"})
		return
	}
	msg := map[string]any{"id": id, "chat_id": chat, "sender_id": uid(r), "body": strings.TrimSpace(in.Body), "created_at": time.Now()}
	b, _ := json.Marshal(msg)
	a.hub.Broadcast(chat, b)
	jsonWrite(w, 201, msg)
}
func (a *App) isMember(ctx context.Context, chatID, userID string) bool {
	var ok bool
	err := a.db.QueryRow(ctx, `select exists(select 1 from chat_members where chat_id=$1 and user_id=$2)`, chatID, userID).Scan(&ok)
	return err == nil && ok
}

func (a *App) createChat(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title     string   `json:"title"`
		Type      string   `json:"type"`
		MemberIDs []string `json:"member_ids"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		jsonWrite(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		jsonWrite(w, http.StatusBadRequest, map[string]string{"error": "title required"})
		return
	}
	if in.Type == "" {
		in.Type = "group"
	}
	if in.Type != "private" && in.Type != "group" && in.Type != "channel" {
		jsonWrite(w, http.StatusBadRequest, map[string]string{"error": "invalid chat type"})
		return
	}
	tx, err := a.db.Begin(r.Context())
	if err != nil {
		jsonWrite(w, 500, map[string]string{"error": "database error"})
		return
	}
	defer tx.Rollback(r.Context())
	var chatID string
	if err = tx.QueryRow(r.Context(), `insert into chats(title,type) values($1,$2) returning id`, in.Title, in.Type).Scan(&chatID); err != nil {
		jsonWrite(w, 500, map[string]string{"error": "chat creation failed"})
		return
	}
	members := append([]string{uid(r)}, in.MemberIDs...)
	seen := map[string]bool{}
	for _, member := range members {
		if member == "" || seen[member] {
			continue
		}
		seen[member] = true
		if _, err = tx.Exec(r.Context(), `insert into chat_members(chat_id,user_id,role) values($1,$2,$3) on conflict do nothing`, chatID, member, map[bool]string{true: "owner", false: "member"}[member == uid(r)]); err != nil {
			jsonWrite(w, 500, map[string]string{"error": "member add failed"})
			return
		}
	}
	if err = tx.Commit(r.Context()); err != nil {
		jsonWrite(w, 500, map[string]string{"error": "commit failed"})
		return
	}
	jsonWrite(w, http.StatusCreated, map[string]string{"id": chatID, "title": in.Title, "type": in.Type})
}

func (a *App) search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	rows, e := a.db.Query(r.Context(), `select id,username,display_name from users where search_vector @@ plainto_tsquery('simple',$1) order by ts_rank(search_vector,plainto_tsquery('simple',$1)) desc limit 30`, q)
	if e != nil {
		jsonWrite(w, 500, map[string]string{"error": "database error"})
		return
	}
	defer rows.Close()
	out := []map[string]string{}
	for rows.Next() {
		var i, u, d string
		if rows.Scan(&i, &u, &d) == nil {
			out = append(out, map[string]string{"id": i, "username": u, "display_name": d})
		}
	}
	jsonWrite(w, 200, out)
}
func (a *App) profile(w http.ResponseWriter, r *http.Request) {
	u := strings.TrimPrefix(r.URL.Path, "/api/v1/profile/")
	var id, un, d, b string
	e := a.db.QueryRow(r.Context(), `select id,username,display_name,bio from users where username=$1`, u).Scan(&id, &un, &d, &b)
	if e != nil {
		jsonWrite(w, 404, map[string]string{"error": "profile not found"})
		return
	}
	jsonWrite(w, 200, map[string]string{"id": id, "username": un, "display_name": d, "bio": b})
}
func (a *App) ws(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := ParseAccessToken(token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true // Restrict this to your web origin in production.
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 32)}
	a.hub.Add(client)
	defer a.hub.Remove(client)

	go func() {
		for msg := range client.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var in struct {
			Type   string `json:"type"`
			ChatID string `json:"chat_id"`
			Body   string `json:"body"`
		}
		if json.Unmarshal(raw, &in) != nil {
			continue
		}
		switch in.Type {
		case "subscribe":
			if in.ChatID != "" && a.isMember(r.Context(), in.ChatID, claims.Subject) {
				client.chatID = in.ChatID
			}
		case "ping":
			select {
			case client.send <- []byte(`{"type":"pong"}`):
			default:
			}
		case "typing":
			payload, _ := json.Marshal(map[string]any{"type": "typing", "chat_id": in.ChatID, "user_id": claims.Subject, "active": in.Body != ""})
			a.hub.Broadcast(in.ChatID, payload)
		}
	}
}

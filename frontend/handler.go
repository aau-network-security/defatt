package frontend

import (
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/aau-network-security/defat/database"
	"github.com/aau-network-security/defat/game"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
)

var (
	sessionName     = "defatt"
	flashSession    = "defatt_flash"
	contextTeamKey  = "defatt-user"
	contextEventKey = "defatt-event"
	ErrUnknownGame  = errors.New("Game does not exist")
)

type content struct {
	Event *game.GameConfig
	User  *database.GameUser
}

type Event struct {
	ID       string
	Name     string
	Tag      string
	Scenario Scenario
}

type Scenario struct {
	ID         string
	Name       string
	Duration   time.Duration
	Difficulty string
	AlertLimit string
	Networks   map[string]string
}

type Web struct {
	m             sync.RWMutex
	Router        *mux.Router
	ServerBind    string
	ServerBindTLS string
	Domain        string
	CertKey       string
	CertFile      string
	cookieStore   *sessions.CookieStore
	Templates     map[string]*template.Template
	Events        map[string]*game.GameConfig
}
type vpnConf struct {
	IPAddress    string
	PrivateKey   string
	ServerPubKey string
	AllowedIPs   string
	Endpoint     string
}

func init() {
	gob.Register(flashMessage{})
	gob.Register(database.GameUser{})
}

func New(serverbind, serverbindTLS, domain, certKey, certFile string) (*Web, error) {
	w := Web{
		Router:        mux.NewRouter(),
		ServerBind:    serverbind,
		ServerBindTLS: serverbindTLS,
		Domain:        domain,
		CertKey:       certKey,
		CertFile:      certFile,
		Events:        make(map[string]*game.GameConfig),
		cookieStore:   sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)),
	}
	if w.Templates == nil {
		w.Templates = make(map[string]*template.Template)

		w.parseTemplate("", "index")
		w.parseTemplate("index", "")
		w.parseTemplate("login", "")
		w.parseTemplate("signup", "")
		w.parseTemplate("landing", "")

	}

	return &w, nil
}

func (w *Web) Run() error {
	// setup routes
	if err := w.Routes(); err != nil {
		return err
	}
	if w.CertKey == "" || w.CertFile == "" {
		log.Info().Str("bind", w.ServerBind).Msg("no cert Files running HTTP only")
		if err := http.ListenAndServe(w.ServerBind, w.Router); err != nil {
			return err
		}
	}
	go http.ListenAndServe(w.ServerBind, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	}))

	if err := http.ListenAndServeTLS(w.ServerBindTLS, w.CertFile, w.CertKey, w.Router); err != nil {
		log.Warn().Msgf("Serving error: %s", err)
	}

	// run the webserver
	// log.Info().Str("bind", w.ServerBind).Msg("running server")
	// if err := http.ListenAndServe(w.ServerBind, w.Router); err != nil {
	// 	return err
	// }

	return nil
}

func (w *Web) Routes() error {
	subrouter(w.Router, "/", func(r *mux.Router) {
		subrouter(r.Host("{subdomain:[A-z0-9]+}.localhost:8080").Subrouter(), "/", func(r *mux.Router) {
			// this one will extract the event from each subdomain
			// and attach it to the context
			r.Use(w.middlewareExtractEvent)
			r.Use(w.teamMiddleware)

			r.HandleFunc("/", w.handleIndex)
			r.HandleFunc("/vpn", w.handleVPN).Methods("GET")
			r.HandleFunc("/login", w.handleLoginGet).Methods("GET")
			r.HandleFunc("/login", w.handleLoginPost).Methods("POST")
			r.HandleFunc("/logout", w.handleLogout).Methods("GET")
			r.HandleFunc("/signup", w.handleSignupGet).Methods("GET")
			r.HandleFunc("/signup", w.handleSignupPost).Methods("POST")
			r.HandleFunc("/start", w.handleStartGame).Methods("GET")
			r.PathPrefix("/assets/").Handler(http.StripPrefix("", http.FileServer(http.FS(fsStatic))))
		})

		r.HandleFunc("/", noEvent)

	})

	return nil
}

func (w *Web) handleIndex(rw http.ResponseWriter, r *http.Request) {
	var content content
	content.Event = EventFromContext(r.Context())
	content.User = UserFromContext(r.Context())

	if content.User.ID == "" {
		w.templateExec(rw, r, "index", content)
		return
	}
	w.templateExec(rw, r, "landing", content)

}

func (w *Web) handleLogout(rw http.ResponseWriter, r *http.Request) {
	session, err := w.cookieStore.Get(r, sessionName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = database.GameUser{}
	session.Options.MaxAge = -1

	err = session.Save(r, rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusFound)

}

func noEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to NAP, you should have gotten a link in the form of https://eventname.localhost:8080, please use that"))
}

func (w *Web) handleVPN(rw http.ResponseWriter, r *http.Request) {
	var content content
	content.User = UserFromContext(r.Context())
	content.Event = EventFromContext(r.Context())

	vpn, err := content.Event.CreateVPNConfig(r.Context(), false, content.Event.Tag, content.User.ID)
	if err != nil {
		log.Error().Err(err).Str("user", content.User.ID).Interface("VPN conf", vpn).Msg("failed to create vpn conf")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Disposition", `inline; filename="wg_deffat.conf"`)
	rw.Header().Set("Content-Type", "application/txt")

	tmpl := template.Must(template.ParseFiles(templatesBasePath + "wireguard.conf" + templatesExt))

	if err := tmpl.Execute(rw, vpn); err != nil {
		log.Error().Err(err).Str("user", content.User.ID).Interface("VPN conf", vpn).Msg("failed to create vpn conf")
		rw.WriteHeader(http.StatusInternalServerError)
	}

}

func (w *Web) handleLoginGet(rw http.ResponseWriter, r *http.Request) {
	var content content
	content.User = UserFromContext(r.Context())
	content.Event = EventFromContext(r.Context())

	if content.User.ID != "" {
		http.Redirect(rw, r, "/", http.StatusSeeOther)
		return
	}

	w.templateExec(rw, r, "login", content)

}

func (w *Web) handleLoginPost(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error().Err(err).Msg("could not parse login form")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := w.cookieStore.Get(r, sessionName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	event := EventFromContext(r.Context())
	user := UserFromContext(r.Context())

	if user.ID != "" {
		http.Redirect(rw, r, "/", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Username cannot be empty"})
		http.Redirect(rw, r, "/login", http.StatusSeeOther)
		return
	}

	pw := r.FormValue("password")
	if pw == "" {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Password cannot be empty"})
		http.Redirect(rw, r, "/login", http.StatusSeeOther)
		return
	}
	auser, err := database.AuthUser(r.Context(), username, pw, event.ID)
	if err != nil {
		return
	}
	if auser.ID == "" {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "User does not exist"})
		http.Redirect(rw, r, "/login", http.StatusSeeOther)
		return
	}

	session.Values["user"] = auser

	if err := session.Save(r, rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)

}

func (w *Web) handleSignupGet(rw http.ResponseWriter, r *http.Request) {
	var content content
	content.User = UserFromContext(r.Context())
	content.Event = EventFromContext(r.Context())

	if content.User.ID != "" {
		http.Redirect(rw, r, "/", http.StatusSeeOther)
		return
	}

	w.templateExec(rw, r, "signup", content)

}

func (w *Web) handleSignupPost(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error().Err(err).Msg("could not parse add user form")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	session, err := w.cookieStore.Get(r, sessionName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user := UserFromContext(r.Context())
	game := EventFromContext(r.Context())

	if user.ID != "" {
		http.Redirect(rw, r, "/", http.StatusBadRequest)
		return
	}
	email := r.FormValue("signupemail")
	if email == "" {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Email cannot be empty"})
		http.Redirect(rw, r, "/signup", http.StatusSeeOther)
		return
	}

	username := r.FormValue("signupusername")
	if username == "" {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Username cannot be empty"})
		http.Redirect(rw, r, "/signup", http.StatusSeeOther)
		return
	}

	pw := r.FormValue("signuppassword")
	if pw == "" {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Password cannot be empty"})
		http.Redirect(rw, r, "/signup", http.StatusSeeOther)
		return
	}
	pwc := r.FormValue("signupcpassword")
	if pwc == "" {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Confirm Password cannot be empty"})
		http.Redirect(rw, r, "/signup", http.StatusSeeOther)
		return
	}
	if pwc != pw {
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Passwords should match"})
		http.Redirect(rw, r, "/signup", http.StatusSeeOther)
		return
	}
	team := r.FormValue("team")
	if team != "red" {
		if team != "blue" {
			w.addFlash(rw, r, flashMessage{flashLevelWarning, "Wrong team"})
			http.Redirect(rw, r, "/signup", http.StatusBadRequest)
			return
		}
	}

	if team == "red" {
		user, err := database.AddUser(r.Context(), username, email, pw, game.ID, database.RedTeam)
		if err != nil {
			w.addFlash(rw, r, flashMessage{flashLevelWarning, "Database error occcured"})
			http.Redirect(rw, r, "/signup", http.StatusInternalServerError)
			return
		}
		session.Values["user"] = user
		if err := session.Save(r, rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if team == "blue" {
		user, err := database.AddUser(r.Context(), username, email, pw, game.ID, database.BlueTeam)
		if err != nil {
			w.addFlash(rw, r, flashMessage{flashLevelWarning, "Database error occcured"})
			http.Redirect(rw, r, "/signup", http.StatusInternalServerError)
			return
		}
		session.Values["user"] = user
		if err := session.Save(r, rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (w *Web) handleStartGame(rw http.ResponseWriter, r *http.Request) {
	var content content
	content.User = UserFromContext(r.Context())
	content.Event = EventFromContext(r.Context())

	if content.User.ID == "" {
		http.Redirect(rw, r, "/", http.StatusSeeOther)
		return
	}

	w.templateExec(rw, r, "signup", content)
}

func subrouter(origRouter *mux.Router, path string, fn func(r *mux.Router)) {
	fn(origRouter.NewRoute().PathPrefix(path).Subrouter())
}

func writeError(rw http.ResponseWriter, err error, msg string) {
	if err != nil {
		rw.Write([]byte(fmt.Sprintf("%s: %v", msg, err)))
	} else {
		rw.Write([]byte(fmt.Sprintf("%s", msg)))
	}
}

func (w *Web) AddGame(e *game.GameConfig) {
	w.m.Lock()
	defer w.m.Unlock()
	w.Events[e.Tag] = e
}

func (w *Web) GetGame(tag string) (*game.GameConfig, error) {
	w.m.RLock()
	ev, ok := w.Events[tag]
	w.m.RUnlock()
	if !ok {
		return nil, ErrUnknownGame
	}

	return ev, nil
}

func (w *Web) RemoveGame(tag string) error {
	w.m.Lock()
	defer w.m.Unlock()

	if _, ok := w.Events[tag]; !ok {
		return ErrUnknownGame
	}

	delete(w.Events, tag)

	return nil
}

// func main() {
// 	ctx := context.Background()
// 	database.New(ctx, "defatt.db")
// 	defer database.Close()
// 	e := Event{Name: "test", Tag: "test", ID: "Testing"}
// 	w, _ := New(":8080", "localhost")
// 	w.AddGame(&e)
// 	for g := range w.Events {
// 		fmt.Println(g)
// 	}
// 	log.Info().Msgf("%v", w.Run())
// }

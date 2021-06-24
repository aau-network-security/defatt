package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"

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

type Event struct {
	Name string
	Tag  string
}

type Web struct {
	m           sync.RWMutex
	Router      *mux.Router
	ServerBind  string
	Domain      string
	cookieStore *sessions.CookieStore
	Templates   map[string]*template.Template
	Events      map[string]*Event
}

func init() {
	gob.Register(flashMessage{})
}

func New(serverbind, domain string) (*Web, error) {
	w := Web{
		Router:      mux.NewRouter(),
		ServerBind:  serverbind,
		Domain:      domain,
		Events:      make(map[string]*Event),
		cookieStore: sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)),
	}
	if w.Templates == nil {
		w.Templates = make(map[string]*template.Template)

		w.parseTemplate("", "index")
		w.parseTemplate("index", "")
		w.parseTemplate("login", "")
		w.parseTemplate("signup", "")

	}

	return &w, nil
}

func (w *Web) Run() error {
	// setup routes
	if err := w.Routes(); err != nil {
		return err
	}

	// run the webserver
	log.Info().Str("bind", w.ServerBind).Msg("running server")
	if err := http.ListenAndServe(w.ServerBind, w.Router); err != nil {
		return err
	}

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
			r.HandleFunc("/login", w.handleLoginGet).Methods("GET")
			r.HandleFunc("/login", w.handleLoginPost).Methods("POST")
			r.PathPrefix("/assets/").Handler(http.StripPrefix("", http.FileServer(http.FS(fsStatic))))
		})

		r.HandleFunc("/", noEvent)

	})

	return nil
}

func (w *Web) handleIndex(rw http.ResponseWriter, r *http.Request) {
	w.templateExec(rw, r, "index", nil)
}

func noEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to NAP, you should have gotten a link in the form of https://eventname.localhost:8080, please use that"))
}

func (w *Web) handleLoginGet(rw http.ResponseWriter, r *http.Request) {
	team := TeamFromContext(r.Context())
	event := EventFromContext(r.Context())

	if team.ID != "" {
		http.Redirect(rw, r, "/", http.StatusSeeOther)
	}

	w.templateExec(rw, r, "login", event)

}

func (w *Web) handleLoginPost(rw http.ResponseWriter, r *http.Request) {
	log.Info().Msg("loginpost")
	if err := r.ParseForm(); err != nil {
		log.Error().Err(err).Msg("could not parse add user form")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	team := TeamFromContext(r.Context())

	if team.ID != "" {
		http.Redirect(rw, r, "/", http.StatusBadRequest)
	}

	username := r.FormValue("username")
	if username == "" {
		log.Info().Msg("1t")
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Username cannot be empty"})
		http.Redirect(rw, r, "/login", http.StatusSeeOther)
		return
	}

	pw := r.FormValue("password")
	if pw == "" {
		log.Info().Msg("2")
		w.addFlash(rw, r, flashMessage{flashLevelWarning, "Password cannot be empty"})
		http.Redirect(rw, r, "/login", http.StatusSeeOther)
		return
	}
	fmt.Println(pw, username)

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

func (w *Web) AddGame(e *Event) {
	w.m.Lock()
	defer w.m.Unlock()
	w.Events[e.Tag] = e
}

func (w *Web) GetGame(tag string) (*Event, error) {
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

func main() {
	e := Event{Name: "test", Tag: "test"}
	w, _ := New(":8080", "localhost")
	w.AddGame(&e)
	for g := range w.Events {
		fmt.Println(g)
	}
	log.Info().Msgf("%v", w.Run())
}

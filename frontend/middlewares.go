package frontend

import (
	"context"
	"net/http"

	"github.com/aau-network-security/defatt/database"
	"github.com/aau-network-security/defatt/game"
	"github.com/gorilla/mux"
)

func (w *Web) middlewareExtractEvent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// extract our vars from our request
		vars := mux.Vars(r)

		// get our subdomain, if any?
		subdomain, hasSubdomain := vars["subdomain"]
		// if we can't find any subdomain, we just throw a error ba ck at them
		if !hasSubdomain {
			writeError(rw, nil, "couldn't find any subdomain, aborting")
			return
		}

		// check if our subdomain has a event that matches
		event, err := w.GetGame(subdomain)
		if err != nil {
			writeError(rw, nil, "no event for that subdomain")
			return
		}

		// make new context with the event in it
		r = r.WithContext(context.WithValue(r.Context(), contextEventKey, event))

		// serve the next route with the event attached to the context
		next.ServeHTTP(rw, r)
	})
}

func (w *Web) teamMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		s, err := w.cookieStore.Get(r, sessionName)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		val := s.Values["user"]
		team, _ := val.(database.GameUser)

		newRequest := r.WithContext(context.WithValue(r.Context(), contextTeamKey, &team))
		s.Save(newRequest, rw)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, newRequest)
	})
}

func EventFromContext(ctx context.Context) *game.GameConfig {
	// we discard our OK value, as we do not need it,
	// but only so we do not panic! As if we cannot
	// get any value, then we get a nil pointer, and
	// we will use that to determine if we have
	// any information on the event
	event, _ := ctx.Value(contextEventKey).(*game.GameConfig)
	return event
}

func UserFromContext(ctx context.Context) *database.GameUser {
	// we discard our OK value, as we do not need it,
	// but only so we do not panic! As if we cannot
	// get any value, then we get a nil pointer, and
	// we will use that to determine if we have
	// any information on the event
	user, _ := ctx.Value(contextTeamKey).(*database.GameUser)
	return user
}

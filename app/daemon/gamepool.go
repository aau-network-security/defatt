package daemon

import (
	"net/http"
	"strings"
	"sync"

	wg "github.com/aau-network-security/defat/app/daemon/vpn-proto"
	"github.com/aau-network-security/defat/frontend"
	"github.com/aau-network-security/defat/store"
)

type gamepool struct {
	m               sync.RWMutex
	host            string
	notFoundHandler http.Handler

	handlers map[store.Tag]http.Handler
	wg       wg.WireguardClient
}

func NewGamePool(host string) *gamepool {
	return &gamepool{
		host:            host,
		notFoundHandler: notFoundHandler(),
		handlers:        map[store.Tag]http.Handler{},
	}
}

//
func (gp *gamepool) AddGame(tag string, gm *frontend.WebSite) {
	gp.m.Lock()
	defer gp.m.Unlock()
	gp.handlers[store.Tag(tag)] = gm.Handler()
}

func (gp *gamepool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainParts := strings.SplitN(r.Host, ".", 2)

	if len(domainParts) != 2 {
		gp.notFoundHandler.ServeHTTP(w, r)
		return
	}

	sub, dom := domainParts[0], domainParts[1]
	if !strings.HasPrefix(dom, gp.host) {
		gp.notFoundHandler.ServeHTTP(w, r)
		return
	}

	gp.m.RLock()
	mux, ok := gp.handlers[store.Tag(sub)]
	gp.m.RUnlock()
	if !ok {
		gp.notFoundHandler.ServeHTTP(w, r)
		return
	}

	mux.ServeHTTP(w, r)
}

//
//func (gp *gamepool) RemoveGame(t store.Tag) error {
//	gp.m.Lock()
//	defer gp.m.Unlock()
//
//	if _, ok := gp.Games[t]; !ok {
//		return UnknownGameErr
//	}
//
//	delete(gp.Games, t)
//	delete(gp.handlers, t)
//
//	return nil
//}
//
//func (gp *gamepool) GetGame(t store.Tag) (game, error) {
//	gp.m.RLock()
//	ev, ok := gp.Games[t]
//	gp.m.RUnlock()
//	if !ok {
//		return nil, UnknownGameErr
//	}
//
//	return ev, nil
//}
//
//func (gp *gamepool) GetAllGames() []game {
//	Games := make([]game, len(gp.Games))
//
//	var i int
//	gp.m.RLock()
//	for _, ev := range gp.Games {
//		Games[i] = ev
//		i++
//	}
//	gp.m.RUnlock()
//
//	return Games
//}
//
//func (gp *gamepool) Close() error {
//	var firstErr error
//
//	for _, ev := range gp.Games {
//		if err := ev.Close(); err != nil && firstErr == nil {
//			firstErr = err
//		}
//	}
//
//	return firstErr
//}
//
//
//func getHost(r *http.Request) string {
//	if r.URL.IsAbs() {
//		host := r.Host
//		// Slice off any port information.
//		if i := strings.Index(host, ":"); i != -1 {
//			host = host[:i]
//		}
//		return host
//	}
//	return r.URL.Host
//
//}
//
func notFoundHandler() http.Handler {
	p := []byte("404 NOT FOUND")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(p)
	})
}

package game

import (
	"context"
	"io"
	"net/http"

	wg "github.com/aau-network-security/defat/app/daemon/vpn-proto"
	"github.com/aau-network-security/defat/frontend"
	"github.com/aau-network-security/defat/store"
	"github.com/aau-network-security/defat/virtual/docker"
)

type Game interface {
	Start(context.Context) error
	Close() error
	Suspend(context.Context) error
	Resume(context.Context) error

	Finish(string)
	//AssignLab(*store.Team, lab.Lab) error
	Handler() http.Handler

	SetStatus(int32)
	GetStatus() int32
	GetConfig() store.GameConfig
	GetTeams() []*store.Team
	//GetHub() lab.Hub
	//GetLabByTeam(teamId string) (lab.Lab, bool)
}

type GamePoint struct {
	web        *frontend.WebSite
	store      store.GameConfig
	ipAddrs    []int
	wg         wg.WireguardClient
	dockerHost docker.Host
	closers    []io.Closer
}

func (g GamePoint) Start(ctx context.Context) error {
	panic("implement me")
}

func (g GamePoint) Close() error {
	panic("implement me")
}

func (g GamePoint) Suspend(ctx context.Context) error {
	panic("implement me")
}

func (g GamePoint) Resume(ctx context.Context) error {
	panic("implement me")
}

func (g GamePoint) Finish(s string) {
	panic("implement me")
}

func (g GamePoint) Handler() http.Handler {
	return g.web.Handler()
}

func (g GamePoint) SetStatus(i int32) {
	panic("implement me")
}

func (g GamePoint) GetStatus() int32 {
	panic("implement me")
}

func (g GamePoint) GetConfig() store.GameConfig {
	return g.store
}

func (g GamePoint) GetTeams() []*store.Team {
	panic("implement me")
}

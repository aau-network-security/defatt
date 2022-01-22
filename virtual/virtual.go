// Copyright (c) 2018-2019 Aalborg University
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package virtual

import (
	"context"
	"io"
)

const (
	Running   = State(0)
	Stopped   = State(1)
	Suspended = State(2)
	Error     = State(3)
)

type NatPortSettings struct {
	HostPort    string
	GuestPort   string
	ServiceName string
	Protocol    string
}

type State int

type InstanceInfo struct {
	Image string
	Type  string
	Id    string
	State State
}

type Instance interface {
	Create(context.Context) error
	Start(context.Context) error
	Run(context.Context) error
	Execute(context.Context, []string, string) error
	Suspend(context.Context) error
	Stop() error
	Info() InstanceInfo
	io.Closer
}

type ResourceResizer interface {
	SetRAM(uint) error
	SetCPU(uint) error
}

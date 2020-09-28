package controller

import (
	"testing"

	"github.com/aau-network-security/openvswitch/ovs"
)

func TestOvsManagement_CreateBridge(t *testing.T) {
	tests := []struct {
		name       string
		bridgeName string
		wantErr    bool
	}{
		{name: "No Error, Valid bridge name", bridgeName: "SW", wantErr: false},
		{name: "Error, Invalid bridge name", bridgeName: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &OvsManagement{
				Client:    ovs.New(ovs.Sudo()),
				NetClient: New(Sudo()),
			}
			if err := c.CreateBridge(tt.bridgeName); (err != nil) != tt.wantErr {
				t.Errorf("CreateBridge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

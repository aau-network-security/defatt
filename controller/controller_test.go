package controller

import (
	"testing"
)

var (
	netClient = New(Sudo())
)

func TestOvsManagement_AddBridge(t *testing.T) {
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
			if err := netClient.Ovs.VSwitch.AddBridge(tt.bridgeName); (err != nil) != tt.wantErr {
				t.Errorf("CreateBridge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIPService_AddTuntap(t *testing.T) {
	if err := netClient.IPService.AddTunTap("test", "tap"); err != nil {
		t.Errorf("AddTunTap() error %v", err)
	}
}

func TestIFConfigService_TapUp(t *testing.T) {
	if err := netClient.IFConfig.TapUp("test"); err != nil {
		t.Errorf("TapUp() error %v", err)
	}
}

func TestIFConfigService_TapDown(t *testing.T) {
	if err := netClient.IFConfig.TapDown("test"); err != nil {
		t.Errorf("TapDown() error %v", err)
	}
}

func TestIPService_DelTuntap(t *testing.T) {
	if err := netClient.IPService.DelTuntap("test", "tap"); err != nil {
		t.Errorf("DelTunTap() error %v", err)
	}

}

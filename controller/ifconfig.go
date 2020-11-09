package controller

type IFConfigService struct {
	c *NetController
}

// TapUp ; up given tap
// i.e ifconfig tap0 up
func (ipc *IFConfigService) TapUp(tap string) error {
	cmds := []string{tap, "up"}
	_, err := ipc.exec(cmds...)
	return err
}

// TapDown ; down given tap
// i.e ifconfig tap0 down
func (ipc *IFConfigService) TapDown(tap string) error {
	cmds := []string{tap, "down"}
	_, err := ipc.exec(cmds...)
	return err
}

// OvsBridgeUp configures an internal ip address on the ovs bridge
// ifconfig ovs-br1 192.168.0.1 netmask 255.255.0.0 up
func (ipc *IFConfigService) OvsBridgeUp(bridge, ip, netmask string) error {
	cmds := []string{bridge, ip, "netmask", netmask, "up"}
	_, err := ipc.exec(cmds...)
	return err
}

// exec executes an ExecFunc using 'ip'.
func (ipc *IFConfigService) exec(args ...string) ([]byte, error) {

	return ipc.c.exec("ifconfig", args...)
}

// todo: if required add more functions below

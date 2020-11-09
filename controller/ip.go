package controller

type IPService struct {
	c *NetController
}

// AddTunTap function adds  tap with given mode
// i.e:  ip tuntap add tap0 mode tap
// ip tuntap add <tap name> mode <tap mode>
func (ipc *IPService) AddTunTap(tap, mode string) error {
	cmds := []string{"tuntap", "add", tap, "mode", mode}
	//_, err := ipc.exec(fmt.Sprintf("tuntap add %s mode %s", tap, mode))
	_, err := ipc.exec(cmds...)
	return err
}

// DeleteTuntap deletes tap with given name and mode
func (ipc *IPService) DelTuntap(tap, mode string) error {
	cmds := []string{"tuntap", "del", tap, "mode", mode}
	//_, err := ipc.exec(fmt.Sprintf("tuntap del %s mode %s", tap, mode))
	_, err := ipc.exec(cmds...)

	return err
}

// exec executes an ExecFunc using 'ip'.
func (ipc *IPService) exec(args ...string) ([]byte, error) {
	return ipc.c.exec("ip", args...)
}

// todo: if required add more functions below

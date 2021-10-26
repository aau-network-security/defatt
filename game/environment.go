package game

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/aau-network-security/defatt/controller"
	"github.com/aau-network-security/defatt/dnet/dhcp"
	"github.com/aau-network-security/defatt/dnet/dhcp/proto"
	dhproto "github.com/aau-network-security/defatt/dnet/dhcp/proto"
	"github.com/aau-network-security/defatt/dnet/wg"
	vpn "github.com/aau-network-security/defatt/dnet/wg/proto"
	"github.com/aau-network-security/defatt/store"
	"github.com/aau-network-security/defatt/virtual"
	"github.com/aau-network-security/defatt/virtual/docker"
	"github.com/aau-network-security/defatt/virtual/vbox"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
)

var (
	// there can be only 50 VPN Interface it means 25 Games *(one for blue one for red )
	// this can be changed
	redListenPort    uint = 5181
	blueListenPort   uint = 5182
	min                   = 7900
	max                   = 7950
	gmin                  = 5350
	gmax                  = 5375
	smin                  = 3000
	smax                  = 3500
	rmin                  = 5000
	rmax                  = 5300
	challengeURLList      = map[string]string{
		"ftp":      "registry.gitlab.com/haaukins/forensics/ftp_bf_login",
		"hb":       "registry.gitlab.com/haaukins/web-exploitation/heartbleed",
		"microcms": "registry.gitlab.com/haaukins/web-exploitation/micro_cms",
		"scan":     "registry.gitlab.com/haaukins/forensics/hidden-server",
		"rot":      "registry.gitlab.com/haaukins/crytopgraphy/rot13",
		"csrf":     "registry.gitlab.com/haaukins/web-exploitation/csrf",
		"uwb":      "registry.gitlab.com/haaukins/web-exploitation/webadmin-1.920-urce",
	}
	ErrVMNotCreated       = errors.New("No VM created")
	ErrGettingContainerID = errors.New("Could not get container ID")
)

type environment struct {
	// challenge microservice should be integrated heres
	controller controller.NetController
	wg         vpn.WireguardClient
	dhcp       dhproto.DHCPClient
	dockerHost docker.Host
	closers    []io.Closer
	vlib       vbox.Library
}

type GameConfig struct {
	ID         string
	ScenarioNo int
	Name       string
	Tag        string
	Host       string
	WgConfig   wg.WireGuardConfig
	env        *environment
	NetworksIP map[string]string
	redVPNIp   string
	blueVPNIp  string
	redPort    uint
	bluePort   uint
}

type VPNConfig struct {
	PeerIP           string
	PrivateKeyClient string
	ServerPublicKey  string
	AllowedIPs       string
	Endpoint         string
}

func NewEnvironment(conf *GameConfig, vlib vbox.Library) (*GameConfig, error) {

	netController := controller.New()
	netController.IPPool = controller.NewIPPoolFromHost()

	dockerHost := docker.NewHost()

	env := &environment{
		controller: *netController,
		dockerHost: dockerHost,
		vlib:       vlib,
	}

	conf.env = env

	log.Info().Msgf("New environment initialized ")
	return conf, nil
}

func (env *environment) Close() error {
	//var wg sync.WaitGroup
	// var closers []io.Closer

	// todo: add closers for other components as well
	return nil
}

func (gc *GameConfig) StartGame(ctx context.Context, tag, name string, scenarioNo int) error {

	log.Info().Str("Game Tag", tag).
		Str("Game Name", name).
		Int("Scenario Number", scenarioNo).
		Msg("starting game")

	scenario, err := store.GetScenarioByID(scenarioNo)
	if err != nil {
		return err
	}

	// bridge name will be same with event tag
	bridgeName := tag

	log.Debug().Str("Game", name).Str("bridgeName", tag).Msg("creating openvswitch bridge")
	if err := gc.env.initializeOVSBridge(bridgeName); err != nil {
		return err
	}

	numNetworks := len(scenario.Networks)
	log.Debug().Str("Game", name).Int("Networks", numNetworks).Msg("Creating randomized networks")
	if err := gc.env.createRandomNetworks(bridgeName, numNetworks); err != nil {
		return err
	}

	log.Debug().Str("Game", name).Msg("configuring monitoring")
	if err := gc.env.configureMonitor(ctx, bridgeName, numNetworks); err != nil {
		log.Error().Err(err).Msgf("configuring monitoring")
		return err
	}

	var vlanPorts []string
	for i := 1; i <= numNetworks; i++ {
		vlan := fmt.Sprintf("%s_vlan%d", tag, i*10)
		vlanPorts = append(vlanPorts, vlan)
	}
	vlanPorts = append(vlanPorts, fmt.Sprintf("%s_monitoring", tag))

	log.Debug().Str("Game", bridgeName).Msgf("Initilizing VPN VM")

	//assign connection port to RED users
	redTeamVPNPort := getRandomPort(min, max)

	// assign grpc port to wg vms
	wgPort := getRandomPort(gmin, gmax)

	routerPort := getRandomPort(rmin, rmax)

	//assign connection port to Blue users
	blueTeamVPNPort := getRandomPort(min, max)

	if err := gc.env.initWireguardVM(ctx, tag, vlanPorts, redTeamVPNPort, blueTeamVPNPort, wgPort, routerPort); err != nil {

		return err
	}

	log.Debug().Str("Game", name).Msg("waiting for wireguard vm to boot")

	dhcpClient, err := dhcp.NewDHCPClient(ctx, gc.WgConfig, wgPort)
	if err != nil {
		log.Error().Err(err).Msg("connecting to DHCP service")
		return err
	}

	gc.env.dhcp = dhcpClient

	log.Debug().Str("Game  ", name).Msg("starting DHCP server")
	gc.NetworksIP, err = gc.env.initDHCPServer(ctx, bridgeName, numNetworks)
	if err != nil {
		return err
	}

	wgClient, err := wg.NewGRPCVPNClient(ctx, gc.WgConfig, wgPort)
	if err != nil {
		log.Error().Err(err).Msg("connecting to wireguard service")
		return err
	}
	gc.env.wg = wgClient

	log.Debug().Str("Game", name).Msg("initializing scenario")
	if err := gc.env.initializeScenario(ctx, bridgeName, scenario); err != nil {
		return err
	}

	ethInterfaceName := "eth0" // can be customized later

	redTeamVPNIp, err := gc.env.getRandomIp()
	if err != nil {
		log.Error().Err(err).Msg("Problem in generating red team VPNip")
		return err
	}

	gc.redVPNIp = fmt.Sprintf("%s.0/24", redTeamVPNIp)
	//Assigning a connection port for Red team

	gc.redPort = redTeamVPNPort

	// create wireguard interface for red team
	wgNICred := fmt.Sprintf("%s_red", tag)

	// initializing VPN endpoint for red team
	if err := gc.env.initVPNInterface(gc.redVPNIp, redListenPort, wgNICred, ethInterfaceName); err != nil {
		return err
	}

	blueTeamVPNIp, err := gc.env.getRandomIp()
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	gc.blueVPNIp = fmt.Sprintf("%s.0/24", blueTeamVPNIp)

	//Assigning a connection port for blue team

	gc.bluePort = blueTeamVPNPort
	// initializing VPN endpoint for blue team

	//create wireguard interface for blue team
	wgNICblue := fmt.Sprintf("%s_blue", tag)

	if err := gc.env.initVPNInterface(gc.blueVPNIp, blueListenPort, wgNICblue, ethInterfaceName); err != nil {
		return err
	}

	macAddress := "04:d3:b0:9b:ea:d6"
	macAddressClean := strings.ReplaceAll(macAddress, ":", "")

	log.Debug().Str("game", tag).Msg("Initalizing SoC")
	socPort := getRandomPort(smin, smax)
	ifaces := []string{fmt.Sprintf("%s_monitoring", tag), fmt.Sprintf("%s_AllBlue", tag)}
	if err := gc.env.initializeSOC(ctx, ifaces, macAddressClean, tag, 2, socPort); err != nil {
		log.Error().Err(err).Str("game", tag).Msg("starting SoC vm")
		return err
	}

	log.Info().Str("Game Tag", tag).
		Str("Game Name", name).
		Int("Scenario Number", scenarioNo).
		Msg("started game")

	return nil
}

func (env *environment) initVPNInterface(ipAddress string, port uint, vpnInterfaceName, ethInterface string) error {

	// ipAddress should be in this format : "45.11.23.1/24"
	// port should be unique per interface

	_, err := env.wg.InitializeI(context.Background(), &vpn.IReq{
		Address:            ipAddress,
		ListenPort:         uint32(port),
		SaveConfig:         true,
		Eth:                ethInterface,
		IName:              vpnInterfaceName,
		DownInterfacesFile: "/etc/network/downinterfaces",
	})
	if err != nil {
		log.Error().Msgf("Error in initializing interface %v", err)
		return err
	}
	return nil
}

func (env *environment) createRandomNetworks(bridge string, numberOfNetworks int) error {
	var wg sync.WaitGroup
	for i := 1; i <= numberOfNetworks; i++ {
		wg.Add(1)
		//TODO() Handle the error from this function
		go env.createPort(&wg, bridge, i)
	}
	wg.Wait()
	return nil
}

func (env *environment) createPort(wg *sync.WaitGroup, bridge string, vlan int) error {
	portName := fmt.Sprintf("%s_vlan%d", bridge, vlan*10)
	vlanTag := fmt.Sprintf("%d", vlan*10)
	defer wg.Done()

	if err := env.controller.IPService.AddTunTap(portName, "tap"); err != nil {
		log.Error().Err(err).Str("port", portName).Msg("creating port")
		return err
	}

	if err := env.controller.IFConfig.TapUp(portName); err != nil {
		log.Error().Err(err).Str("port", portName).Msg("setting port as UP")
		return err
	}

	if err := env.controller.Ovs.VSwitch.AddPortTagged(bridge, portName, vlanTag); err != nil {
		log.Error().Err(err).Str("port", portName).Msgf("adding port")
		return err
	}

	return nil
}

func (env *environment) initializeOVSBridge(bridgeName string) error {
	if err := env.controller.Ovs.VSwitch.AddBridge(bridgeName); err != nil {
		log.Error().Err(err).Msg("creating OVS bridge")
		return err
	}
	return nil
}

func (env *environment) attachChallenge(ctx context.Context, wg *sync.WaitGroup, bridge string, challenge, vlan string) error {
	defer wg.Done()

	container := docker.NewContainer(docker.ContainerConfig{
		Image: challengeURLList[challenge],
		Labels: map[string]string{
			"nap": "challenges",
		}})

	if err := container.Create(ctx); err != nil {
		log.Error().Err(err).Str("challenge name", challenge).Msg("creating container")
		return err
	}

	if err := container.Start(ctx); err != nil {
		log.Error().Err(err).Str("challenge name", challenge).Msg("starting container")
		return err
	}

	cid := container.ID()
	if cid == "" {
		log.Error().Str("challenge name", challenge).Msg("getting ID for container")
		return ErrGettingContainerID
	}

	if err := env.controller.Ovs.Docker.AddPort(bridge, "eth0", cid, ovs.DockerOptions{DHCP: true, VlanTag: vlan}); err != nil {
		log.Error().Err(err).Str("container", cid).Msg("adding port to container")
		return err
	}

	return nil
}

func (env *environment) initializeScenario(ctx context.Context, bridge string, scenario store.Scenario) error {
	var wg sync.WaitGroup
	var i int
	for _, vlan := range scenario.Networks {
		i++
		vtag := fmt.Sprintf("%d", i*10)
		log.Info().Str("vlan", vtag).Str("Game tag", bridge).Msg("adding challenges to network")
		for _, chal := range vlan.Chals {
			wg.Add(1)
			//TODO() Handle the error from this function
			go env.attachChallenge(ctx, &wg, bridge, chal, vtag)
		}
	}
	wg.Wait()

	return nil
}

func (env *environment) initDHCPServer(ctx context.Context, bridge string, numberNetworks int) (map[string]string, error) {
	var networks []*dhproto.Network
	var staticHosts []*dhproto.StaticHost
	ipList := make(map[string]string)

	for i := 1; i <= numberNetworks; i++ {
		var network dhproto.Network
		randIP, _ := env.controller.IPPool.Get()
		network.Network = randIP + ".0"
		network.Min = randIP + ".6"
		network.Max = randIP + ".254"
		network.Router = randIP + ".1"

		ipList[fmt.Sprintf("vlan_%d", 10*i)] = randIP + ".0/24"
		networks = append(networks, &network)
	}

	// Setup monitoring network

	monitoringNet := dhproto.Network{
		Network: "10.10.10.0",
		Min:     "10.10.10.6",
		Max:     "10.10.10.254",
		Router:  "10.10.10.1",
	}
	ipList[""] = "10.10.10.0/24"

	networks = append(networks, &monitoringNet)

	host := dhproto.StaticHost{
		Name:       "SOC",
		MacAddress: "04:d3:b0:9b:ea:d6",
		Address:    "10.10.10.200",
	}

	staticHosts = append(staticHosts, &host)

	_, err := env.dhcp.StartDHCP(ctx, &proto.StartReq{Networks: networks, StaticHosts: staticHosts})
	if err != nil {
		return ipList, err
	}

	return ipList, nil
}

//configureMonitor will configure the monitoring VM by attaching the correct interfaces
func (env *environment) configureMonitor(ctx context.Context, bridge string, numberNetworks int) error {

	log.Info().Str("game tag", bridge).Msg("creating monitoring network")
	monitoring := fmt.Sprintf("%s_monitoring", bridge)

	if err := env.controller.IPService.AddTunTap(monitoring, "tap"); err != nil {
		log.Error().Err(err).Str("port", monitoring).Msg("creating port")
		return err
	}

	if err := env.controller.IFConfig.TapUp(monitoring); err != nil {
		log.Error().Err(err).Str("port", monitoring).Msg("setting port as UP")
		return err
	}

	if err := env.controller.Ovs.VSwitch.AddPort(bridge, monitoring); err != nil {
		log.Error().Err(err).Str("port", monitoring).Str("bridge", bridge).Msg("adding port to bridge")
		return err
	}

	mirror := fmt.Sprintf("%s_mirror", bridge)
	AllTraffic := fmt.Sprintf("%s_AllBlue", bridge)

	log.Info().Str("game tag", bridge).Msg("Creating the network mirror")
	if err := env.controller.Ovs.VSwitch.CreateMirrorforBridge(mirror, bridge); err != nil {
		log.Error().Err(err).Msg("creating mirror")
		return err
	}

	if err := env.controller.IPService.AddTunTap(AllTraffic, "tap"); err != nil {
		log.Error().Err(err).Str("port", AllTraffic).Msg("creating port")
		return err
	}

	if err := env.controller.IFConfig.TapUp(AllTraffic); err != nil {
		log.Error().Err(err).Str("port", AllTraffic).Msg("setting port as UP")
		return err
	}

	if err := env.controller.Ovs.VSwitch.AddPort(bridge, AllTraffic); err != nil {
		log.Error().Err(err).Str("port", AllTraffic).Str("bridge", bridge).Msg("adding port to bridge")
		return err
	}

	portUUID, err := env.controller.Ovs.VSwitch.GetPortUUID(AllTraffic)
	if err != nil {
		log.Error().Err(err).Str("port", AllTraffic).Msg("getting port uuid")
		return err
	}

	var vlans []string
	for i := 1; i <= numberNetworks; i++ {
		vlan := fmt.Sprintf("%d", i*10)
		vlans = append(vlans, vlan)
	}

	if err := env.controller.Ovs.VSwitch.MirrorAllVlans(mirror, portUUID, vlans); err != nil {
		log.Error().Err(err).Msgf("mirroring traffic")
		return err
	}

	return nil
}

func (env *environment) initializeSOC(ctx context.Context, networks []string, mac string, tag string, nic int, socPort uint) error {

	//TODO: Add random port here | soc

	//TODO: Solve problem with the soc ovaFile
	vm, err := env.vlib.GetCopy(ctx, tag,
		vbox.InstanceConfig{Image: "soc.ova",
			CPU:      2,
			MemoryMB: 8096},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				HostPort:    strconv.FormatUint(uint64(socPort), 10),
				GuestPort:   "22",
				ServiceName: "sshd",
				Protocol:    "tcp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(networks, false),
		vbox.SetMAC(mac, nic),
	)

	if err != nil {
		log.Error().Err(err).Msg("creating copy of SoC VM")
		return err
	}
	if vm == nil {
		return ErrVMNotCreated
	}
	log.Debug().Str("VM", vm.Info().Id).Msg("starting VM")

	if err := vm.Start(ctx); err != nil {
		log.Error().Err(err).Msgf("starting virtual machine")
		return err
	}
	return nil
}

func (env *environment) initWireguardVM(ctx context.Context, tag string, vlanPorts []string, redTeamVPNport, blueTeamVPNport, wgPort uint, routerPort uint) error {

	//TODO: Add random port here | Router
	//* Just needs to be tested now.
	// socPort := getRandomPort(5000, 6000)

	vm, err := env.vlib.GetCopy(ctx,
		tag,
		vbox.InstanceConfig{Image: "Router.ova",
			CPU:      1,
			MemoryMB: 2048},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				// this is for gRPC service
				HostPort:    strconv.FormatUint(uint64(wgPort), 10),
				GuestPort:   "5353",
				ServiceName: "wgservice",
				Protocol:    "tcp",
			},
			{
				HostPort:    strconv.FormatUint(uint64(redTeamVPNport), 10),
				GuestPort:   strconv.FormatUint(uint64(redListenPort), 10),
				ServiceName: "wgRedConnection",
				Protocol:    "udp",
			},
			{
				HostPort:    strconv.FormatUint(uint64(blueTeamVPNport), 10),
				GuestPort:   strconv.FormatUint(uint64(blueListenPort), 10),
				ServiceName: "wgBlueConnection",
				Protocol:    "udp",
			},

			{
				HostPort:    strconv.FormatUint(uint64(routerPort), 10),
				GuestPort:   "22",
				ServiceName: "sshd",
				Protocol:    "tcp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(vlanPorts, false),
		//vbox.PortForward(min, max), // this is added to enable range of port to be used in Wireguard Interface initializing
	)

	if err != nil {
		log.Error().Err(err).Msg("creating VPN VM")
		return err
	}
	if vm == nil {
		return ErrVMNotCreated
	}
	log.Debug().Str("VM", vm.Info().Id).Msg("starting VM")

	if err := vm.Start(ctx); err != nil {
		log.Error().Err(err).Msgf("starting virtual machine")
		return err
	}
	return nil
}

func (gc *GameConfig) CreateVPNConfig(ctx context.Context, isRed bool, idUser string) (VPNConfig, error) {

	var nicName string

	var allowedIps []string
	var peerIP string
	var endpoint string
	if isRed {
		nicName = fmt.Sprintf("%s_red", gc.Tag)

		for key := range gc.NetworksIP {
			if gc.NetworksIP[key] == "10.10.10.0/24" {
				continue
			}
			allowedIps = append(allowedIps, gc.NetworksIP[key])
			break
		}

		peerIP = gc.redVPNIp
		allowedIps = append(allowedIps, peerIP)

		endpoint = fmt.Sprintf("%s.%s:%d", gc.Tag, gc.Host, gc.redPort)
	} else {

		nicName = fmt.Sprintf("%s_blue", gc.Tag)
		for key := range gc.NetworksIP {
			allowedIps = append(allowedIps, gc.NetworksIP[key])
		}

		peerIP = gc.blueVPNIp
		allowedIps = append(allowedIps, peerIP)
		endpoint = fmt.Sprintf("%s.%s:%d", gc.Tag, gc.Host, gc.bluePort)
	}

	serverPubKey, err := gc.env.wg.GetPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: nicName, PrivKeyName: nicName})
	if err != nil {
		log.Error().Err(err).Str("User", idUser).Msg("Err get public nicName wireguard")
		return VPNConfig{}, err
	}

	_, err = gc.env.wg.GenPrivateKey(ctx, &vpn.PrivKeyReq{PrivateKeyName: gc.Tag + "_" + idUser + "_"})
	if err != nil {
		//fmt.Printf("Err gen private nicName wireguard  %v", err)
		log.Error().Err(err).Str("User", idUser).Msg("Err gen private nicName wireguard")
		return VPNConfig{}, err
	}

	//generate client public nicName
	//log.Info().Msgf("Generating public nicName for team %s", evTag+"_"+team+"_"+strconv.Itoa(ipAddr))
	_, err = gc.env.wg.GenPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: gc.Tag + "_" + idUser + "_", PrivKeyName: gc.Tag + "_" + idUser + "_"})
	if err != nil {
		log.Error().Err(err).Str("User", idUser).Msg("Err gen public nicName client")
		return VPNConfig{}, err
	}

	clientPubKey, err := gc.env.wg.GetPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: gc.Tag + "_" + idUser + "_"})
	if err != nil {
		fmt.Printf("Error on GetPublicKey %v", err)
		return VPNConfig{}, err
	}

	//hitNetworks = "get all networks here"
	//TODO from DAtabase/teamStore or something

	// Todo: get events team length from environment --- //pIP := fmt.Sprintf("%d/32", len(ev.GetTeams())+2)
	pIP := fmt.Sprintf("%d/32", IPcounter())
	//pIP := fmt.Sprintf("%s/32", idUser )
	//todo: Keep track of what IPs are added.s

	peerIP = strings.Replace(peerIP, "0/24", pIP, 1)

	_, err = gc.env.wg.AddPeer(ctx, &vpn.AddPReq{
		Nic:        nicName,
		AllowedIPs: peerIP,
		PublicKey:  clientPubKey.Message,
	})

	if err != nil {
		log.Error().Err(err).Msg("Error on adding peer to interface")
		return VPNConfig{}, err

	}

	clientPrivKey, err := gc.env.wg.GetPrivateKey(ctx, &vpn.PrivKeyReq{PrivateKeyName: gc.Tag + "_" + idUser + "_"})
	if err != nil {
		log.Error().Err(err).Msg("getting priv NIC")
		return VPNConfig{}, err
	}

	return VPNConfig{
		ServerPublicKey:  serverPubKey.Message,
		PrivateKeyClient: clientPrivKey.Message,
		Endpoint:         endpoint,
		AllowedIPs:       strings.Join(allowedIps, ", "),
		PeerIP:           peerIP,
	}, nil

}

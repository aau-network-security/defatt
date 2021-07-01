package game

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	vpn "github.com/aau-network-security/defatt/app/daemon/vpn-proto"
	"github.com/aau-network-security/defatt/config"
	"github.com/aau-network-security/defatt/controller"
	"github.com/aau-network-security/defatt/dnet/dhcp"
	"github.com/aau-network-security/defatt/dnet/wg"

	//TODO:WAITING FOR FRONTEND
	//"github.com/aau-network-security/defat/frontend"
	//"github.com/aau-network-security/defat/store"
	"github.com/aau-network-security/defatt/virtual"
	"github.com/aau-network-security/defatt/virtual/docker"
	"github.com/aau-network-security/defatt/virtual/vbox"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
)

const (
	MAX_NET_CONN = 7
)

var (
	// there can be only 50 VPN Interface it means 25 Games *(one for blue one for red )
	// this can be changed
	min              = 7900
	max              = 7950
	challengeURLList = map[string]string{
		"ftp":      "registry.gitlab.com/haaukins/forensics/ftp_bf_login",
		"hb":       "registry.gitlab.com/haaukins/web-exploitation/heartbleed",
		"microcms": "registry.gitlab.com/haaukins/web-exploitation/micro_cms",
		"scan":     "registry.gitlab.com/haaukins/forensics/hidden-server",
		"rot":      "registry.gitlab.com/haaukins/crytopgraphy/rot13",
		"csrf":     "registry.gitlab.com/haaukins/web-exploitation/csrf",
		"uwb":      "registry.gitlab.com/haaukins/web-exploitation/webadmin-1.920-urce",
	}
	TemporaryScenariosPlaceHolder = map[int]Scenario{
		1: {
			ID: 1,
			Networks: []network{
				{
					Chals: []string{"hb", "ftp", "scan"},
					Vlan:  "vlan20",
				},
				{
					Chals: []string{"scan", "csrf"},
					Vlan:  "vlan30",
				},
				{
					Chals: []string{"rot", "uwb"},
					Vlan:  "vlan10",
				},
			},
			Story:    "Scenario 1 Story placeholder",
			Duration: "2",

			Difficulty: "Easy",
		},
		2: {
			ID: 2,
			Networks: []network{
				{
					Chals: []string{"microcms", "joomla", "uwb"},
					Vlan:  "vlan10",
				},
				{
					Chals: []string{"jwt", "csrf"},
					Vlan:  "vlan20",
				},
				{
					Chals: []string{"rot", "uwb"},
					Vlan:  "vlan40",
				},
				{
					Chals: []string{"rot", "uwb"},
					Vlan:  "vlan30",
				},
			},
			Story:      "Scenario 2 Story placeholder",
			Duration:   "3",
			Difficulty: "Moderate",
		},
	}
)

type Environment interface {
	GetScenarios() string
}

type network struct {
	Vlan  string
	Chals []string
}

type Scenario struct {
	ID         int
	Networks   []network
	Story      string
	Duration   string
	Difficulty string
}

type environment struct {
	// web interface microservice should stay here
	// challenge microservice should be integrated heres
	controller controller.NetController
	wg         vpn.WireguardClient
	dockerHost docker.Host
	closers    []io.Closer
	config     GameConfig
	vlib       vbox.Library
	dhcp       *dhcp.Server
	//web        *frontend.WebSite
}

type GameConfig struct {
	ID         string
	ScenarioNo int
	Name       string
	Tag        string
	WgConfig   wg.WireGuardConfig
}

func NewEnvironment(conf GameConfig, vboxConf config.VmConfig) (*environment, error) {
	//if len(conf.Scenario.Networks) > MAX_NET_CONN {
	//	return nil, fmt.Errorf("exceeds maximum number of Networks for a environment. Max is %d", MAX_NET_CONN)
	//}

	wgClient, err := wg.NewGRPCVPNClient(conf.WgConfig)
	if err != nil {
		log.Error().Msgf("Connection error on wireguard service error %v ", err)
		return nil, err
	}
	netController := controller.New()
	vlib := vbox.NewLibrary(vboxConf.OvaDir)
	if vlib == nil {
		log.Error().Msgf("Library could not be created properly...")
		return nil, fmt.Errorf("Error on new library")
	}

	//TODO:WAITING FOR FRONTEND
	//webUI := frontend.NewFrontend(store.GameConfig{
	//	Name:       conf.Name,
	//	Tag:        conf.Tag,
	//	ScenarioID: conf.ScenarioNo,
	//	//StartedAt:  nil,
	//	//FinishedAt: nil,
	//}, wgClient)

	dockerHost := docker.NewHost()
	env := &environment{
		controller: *netController,
		wg:         wgClient,
		dockerHost: dockerHost,
		config:     conf,
		vlib:       vlib,
		//TODO:WAITING FOR FRONTEND
		//web:        webUI,
	}
	log.Info().Msgf("New environment initialized ")
	return env, nil
}

func (g *environment) Close() error {
	//var wg sync.WaitGroup
	var closers []io.Closer

	if g.dhcp != nil {
		closers = append(closers, g.dhcp)
	}
	// todo: add closers for other components as well
	return nil
}

func (g *environment) StartGame(tag, name string, scenarioNo int) error {
	// red team wireguard interface port is : 87878

	log.Info().Str("GamePoint Tag", tag).
		Str("GamePoint Name", name).
		Int("Scenario Number", scenarioNo).
		Msgf("Staring GamePoint")
	// bridge name will be same with event tag
	bridgeName := tag
	var mainErr error
	// todo: scenarios should be taken from a resource such as postresql / mongodb !
	selectedScenario := TemporaryScenariosPlaceHolder[scenarioNo]

	numNetworks := len(selectedScenario.Networks)
	log.Info().Msgf("Setting openvswitch bridge %s", bridgeName)

	if err := g.initializeOVSBridge(bridgeName); err != nil {
		mainErr = err
	}

	if err := g.createRandomNetworks(bridgeName, numNetworks); err != nil {
		mainErr = err
	}

	if err := g.configureMonitor(bridgeName, numNetworks); err != nil {
		log.Error().Err(err).Msgf("Error to configure monitoring")
		return nil
	}

	if err := g.initializeScenarios(bridgeName, &g.controller, scenarioNo); err != nil {
		mainErr = err
	}

	ethInterfaceName := "eth0" // can be customized later

	redTeamVPNIp, err := g.getRandomIp()
	if err != nil {
		mainErr = err
	}
	redTeamVPNIp = fmt.Sprintf("%s.0/24", redTeamVPNIp)
	// initializing VPN endpoint for red team
	if err := g.initVPNInterface(redTeamVPNIp, getRandomPort(), fmt.Sprintf("%s_red", tag), ethInterfaceName); err != nil {
		mainErr = err
	}
	blueTeamVPNIp, err := g.getRandomIp()
	if err != nil {
		mainErr = err
	}
	blueTeamVPNIp = fmt.Sprintf("%s.0/24", blueTeamVPNIp)
	// initializing VPN endpoint for blue team
	if err := g.initVPNInterface(blueTeamVPNIp, getRandomPort(), fmt.Sprintf("%s_blue", tag), ethInterfaceName); err != nil {
		mainErr = err
	}

	return mainErr

}

func (g *environment) getRandomIp() (string, error) {
	var ip string
	if g.controller.IPPool != nil {
		ipAddress, err := g.controller.IPPool.Get()
		if err != nil {
			return "", err
		}
		ip = ipAddress
	}
	return ip, nil
}

func getRandomPort() uint {
	port := rand.Intn(max-min) + min
	for checkPort(port) {
		port = rand.Intn(max-min) + min
	}
	return uint(port)
}

func (g *environment) initVPNInterface(ipAddress string, port uint, vpnInterfaceName, ethInterface string) error {

	// ipAddress should be in this format : "45.11.23.1/24"
	// port should be unique per interface

	_, err := g.wg.InitializeI(context.Background(), &vpn.IReq{
		Address:    ipAddress,
		ListenPort: uint32(port),
		SaveConfig: true,
		Eth:        ethInterface,
		IName:      vpnInterfaceName,
	})
	if err != nil {
		log.Error().Msgf("Error in initializing interface %v", err)
		return err
	}
	return nil
}

//TODO:WAITING FOR FRONTEND
//func (g *environment) GetFrontend() *frontend.WebSite {
//	return g.web
//}

func (g *environment) createRandomNetworks(bridge string, numberOfNetworks int) error {
	vlanTags := make(map[string]string)
	var waitGroups sync.WaitGroup
	log.Info().Msgf("Creating randomized Networks for chosen number of Networks %d", numberOfNetworks)
	var mainErr error
	for i := 1; i < numberOfNetworks+1; i++ {
		waitGroups.Add(1)
		vlan := fmt.Sprintf("vlan%d", i*10)
		vlanTags[vlan] = fmt.Sprintf("%d", i*10)
		go func() {
			defer waitGroups.Done()
			if err := g.controller.Ovs.VSwitch.AddPortTagged(bridge, vlan, fmt.Sprintf("%d", i*10)); err != nil {
				log.Error().Msgf("Error on adding port with tag err %v", err)
				mainErr = err
			}
			log.Info().Msgf("AddPort Set Interface Options %s", vlan)
			if err := g.controller.Ovs.VSwitch.Set.Interface(vlan, ovs.InterfaceOptions{Type: ovs.InterfaceTypeInternal}); err != nil {
				log.Error().Msgf("Error on matching interface error %v", err)
				mainErr = err
			}

			//ip tuntap add tap0 mode tap
			//ifconfig tap0 up
			//ip tuntap add tap2 mode tap
			// ifconfig tap2 up
			//ip tuntap add tap4 mode tap
			//ifconfig tap4 up
			t := fmt.Sprintf("tap%d", i)
			if err := g.controller.IPService.AddTunTap(t, "tap"); err != nil {
				log.Error().Msgf("Error happened on adding tuntap %v", err)
				mainErr = err
			}
			if err := g.controller.IFConfig.TapUp(t); err != nil {
				log.Error().Msgf("Error happened on making up tap %s %v", t, err)
				mainErr = err
			}

			tag := fmt.Sprintf("%d", i*10)
			//ovs-vsctl add-port SW tap0 tag=10
			//ovs-vsctl add-port SW tap2 tag=20
			//ovs-vsctl add-port SW tap4 tag=30
			if err := g.controller.Ovs.VSwitch.AddPortTagged(bridge, t, tag); err != nil {
				log.Error().Msgf("Error on adding port with tag err %v", err)
				mainErr = err
			}

			if err := g.controller.IFConfig.TapUp(vlan); err != nil {
				log.Error().Msgf("Error happened on making up tap %s %v", vlan, err)
				mainErr = err
			}
		}()
		waitGroups.Wait()
	}

	log.Info().Msgf("Creating the monitoring network")
	//Always creating +1 network for the monitoring machine.

	//TODO: Make assign the monitoring network smarter ! Now is hardcoded.

	//How it is happening now will be a problem for multiple games
	i := 1

	monitor := fmt.Sprintf("mon%d", i*10)

	if err := g.controller.Ovs.VSwitch.AddPort(bridge, monitor); err != nil {
		log.Error().Msgf("Error on adding port with tag err %v", err)
		return err
	}

	m := fmt.Sprintf("mon%d", i*10)
	if err := g.controller.IPService.AddTunTap(m, "tap"); err != nil {
		log.Error().Msgf("Error happened on adding monitor tuntap %v", err)
		return err
	}
	if err := g.controller.IFConfig.TapUp(m); err != nil {
		log.Error().Msgf("Error happened on making up monitor %s %v", m, err)
		return err
	}
	//adding the monitoring port in the networks
	vlanTags["monitor"] = ""

	server, err := dhcp.New(context.TODO(), vlanTags, bridge, &g.controller)
	if err != nil {
		log.Error().Msgf("Error creating DHCP server %v", err)
		return err
	}
	if err := server.Run(context.Background()); err != nil {
		log.Error().Msgf("Error in starting DHCP  %v", err)
		return err
	}
	g.dhcp = server

	return mainErr
}

func (g *environment) initializeOVSBridge(bridgeName string) error {
	log.Info().Msgf("Game brigde name is set to game tag %s", bridgeName)
	if err := g.controller.Ovs.VSwitch.AddBridge(bridgeName); err != nil {
		log.Error().Msgf("Error on creating OVS bridge %v", err)
		return err
	}
	return nil
}

func (g *environment) attachChallenge(bridge string, challengeList []string, cli *controller.NetController, vlan string) error {
	ctx := context.Background()
	log.Info().Msgf("Starting challenges for the game %s", bridge)
	for _, ch := range challengeList {
		container := docker.NewContainer(docker.ContainerConfig{
			Image: challengeURLList[ch],
			Labels: map[string]string{
				"nap": "challenges",
			}})
		if err := container.Create(ctx); err != nil {
			log.Error().Msgf("Error in creating container  %v", err)
			return err
		}
		if err := container.Start(ctx); err != nil {
			log.Error().Msgf("Error in creating container  %v", err)
			return err
		}

		cid := container.ID()
		if cid == "" {
			return fmt.Errorf("Container ID could be fetched correctly")
		}

		if err := cli.Ovs.Docker.AddPort(bridge, "eth0", cid, ovs.DockerOptions{DHCP: true, VlanTag: vlan}); err != nil {
			log.Error().Msgf("Error on adding port on docker %v", err)
			return err
		}

		if err := cli.Ovs.Docker.SetVlan(bridge, "eth0", cid, vlan); err != nil {
			log.Error().Msgf("Error on ovs-docker SetVlan %v", err)
			return err
		}

	}

	return nil

}

func (g *environment) initializeScenarios(bridge string, cli *controller.NetController, scenarioNumber int) error {
	log.Debug().Msgf("Inializing scenarios for game [ %s ]", bridge)
	networks := TemporaryScenariosPlaceHolder[scenarioNumber].Networks
	var vlans []string
	var initScenErr error
	var waitGroup sync.WaitGroup
	if scenarioNumber > 3 || scenarioNumber < 0 {
		return fmt.Errorf("Invalid senario selection, make a selection between 1 to 3 ")
	}
	for _, net := range networks {
		vlans = append(vlans, net.Vlan)

	}
	log.Debug().Strs("Network Vlans", vlans).Msgf("Vlans")
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := g.initWireguardVM(vlans, min, max); err != nil {
			initScenErr = err
		}
	}()
	waitGroup.Wait()
	//Todo:Fix Mac address problem
	//// initializing SOC all networks
	//if err := g.initializeSOC(vlans); err != nil {
	//	return err
	//}

	// initializing scenarios by attaching correct challenge to correct network
	for _, net := range networks {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if err := g.attachChallenge(bridge, net.Chals, cli, net.Vlan[len(net.Vlan)-2:]); err != nil {
				fmt.Printf("Error in attach challenge %v", err)
				initScenErr = err
			}
		}()
		waitGroup.Wait()
	}

	return initScenErr
}

//configureMonitor will configure the monitoring VM by attaching the correct interfaces
func (g *environment) configureMonitor(bridge string, numberNetworks int) error {

	var ifaces []string
	var vlanTags []string
	var getBlue string  // mirrorName
	var bluePort string // port in OVS for mirror traffic

	getBlue = "blueMirror"
	if err := g.controller.Ovs.VSwitch.CreateMirrorforBridge(getBlue, bridge); err != nil {
		log.Error().Err(err).Msgf("Error on creating mirror")
		return err

	}

	for i := 1; i <= numberNetworks; i++ {
		tag := fmt.Sprintf("%d", i*10)
		vlanTags = append(vlanTags, tag)

	}

	bluePort = "ALLblue"

	if err := g.controller.IPService.AddTunTap(bluePort, "tap"); err != nil {
		log.Error().Msgf("Error happened on adding monitor tuntap %v", err)
		return err
	}
	if err := g.controller.IFConfig.TapUp(bluePort); err != nil {
		log.Error().Msgf("Error happened on making up monitor %s %v", bluePort, err)
		return err
	}

	if err := g.controller.Ovs.VSwitch.AddPort(bridge, bluePort); err != nil {
		log.Error().Err(err).Msgf("Error on adding port to mirror traffic, err %v", err)
		return err
	}
	//
	//log.Info().Msgf("AddPort for mirroring Set Interface Options %s", bluePort)
	//if err := g.controller.Ovs.VSwitch.Set.Interface(bluePort, ovs.InterfaceOptions{Type: ovs.InterfaceTypeInternal}); err != nil {
	//	log.Error().Msgf("Error on matching interface error %v", err)
	//	return err
	//}

	portUUID, err := g.controller.Ovs.VSwitch.GetPortUUID(bluePort)
	if err != nil {
		log.Error().Err(err).Msgf("Error on getting port uuid")
		return err
	}

	if err := g.controller.Ovs.VSwitch.MirrorAllVlans(getBlue, portUUID, vlanTags); err != nil {
		log.Error().Err(err).Msgf("Error on adding port to mirror traffic")
		return err

	}

	//if err := g.controller.Ovs.VSwitch.MirrorAllVlans(getBlue, bluePort, vlanTags); err != nil {
	//	log.Error().Err(err).Msgf("Error on adding port to mirror traffic")
	//	return err
	//
	//}

	//err, monitoringNetwr
	ifaces = append(ifaces, bluePort)

	ineti, err := net.Interfaces()
	if err != nil {
		log.Error().Err(err).Msgf("Error getting the system interfaces")
		panic(err)

	}

	for _, inter := range ineti {
		if strings.Contains(inter.Name, "mon") {
			ifaces = append(ifaces, inter.Name)
			if len(ifaces) != 2 {
				log.Error().Err(err).Msgf("error on creating the list of interfaces")

			}

		}
		continue

	}

	macAddress := g.dhcp.GetMAC()
	macAddressClean := strings.ReplaceAll(macAddress, ":", "")
	nicNumber := len(ifaces) + 1

	fmt.Println(macAddressClean)
	fmt.Println(nicNumber)
	//
	fmt.Println(ifaces)
	if err := g.initializeSOC(ifaces, macAddressClean, nicNumber); err != nil {
		log.Error().Err(err).Msgf("error starting VM with given interfaces")
		return err
	}

	return nil

}

func (env *environment) initWireguardVM(networks []string, min, max int) error {
	log.Debug().Str("Service Port", "5353").Str("VPN Port", "51820").Msgf("Initalizing VPN endpoint for the game")
	vm, err := env.vlib.GetCopy(context.Background(),
		vbox.InstanceConfig{Image: "ubuntu.ova",
			CPU:      1,
			MemoryMB: 2048},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				// this is for gRPC service
				HostPort:    "5353",
				GuestPort:   "5353",
				ServiceName: "wgservice",
				Protocol:    "tcp",
			},
			{
				HostPort:    "5555",
				GuestPort:   "22",
				ServiceName: "sshd",
				Protocol:    "tcp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(networks, false),
		vbox.PortForward(min, max), // this is added to enable range of port to be used in Wireguard Interface initializing
		//vbox.SetNameofVM(),
	)

	if err != nil {
		log.Error().Msgf("Error while getting copy of VM err : %v", err)
		return err
	}
	if vm != nil {

		log.Debug().Msgf("VM [ %s ] is starting .... ", vm.Info().Id)

		if err := vm.Start(context.Background()); err != nil {
			log.Error().Msgf("Failed to start virtual machine on Vlan ")
			return err
		}
	}
	return nil
}

func (env *environment) initializeSOC(networks []string, mac string, nic int) error {

	log.Debug().Str("Elastic Port", "9200").
		Str("Kibana Port", "5601").
		Msgf("Initalizing SoC for the game")

	// todo: Solve problem with the soc ovaFile
	vm, err := env.vlib.GetCopy(context.Background(),
		vbox.InstanceConfig{Image: "soc.ova",
			CPU:      2,
			MemoryMB: 8096},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				HostPort:    "3334",
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
		log.Error().Msgf("Error while getting copy of VM err : %v", err)
		return err
	}
	if vm != nil {
		log.Debug().Msgf("VM [ %s ] is starting .... ", vm.Info().Id)

		if err := vm.Start(context.Background()); err != nil {
			log.Error().Msgf("Failed to start virtual machine on Vlan ")
			return err
		}
	}
	return nil
}

func checkPort(port int) bool {
	portAllocated := fmt.Sprintf(":%d", port)
	// ensure that VPN port is free to allocate
	conn, _ := net.DialTimeout("tcp", portAllocated, time.Second)
	if conn != nil {
		_ = conn.Close()
		fmt.Printf("Checking VPN port %s\n", portAllocated)
		// true means port is already allocated
		return true
	}
	return false
}

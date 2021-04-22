package game

import (
	"context"
	"fmt"
	"io"

	vpn "github.com/aau-network-security/defat/app/daemon/vpn-proto"
	"github.com/aau-network-security/defat/config"
	"github.com/aau-network-security/defat/controller"
	"github.com/aau-network-security/defat/dnet/dhcp"
	"github.com/aau-network-security/defat/dnet/wg"
	"github.com/aau-network-security/defat/virtual"
	"github.com/aau-network-security/defat/virtual/docker"
	"github.com/aau-network-security/defat/virtual/vbox"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/rs/zerolog/log"
)

const (
	MAX_NET_CONN = 7
)

var (
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
}

type GameConfig struct {
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
	dockerHost := docker.NewHost()
	env := &environment{
		controller: *netController,
		wg:         wgClient,
		dockerHost: dockerHost,
		config:     conf,
		vlib:       vlib,
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
	log.Info().Str("Game Tag", tag).
		Str("Game Name", name).
		Int("Scenario Number", scenarioNo).
		Msgf("Staring game")
	// bridge name will be same with event tag
	bridgeName := tag
	selectedScenario := TemporaryScenariosPlaceHolder[scenarioNo]
	numNetworks := len(selectedScenario.Networks)
	log.Info().Msgf("Setting openvswitch bridge %s", bridgeName)
	if err := g.initializeOVSBridge(bridgeName); err != nil {
		return err
	}
	if err := g.createRandomNetworks(bridgeName, int(numNetworks)); err != nil {
		return err
	}

	if err := g.initializeScenarios(bridgeName, &g.controller, scenarioNo); err != nil {
		return err
	}
	return nil

}

func (g *environment) createRandomNetworks(bridge string, numberOfNetworks int) error {
	vlanTags := make(map[string]string)

	log.Info().Msgf("Creating randomized Networks for chosen number of Networks %d", numberOfNetworks)
	for i := 1; i < numberOfNetworks+1; i++ {
		vlan := fmt.Sprintf("vlan%d", i*10)
		vlanTags[vlan] = fmt.Sprintf("%d", i*10)
		if err := g.controller.Ovs.VSwitch.AddPortTagged(bridge, vlan, fmt.Sprintf("%d", i*10)); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
			return err
		}
		log.Info().Msgf("AddPort Set Interface Options %s", vlan)
		if err := g.controller.Ovs.VSwitch.Set.Interface(vlan, ovs.InterfaceOptions{Type: ovs.InterfaceTypeInternal}); err != nil {
			log.Error().Msgf("Error on matching interface error %v", err)
			return err
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
			return err
		}
		if err := g.controller.IFConfig.TapUp(t); err != nil {
			log.Error().Msgf("Error happened on making up tap %s %v", t, err)
			return err
		}

		tag := fmt.Sprintf("%d", i*10)
		//ovs-vsctl add-port SW tap0 tag=10
		//ovs-vsctl add-port SW tap2 tag=20
		//ovs-vsctl add-port SW tap4 tag=30
		if err := g.controller.Ovs.VSwitch.AddPortTagged(bridge, t, tag); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
			return err
		}

		if err := g.controller.IFConfig.TapUp(vlan); err != nil {
			log.Error().Msgf("Error happened on making up tap %s %v", vlan, err)
			return err
		}
	}

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

	return nil
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
			},})
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
	if scenarioNumber > 3 || scenarioNumber < 0 {
		return fmt.Errorf("Invalid senario selection, make a selection between 1 to 3 ")
	}
	for _, net := range networks {
		vlans = append(vlans, net.Vlan)

	}
	log.Debug().Strs("Network Vlans", vlans).Msgf("Vlans")
	if err := g.initializeWireguard(vlans); err != nil {
		return err
	}
	// initializing scenarios by attaching correct challenge to correct network
	for _, net := range networks {
		if err := g.attachChallenge(bridge, net.Chals, cli, net.Vlan[len(net.Vlan)-2:]); err != nil {
			fmt.Printf("Error in attach challenge %v", err)
			return err
		}

	}

	return nil
}

func (env *environment) initializeWireguard(networks []string) error {
	log.Debug().Str("Service Port", "5353").Str("VPN Port", "51820").Msgf("Initalizing VPN endpoint for the game")
	vm, err := env.vlib.GetCopy(context.Background(),
		vbox.InstanceConfig{Image: "ubuntu.ova",
			CPU:      1,
			MemoryMB: 2048},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				HostPort:    "9200",
				GuestPort:   "9200",
				ServiceName: "elasticsearch",
				Protocol:    "tcp",
			},
			{
				HostPort:    "5601",
				GuestPort:   "5601",
				ServiceName: "kibana",
				Protocol:    "tcp",
			},
			{
				// this is for gRPC service
				HostPort:    "5353",
				GuestPort:   "5353",
				ServiceName: "wgservice",
				Protocol:    "tcp",
			},
			{
				// this is for VPN Connection
				HostPort:    "51820",
				GuestPort:   "51820",
				ServiceName: "wireguard",
				Protocol:    "udp",
			},
			{
				HostPort:    "2222",
				GuestPort:   "22",
				ServiceName: "sshd",
				Protocol:    "tcp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(networks, false),
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


func (env *environment) initializeSOC(networks []string) error {
	log.Debug().Str("Elastic Port", "9200").
		Str("Kibana Port", "5601").
		Msgf("Initalizing SoC for the game")
	vm, err := env.vlib.GetCopy(context.Background(),
		vbox.InstanceConfig{Image: "soc.ova",
			CPU:      2,
			MemoryMB: 8096},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				HostPort:    "2222",
				GuestPort:   "22",
				ServiceName: "sshd",
				Protocol:    "tcp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(networks, false),

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
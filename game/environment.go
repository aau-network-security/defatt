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

	vpn "github.com/aau-network-security/defat/app/daemon/vpn-proto"
	"github.com/aau-network-security/defat/config"
	"github.com/aau-network-security/defat/controller"
	"github.com/aau-network-security/defat/dnet/dhcp"
	"github.com/aau-network-security/defat/dnet/wg"
	//TODO:WAITING FOR FRONTEND
	//"github.com/aau-network-security/defat/frontend"
	//"github.com/aau-network-security/defat/store"
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
	ScenarioNo int
	Name       string
	Tag        string
	WgConfig   wg.WireGuardConfig
	Networks   map[string]string
	redVPNIp   string
	blueVPNIp  string
	redPort    uint
	bluePort   uint
}

type VPNConfig struct {
	peerIP           string
	PrivateKeyClient string
	ServerPublicKey  string
	AllowedIps       string
	Endpoint         string
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
	g.config.redVPNIp = redTeamVPNIp
	//Assigning a connection port for Red team
	redTeamVPNPort := getRandomPort()
	g.config.redPort = redTeamVPNPort

	//create wireguard interface for red team
	wgNICred := fmt.Sprintf("%s_red", tag)

	// initializing VPN endpoint for red team
	if err := g.initVPNInterface(redTeamVPNIp, redTeamVPNPort, wgNICred, ethInterfaceName); err != nil {
		mainErr = err
	}

	blueTeamVPNIp, err := g.getRandomIp()
	if err != nil {
		log.Error().Err(err).Msg("")
		mainErr = err
	}

	blueTeamVPNIp = fmt.Sprintf("%s.0/24", blueTeamVPNIp)
	g.config.blueVPNIp = blueTeamVPNIp

	//Assigning a connection port for blue team
	blueTeamVPNPort := getRandomPort()
	g.config.bluePort = blueTeamVPNPort
	// initializing VPN endpoint for blue team

	//create wireguard interface for blue team
	wgNICblue := fmt.Sprintf("%s_blue", tag)

	if err := g.initVPNInterface(blueTeamVPNIp, blueTeamVPNPort, wgNICblue, ethInterfaceName); err != nil {
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

func (g *environment) CreateVPNConfig(ctx context.Context, isRed bool, gameTag string, idUser string) (VPNConfig, error) {

	var nicName string

	var allowedIps string

	var endpoint string
	var hitNetworks string

	if isRed {
		nicName = fmt.Sprintf("%s_red", gameTag)

		for key, _ := range g.config.Networks {
			hitNetworks = g.config.Networks[key]
			allowedIps = fmt.Sprintf("%s, %s", hitNetworks, g.config.redVPNIp)
			continue
		}

		allowedIps = g.config.redVPNIp
		endpoint = fmt.Sprintf("%s.defatt.haaukins.com:%d", gameTag, g.config.redPort)
	} else {

		nicName = fmt.Sprintf("%s_blue", gameTag)
		allowedIps = g.config.blueVPNIp
		endpoint = fmt.Sprintf("%s.defatt.haaukins.com:%d", gameTag, g.config.bluePort)

		//	10.20.30.
	}

	serverPubKey, err := g.wg.GetPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: nicName, PrivKeyName: nicName})
	if err != nil {
		log.Error().Err(err).Str("User", idUser).Msg("Err get public nicName wireguard")
		return VPNConfig{}, err
	}

	_, err = g.wg.GenPrivateKey(ctx, &vpn.PrivKeyReq{PrivateKeyName: gameTag + "_" + idUser + "_"})
	if err != nil {
		//fmt.Printf("Err gen private nicName wireguard  %v", err)
		log.Error().Err(err).Str("User", idUser).Msg("Err gen private nicName wireguard")
		return VPNConfig{}, err
	}

	//generate client public nicName
	//log.Info().Msgf("Generating public nicName for team %s", evTag+"_"+team+"_"+strconv.Itoa(ipAddr))
	_, err = g.wg.GenPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: gameTag + "_" + idUser + "_", PrivKeyName: gameTag + "_" + idUser + "_"})
	if err != nil {
		log.Error().Err(err).Str("User", idUser).Msg("Err gen public nicName client")
		return VPNConfig{}, err
	}
	// get client public nicName
	//log.Info().Msgf("Retrieving public nicName for client %s", idUser)
	clientPubKey, err := g.wg.GetPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: gameTag + "_" + idUser + "_"})
	if err != nil {
		fmt.Printf("Error on GetPublicKey %v", err)
		return VPNConfig{}, err
	}

	//hitNetworks = "get all networks here"
	//TODO from DAtabase/teamStore or something
	pIP := fmt.Sprintf("%d/32", 3)

	//todo: Keep track of what IPs are added.

	peerIP := strings.Replace(allowedIps, "0/24", pIP, 1)
	//log.Info().Str("NIC", evTag).
	//	Str("AllowedIPs", peerIP).
	//	Str("PublicKey ", resp.Message).Msgf("Generating ip address for peer %s, ip address of peer is %s ", team, peerIP)
	addPeerResp, err := g.wg.AddPeer(ctx, &vpn.AddPReq{
		Nic:        nicName,
		AllowedIPs: peerIP, // Todo: get events team length from environment --- //pIP := fmt.Sprintf("%d/32", len(ev.GetTeams())+2)
		PublicKey:  clientPubKey.Message,
	})

	if err != nil {
		fmt.Sprintf("Error on adding peer to interface %v\n", err)
		log.Error().Err(err).Msg("Error on adding peer to interface")
		return VPNConfig{}, err

	}

	fmt.Printf("AddPEER RESPONSE:  %s", addPeerResp.Message)

	clientPrivKey, err := g.wg.GetPrivateKey(ctx, &vpn.PrivKeyReq{PrivateKeyName: gameTag + "_" + idUser + "_"})
	if err != nil {
		fmt.Sprintf("Error on getting priv nicName for team  %v\n", err)
		log.Error().Err(err).Msg("Error on getting priv nicName for team")
		return VPNConfig{}, err
	}
	//log.Info().Msgf("Privatee nicName for team %s is %s ", team, teamPrivKey.Message)
	//log.Info().Msgf("Client configuration is created for server %s", endpoint)
	// creating client configuration file
	// fmt.Sprintf("%s/24", "10.4.2.1") > this should be the lab subnet, necessry subnet which is assigned to team as a lab when they signed up...
	// 87878 > value should be changed with the randomized port where is it created before initializing the interface of wireguard...
	// fmt.Sprintf("%s.defatt.haaukins.com:%d", f.globalInfo.GameTag, 87878) > the dns address of host should be taken from configuration file of defat.

	//		clientConfig := fmt.Sprintf(
	//			`[Interface]
	//Address = %s
	//PrivateKey = %s
	//DNS = 1.1.1.1
	//MTU = 1500
	//[Peer]
	//PublicKey = %s
	//AllowedIps = %s
	//Endpoint =  %s
	//PersistentKeepalive = 25
	//`, allowedIps, clientPrivKey.Message, serverPubKey.Message, fmt.Sprintf("%s/24", "10.4.2.1"), fmt.Sprintf("%s.defatt.haaukins.com:%d", gameTag, g.config.redPort))
	//
	//
	//
	//
	//
	//}

	//Blue Team

	return VPNConfig{
		ServerPublicKey:  serverPubKey.Message,
		PrivateKeyClient: clientPrivKey.Message,
		Endpoint:         endpoint,
		AllowedIps:       allowedIps,
		peerIP:           peerIP,
	}, nil

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
	//get the IPS of the game

	for key, _ := range vlanTags {
		//fmt.Println("Key:", key, "=>", "Element:", value)
		g.config.Networks[key] = server.GetVlanIP(key)
	}
	//g.config.Networks = server
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
		log.Error().Err(err).Msg("error starting VM with given interfaces")
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

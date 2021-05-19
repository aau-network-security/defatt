package frontend

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	wg "github.com/aau-network-security/defat/app/daemon/vpn-proto"
	"github.com/aau-network-security/defat/store"
)

var (
	wd, _      = os.Getwd()
	signingKey = "random"
)

type siteInfo struct {
	GameName string
	Content  string
	Team     store.Team
	GameTag  string
}

type WebSite struct {
	maxReadBytes int64
	signingKey   []byte
	cookieTTL    int
	globalInfo   siteInfo
	TeamStore    store.TeamStore
	wgClient     wg.WireguardClient
}
type SiteOpt func(*WebSite)

func WithMaxReadBytes(b int64) SiteOpt {
	return func(am *WebSite) {
		am.maxReadBytes = b
	}
}

func WithGameName(gameName string) SiteOpt {
	return func(am *WebSite) {
		am.globalInfo.GameName = gameName
	}
}

func NewFrontend(config store.GameConfig, wgClient wg.WireguardClient, opts ...SiteOpt) *WebSite {

	frntend := &WebSite{
		maxReadBytes: 1024 * 1024,
		signingKey:   []byte(signingKey),
		cookieTTL:    int((3 * 24 * time.Hour).Seconds()),
		globalInfo: siteInfo{
			GameName: config.Name,
			GameTag:  config.Tag,
		},
		wgClient: wgClient,
	}
	for _, opt := range opts {
		opt(frntend)
	}

	return frntend
}

func (f *WebSite) Handler() http.Handler {
	h := http.NewServeMux()
	h.HandleFunc("/", f.handleIndex())
	h.HandleFunc("/signup", f.handleSignup())
	h.HandleFunc("/logout", f.handleLogout())
	h.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir(wd+"/frontend/public"))))

	return h
}

func (f *WebSite) handleIndex() http.HandlerFunc {
	indexTemplate := wd + "/frontend/private/index.tmpl.html"
	tmpl, err := parseTemplates(indexTemplate)
	if err != nil {
		log.Println("error index tmpl: ", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data := f.globalInfo
		if err := tmpl.Execute(w, data); err != nil {
			log.Println("template err index: ", err)
		}
	}
}
func (f *WebSite) handleSignup() http.HandlerFunc {
	get := f.handleSignupGet()
	post := f.handleSignupPOST()
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get(w, r)
			return

		case http.MethodPost:
			post(w, r)
			return
		}

		http.NotFound(w, r)
	}
}
func (f *WebSite) handleSignupGet() http.HandlerFunc {

	indexTemplate := wd + "/frontend/private/signup.tmpl.html"
	tmpl, err := parseTemplates(indexTemplate)
	if err != nil {
		log.Println("error index tmpl: ", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		if r.URL.Path != "/signup" {
			http.NotFound(w, r)
			return
		}

		if err := tmpl.Execute(w, f.globalInfo); err != nil {
			log.Println("template err index: ", err)
		}
	}
}
func (f *WebSite) handleSignupPOST() http.HandlerFunc {
	signupTemplate := wd + "/frontend/private/signup.tmpl.html"
	tmpl, err := parseTemplates(signupTemplate)
	if err != nil {
		log.Println("error index tmpl: ", err)
	}

	type signupData struct {
		Email       string
		TeamName    string
		Password    string
		SignupError string
	}

	readParams := func(r *http.Request) (signupData, error) {
		data := signupData{
			Email:    r.PostFormValue("signupemail"), // removed due to GDPR
			TeamName: strings.TrimSpace(r.PostFormValue("signupusername")),
			Password: r.PostFormValue("signuppassword"),
		}

		if len(data.Password) < 6 {
			return data, fmt.Errorf("Password needs to be at least %d characters", 6)
		}

		if len(data.Password) > 20 {
			return data, fmt.Errorf("The maximum password length is %d characters", 20)
		}

		if data.Password != r.PostFormValue("signupcpassword") {
			return data, fmt.Errorf("Password needs to match")
		}

		return data, nil
	}

	displayErr := func(w http.ResponseWriter, params signupData, err error) {
		tmplData := f.globalInfo
		params.SignupError = err.Error()
		//tmplData.Content = params
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Println("template err signup: ", err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, f.maxReadBytes)
		params, err := readParams(r)
		if err != nil {
			displayErr(w, params, err)
			return
		}
		// max number for a team will be six so there is no capacity concept exactly.

		// todo: add recaptcha // this is no very important at the moment

		// there is no lab concept here !!!! so, all resources are shared !!

		//
		// todo: create VPN config per team
		// todo: the declaration of team who is red or blue will be gathered from signup page !!!
		// todo: this is crucial to create correct VPN configs on different interfaces
		// todo: when creating different interfaces per game, the port along IP should be randomized each time

		// todo: save team to cache

		// todo: save team to database

		// todo: set cookie time for team to be logged in

		// todo: forward to team to logged in interface of web

		// todo: when they signup they should be  automatically logged in as well.

		t := store.NewTeam("", strings.TrimSpace(params.TeamName), params.Password)

		ctx := context.TODO()
		serverPubKey, err := f.wgClient.GetPublicKey(ctx, &wg.PubKeyReq{PubKeyName: f.globalInfo.GameTag, PrivKeyName: f.globalInfo.GameTag})
		if err != nil {
			fmt.Printf("Err get public key wireguard  %v", err)
			return
		}
		_, err = f.wgClient.GenPrivateKey(ctx, &wg.PrivKeyReq{PrivateKeyName: f.globalInfo.GameTag + "_" + t.Id + "_"})
		if err != nil {
			fmt.Printf("Err gen private key wireguard  %v", err)
			return
		}

		//generate client public key
		//log.Info().Msgf("Generating public key for team %s", evTag+"_"+team+"_"+strconv.Itoa(ipAddr))
		_, err = f.wgClient.GenPublicKey(ctx, &wg.PubKeyReq{PubKeyName: f.globalInfo.GameTag + "_" + t.Id + "_", PrivKeyName: f.globalInfo.GameTag + "_" + t.Id + "_"})
		if err != nil {
			fmt.Printf("Err gen public key wireguard  %v", err)
			return
		}
		// get client public key
		//log.Info().Msgf("Retrieving public key for teaam %s", evTag+"_"+team+"_"+strconv.Itoa(ipAddr))
		resp, err := f.wgClient.GetPublicKey(ctx, &wg.PubKeyReq{PubKeyName: f.globalInfo.GameTag + "_" + t.Id + "_"})
		if err != nil {
			fmt.Printf("Error on GetPublicKey %v", err)
			return
		}

		//pIP := fmt.Sprintf("%d/32", len(ev.GetTeams())+2)
		peerIP := "45.11.23.4/32"
		////peerIP := strings.Replace(subnet, "1/24", fmt.Sprintf("%d/32", ipAddr), 1)
		//log.Info().Str("NIC", evTag).
		//	Str("AllowedIPs", peerIP).
		//	Str("PublicKey ", resp.Message).Msgf("Generating ip address for peer %s, ip address of peer is %s ", team, peerIP)
		addPeerResp, err := f.wgClient.AddPeer(ctx, &wg.AddPReq{
			Nic:        f.globalInfo.GameTag,
			AllowedIPs: peerIP,
			PublicKey:  resp.Message,
		})
		if err != nil {
			fmt.Sprintf("Error on adding peer to interface %v\n", err)
			return
		}
		fmt.Printf("AddPEER RESPONSE:  %s", addPeerResp.Message)
		//log.Info().Str("Event: ", evTag).
		//	Str("Peer: ", team).Msgf("Message : %s", addPeerResp.Message)
		////get client privatekey
		//log.Info().Msgf("Retrieving private key for team %s", team)
		teamPrivKey, err := f.wgClient.GetPrivateKey(ctx, &wg.PrivKeyReq{PrivateKeyName: f.globalInfo.GameTag + "_" + t.Id + "_"})
		if err != nil {
			fmt.Sprintf("Error on getting priv key for team  %v\n", err)
			return
		}
		//log.Info().Msgf("Privatee key for team %s is %s ", team, teamPrivKey.Message)
		//log.Info().Msgf("Client configuration is created for server %s", endpoint)
		// creating client configuration file
		// fmt.Sprintf("%s/24", "10.4.2.1") > this should be the lab subnet, necessry subnet which is assigned to team as a lab when they signed up...
		// 87878 > value should be changed with the randomized port where is it created before initializing the interface of wireguard...
		// fmt.Sprintf("%s.defatt.haaukins.com:%d", f.globalInfo.GameTag, 87878) > the dns address of host should be taken from configuration file of defat.

		clientConfig := fmt.Sprintf(
			`[Interface]
Address = %s
PrivateKey = %s
DNS = 1.1.1.1
MTU = 1500
[Peer]
PublicKey = %s
AllowedIps = %s
Endpoint =  %s
PersistentKeepalive = 25
`, peerIP, teamPrivKey.Message, serverPubKey.Message, fmt.Sprintf("%s/24", "10.4.2.1"), fmt.Sprintf("%s.defatt.haaukins.com:%d", f.globalInfo.GameTag, 87878))
		f.globalInfo.Team.VPNConfig = clientConfig
		t.VPNConfig = clientConfig

		//// todo: Save team
		//if err := f.TeamStore.SaveTeam(t); err != nil {
		//	displayErr(w, params, err)
		//	return
		//}

		if err := f.loginTeam(w, r, &t); err != nil {
			displayErr(w, params, err)
			return
		}

		if err := tmpl.Execute(w, f.globalInfo); err != nil {
			log.Println("template err index: ", err)
		}
		//token, err := store.GetTokenForTeam(f.signingKey, &t)
		//if err != nil {
		//	fmt.Printf("Error on getting token from amigo %v", token)
		//	return
		//}
		//
		//if err := f.TeamStore.SaveTokenForTeam(token, &t); err != nil {
		//	fmt.Printf("Create token for team error %v", err)
		//	return
		//}
	}
}

func (f *WebSite) loginTeam(w http.ResponseWriter, r *http.Request, t *store.Team) error {
	token, err := store.GetTokenForTeam(f.signingKey, t)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{Name: "session", Value: token, MaxAge: f.cookieTTL})
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func parseTemplates(givenTemplate string) (*template.Template, error) {
	var tmpl *template.Template
	var err error
	tmpl, err = template.ParseFiles(
		wd+"/frontend/private/main.tmpl.html",
		wd+"/frontend/private/navbar.tmpl.html",
		givenTemplate,
	)
	return tmpl, err
}

func (f *WebSite) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

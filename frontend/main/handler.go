package main

import (
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
	h.HandleFunc("/login", f.handleLogin())
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
		// todo: check capacity of game

		// todo: add recaptcha

		// todo:  create new team with challenges

		t := store.NewTeam("", strings.TrimSpace(params.TeamName), params.Password)

		// todo: Save team
		if err := f.TeamStore.SaveTeam(t); err != nil {
			displayErr(w, params, err)
			return
		}

		if err := f.loginTeam(w, r, &t); err != nil {
			displayErr(w, params, err)
			return
		}
		token, err := store.GetTokenForTeam(f.signingKey, &t)
		if err != nil {
			fmt.Printf("Error on getting token from amigo %v", token)
			return
		}

		if err := f.TeamStore.SaveTokenForTeam(token, &t); err != nil {
			fmt.Printf("Create token for team error %v", err)
			return
		}
		//serverPubKey, err := f.wgClient.GetPublicKey(context.TODO(), &wg.PubKeyReq{PubKeyName: f.globalInfo.GameTag, PrivKeyName: f.globalInfo.GameTag})
		//if err != nil {
		//	fmt.Printf("Err get public key wireguard  %v", err)
		//	return
		//}

		//if err := hook(&t); err != nil { // assigning lab
		//	fmt.Printf("Problem in creating configuration files for team %v ", err)
		//}
	}
}

func (f *WebSite) handleLogin() http.HandlerFunc {
	get := f.handleLoginGet()
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get(w, r)
			return

		case http.MethodPost:

			return
		}

		http.NotFound(w, r)
	}
}

func (f *WebSite) handleLoginGet() http.HandlerFunc {

	indexTemplate := wd + "/frontend/private/login.tmpl.html"
	tmpl, err := parseTemplates(indexTemplate)
	if err != nil {
		log.Println("error index tmpl: ", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		if r.URL.Path != "/login" {
			http.NotFound(w, r)
			return
		}

		if err := tmpl.Execute(w, f.globalInfo); err != nil {
			log.Println("template err index: ", err)
		}
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

func main() {

	g := store.GameConfig{
		Name:       "Testing frontend !!",
		Tag:        "test",
		ScenarioID: 2,
		StartedAt:  nil,
		FinishedAt: nil,
	}
	thor := NewFrontend(g, nil)
	http.ListenAndServe(":7676", thor.Handler())

}

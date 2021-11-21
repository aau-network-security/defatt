package models

type Network struct {
	Name string `yaml:"name"`
	Tag  int    `yaml:"tag"`
}
type Host struct {
	Name string `yaml:"name"`
	// Type is the type of virtulization used, should be set to either docker or vbox
	Type string `yaml:"type"`
	// Networks is should be an array of the networks the hosts should be in
	Networks []string `yaml:"networks"`
	// Image should be either the docker image e.g. "registry.gitlab.com/haaukins/forensics/hidden-server" or a vbox ova e.g. "winxp.ova"
	Image string `yaml:"image"`
	// Resources should only be included for vboxes which has specific requirements for the host
	Resources Recources `yaml:"resources"`
}

type Recources struct {
	CPU uint `yaml:"cpu"`
	RAM uint `yaml:"ram"`
}

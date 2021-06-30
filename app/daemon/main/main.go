package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aau-network-security/defatt/app/daemon"
	"github.com/aau-network-security/defatt/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultConfigFile = "config.yml"
)

func handleCancel(clean func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info().Msgf("Shutting down gracefully...")
		if err := clean(); err != nil {
			log.Error().Msgf("Error while shutting down: %s", err)
			os.Exit(1)
		}
		log.Info().Msgf("Closed daemon")
		os.Exit(0)
	}()
}

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	confFilePtr := flag.String("config", defaultConfigFile, "configuration file")
	flag.Parse()

	conf, err := config.NewConfig(*confFilePtr)
	if err != nil {
		fmt.Printf("unable to read configuration file \"%s\": %s\n", *confFilePtr, err)
		return
	}

	// ensure that gRPC port is free to allocate
	conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", conf.DefatConfig.Port), time.Second)
	if conn != nil {
		_ = conn.Close()
		fmt.Printf("Checking gRPC port %s report: %v\n", fmt.Sprintf(":%d", conf.DefatConfig.Port), daemon.PortIsAllocatedError)
		return
	}

	d, err := daemon.New(conf)
	if err != nil {
		fmt.Printf("unable to create daemon: %s\n", err)
		return
	}

	handleCancel(func() error {
		return d.Close()
	})
	log.Info().Msgf("Started daemon")

	if err := d.Run(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/harlow/go-micro-services/registry"
	"github.com/harlow/go-micro-services/services/frontend"
	"github.com/harlow/go-micro-services/tracing"
	"github.com/harlow/go-micro-services/tune"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pyroscope "github.com/grafana/pyroscope-go"
)

func main() {

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()

	serverAddress := os.Getenv("PYROSCOPE_SERVER_ADDRESS")
	log.Info().Msg(serverAddress)
	log.Info().Msg("THIS WORKS deadbeef")

	if serverAddress == "" {
		serverAddress = "http://pyroscope:4040"
	}
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "frontend.service",
		ServerAddress:   serverAddress,
		Logger:          pyroscope.StandardLogger,
	})

	tune.Init()

	log.Info().Msg("Reading config...")
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Error().Msgf("Got error while reading config: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	serv_port, _ := strconv.Atoi(result["FrontendPort"])
	serv_ip := result["FrontendIP"]
	knative_dns := result["KnativeDomainName"]

	log.Info().Msgf("Read target port: %v", serv_port)
	log.Info().Msgf("Read consul address: %v", result["consulAddress"])
	log.Info().Msgf("Read jaeger address: %v", result["jaegerAddress"])
	var (
		// port       = flag.Int("port", 5000, "The server port")
		jaegeraddr = flag.String("jaegeraddr", result["jaegerAddress"], "Jaeger address")
		consuladdr = flag.String("consuladdr", result["consulAddress"], "Consul address")
	)
	flag.Parse()
	log.Info().Msgf("Initializing jaeger agent [service name: %v | host: %v]...", "frontend", *jaegeraddr)
	tracer, err := tracing.Init("frontend", *jaegeraddr)
	if err != nil {
		log.Panic().Msgf("Got error while initializing jaeger agent: %v", err)
	}
	log.Info().Msg("Jaeger agent initialized")

	log.Info().Msgf("Initializing consul agent [host: %v]...", *consuladdr)
	registry, err := registry.NewClient(*consuladdr)
	if err != nil {
		log.Panic().Msgf("Got error while initializing consul agent: %v", err)
	}
	log.Info().Msg("Consul agent initialized")

	srv := &frontend.Server{
		KnativeDns: knative_dns,
		Registry:   registry,
		Tracer:     tracer,
		IpAddr:     serv_ip,
		Port:       serv_port,
	}

	log.Info().Msg("Starting server...")
	log.Fatal().Msg(srv.Run().Error())
}

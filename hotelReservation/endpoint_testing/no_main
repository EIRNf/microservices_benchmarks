package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"hotelReservation/registry"
	"io/ioutil"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	profile "github.com/harlow/go-micro-services/services/profile/proto"
	recommendation "github.com/harlow/go-micro-services/services/recommendation/proto"
	reservation "github.com/harlow/go-micro-services/services/reservation/proto"
	search "github.com/harlow/go-micro-services/services/search/proto"
	user "github.com/harlow/go-micro-services/services/user/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	consul "github.com/hashicorp/consul/api"
)

var knative_dns string
var registryClient *registry.Client

func main() {

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()

	log.Info().Msg("Starting tests...")

	log.Info().Msg("Reading config...")
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Error().Msgf("Got error while reading config: %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	log.Info().Msg("Successfull")

	var (
		// port       = flag.Int("port", 8083, "Server port")
		jaegeraddr = flag.String("jaegeraddr", result["jaegerAddress"], "Jaeger address")
		consuladdr = flag.String("consuladdr", result["consulAddress"], "Consul address")
		searchaddr = flag.String("SearchPort", result["SearchPort"], "Search port")
	)
	knative_dns = result["KnativeDomainName"]

	log.Info().Msgf("Read consul address: %v", jaegeraddr)
	log.Info().Msgf("Read jaeger address: %v", consuladdr)
	log.Info().Msgf("Read search address: %v", searchaddr)

	flag.Parse()

	log.Info().Msgf("Initializing consul agent [host: %v]...", *consuladdr)
	registryClient, err = registry.NewClient(*consuladdr)
	if err != nil {
		log.Panic().Msgf("Got error while initializing consul agent: %v", err)
	}

	log.Info().Msg("Initializing gRPC clients...")
	searchClient, err := initSearchClient("srv-search")
	if err != nil {
		log.Err(err)
	}

	profileClient, err := initProfileClient("srv-profile")
	if err != nil {
		log.Err(err)
	}

	recommendationClient, err := initRecommendationClient("srv-recommendation")
	if err != nil {
		log.Err(err)
	}

	userClient, err := initUserClient("srv-user")
	if err != nil {
		log.Err(err)
	}

	reservationClient, err := initReservation("srv-reservation")
	if err != nil {
		log.Err(err)
	}
	log.Info().Msg("Successfull")

	fmt.Printf("searchClient: %v\n", searchClient)
	fmt.Printf("profileClient: %v\n", profileClient)
	fmt.Printf("recommendationClient: %v\n", recommendationClient)
	fmt.Printf("userClient: %v\n", userClient)
	fmt.Printf("reservationClient: %v\n", reservationClient)

}

func getGprcConn(name string) (*grpc.ClientConn, error) {
	log.Info().Msg("get Grpc conn is :")
	log.Info().Msg(knative_dns)
	log.Info().Msg(fmt.Sprintf("%s.%s", name, knative_dns))

	cfg := consul.DefaultConfig()
	target := fmt.Sprintf("consul://%s/%s", cfg.Address, name)

	return grpc.Dial(
		target,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// return grpc.Dial(
	// 	fmt.Sprintf("%s:///%s", knative_dns, name),
	// 	grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if s.KnativeDns != "" {
	// 	return dialer.Dial(
	// 		fmt.Sprintf("%s.%s", name, s.KnativeDns),
	// 		dialer.WithTracer(s.Tracer))
	// } else {
	// 	return dialer.Dial(
	// 		name,
	// 		dialer.WithTracer(s.Tracer),
	// 		dialer.WithBalancer(s.Registry.Client),
	// 	)
	// }
}

func initSearchClient(name string) (search.SearchClient, error) {
	conn, err := getGprcConn(name)
	if err != nil {
		log.Err(err)
	}
	return search.NewSearchClient(conn), err
}

func initProfileClient(name string) (profile.ProfileClient, error) {
	conn, err := getGprcConn(name)
	if err != nil {
		log.Err(err)
	}
	return profile.NewProfileClient(conn), err
}

func initRecommendationClient(name string) (recommendation.RecommendationClient, error) {
	conn, err := getGprcConn(name)
	if err != nil {
		log.Err(err)
	}
	return recommendation.NewRecommendationClient(conn), err
}

func initUserClient(name string) (user.UserClient, error) {
	conn, err := getGprcConn(name)
	if err != nil {
		log.Err(err)
	}
	return user.NewUserClient(conn), err
}

func initReservation(name string) (reservation.ReservationClient, error) {
	conn, err := getGprcConn(name)
	if err != nil {
		log.Err(err)
	}
	return reservation.NewReservationClient(conn), err

}

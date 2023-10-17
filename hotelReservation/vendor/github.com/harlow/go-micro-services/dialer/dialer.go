package dialer

import (
	"fmt"

	"github.com/harlow/go-micro-services/registry"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"

	// lb "github.com/olivere/grpc/lb/consul"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

// DialOption allows optional config for dialer
type DialOption func(name string) (grpc.DialOption, error)

// Dial returns a load balanced grpc client conn with tracing interceptor
func Dial(name string, registry *registry.Client, tracer opentracing.Tracer) (*grpc.ClientConn, error) {

	client := registry
	address := client.Config.Address
	target := fmt.Sprintf("consul://%s/%s", address, name)

	log.Info().Msg("get Grpc conn is :")
	// log.Info().Msgf("Consul Client [host: %v]...", *client)
	// log.Info().Msgf("Consul Config [host: %v]...", client.Config)
	// log.Info().Msgf("Address String [host: %sv]...", address)
	log.Info().Msg(fmt.Sprintf("consul://%s/%s", address, name))

	conn, err := grpc.Dial(
		target,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
	)

	// log.Info().Msgf("Conn: [%v]...", conn)
	// log.Info().Msgf("err: [%v]...", err)

	return conn, err
}

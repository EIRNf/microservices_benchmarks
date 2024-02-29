package endpoint_testing

import (
	"context"
	"testing"

	"github.com/EIRNf/notnets_grpc"
	profile_service "github.com/harlow/go-micro-services/services/profile"
	profile "github.com/harlow/go-micro-services/services/profile/proto"
)

func Test_ser_deser(t *testing.T) {

	// w := http.Wrote
	// var r *http.Request

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := context.Background()

	// sLat, sLon := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	// if sLat == "" || sLon == "" {
	// 	http.Error(w, "Please specify location params", http.StatusBadRequest)
	// 	return
	// }
	// // Lat, _ := strconv.ParseFloat(sLat, 64)
	// // // lat := float64(Lat)
	// // Lon, _ := strconv.ParseFloat(sLon, 64)
	// // // lon := float64(Lon)

	// require := r.URL.Query().Get("require")
	// if require != "dis" && require != "rate" && require != "price" {
	// 	http.Error(w, "Please specify require params", http.StatusBadRequest)
	// 	return
	// }

	// recommend hotels
	// recResp, err := s.recommendationClient.GetRecommendations(ctx, &recommendation.Request{
	// 	Require: require,
	// 	Lat:     float64(lat),
	// 	Lon:     float64(lon),
	// })
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	svc := &profile_service.Server{}
	// svc := &channel_tests_service.TestServer{}
	svr := notnets_grpc.NewNotnetsServer()

	//Register Server and instantiate with necessary information
	// test_hello_service.RegisterTestServiceServer(svr, svc)
	profile.RegisterProfileServer(svr, svc)

	//Create Listener
	lis := notnets_grpc.Listen("http://127.0.0.1:8080/hello")

	go svr.Serve(lis)

	// grab locale from query params or default to en
	// locale := r.URL.Query().Get("locale")
	// if locale == "" {
	locale := "en"
	// }

	cc, _ := notnets_grpc.Dial("testtest", "http://127.0.0.1:8080/hello", notnets_grpc.MESSAGE_SIZE)
	profileClient := profile.NewProfileClient(cc)

	hotelIds := []string{"39", "42"}

	// hotel profiles
	profileResp, _ := profileClient.GetProfiles(ctx, &profile.Request{
		HotelIds: hotelIds,
		Locale:   locale,
	})

	_ = profileResp.Hotels
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// json.NewEncoder(w).Encode(geoJSONResponse(profileResp.Hotels))

	// profile.Request()

}

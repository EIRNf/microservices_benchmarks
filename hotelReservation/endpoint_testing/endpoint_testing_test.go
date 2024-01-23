package endpoint_testing

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"time"

	// "hotelReservation/registry"
	"math/big"
	"net/http"
	"testing"

	"github.com/loov/hrtime"
)

// const localURL = "http://localhost:5000"
// const apiURL = "http://192.168.49.2:30052"

var apiURL = flag.String("url", "http://192.168.49.2:32567", "IP/Port of HotelReservation Frontend")

var endpoint_bench = flag.String("endpoint", "all", "Test to run:hotels,recommendations,user,reservation")

var num_requests = flag.Int("reqs", 1000, "Number of requests per instance")

func Benchmark_All(b *testing.B) {

	//Stochastic Bench
	// Hotels: 0.6
	// recommend: 0.39
	// user: 0.005
	// reserve: 0.005

	switch *endpoint_bench {
	case "hotels":
		// b.Run("hotels_endpoint", Bench_Hotels_Endpoint)
		b.Run("hotels_endpoint", Bench_Hotels_Endpoint_Histogram)
	case "recommendations":
	case "user":
	case "reservation":
	case "all":
	default:
		fmt.Println("Unknown bench")
	}

	// b.Run("recommendations_endpoint", Test_Recommendation_Endpoint)
	// b.Run("user_endpoint", Test_User_Endpoint)
	// b.Run("reservation_endpoint", Test_Reservation_Endpoint)

}

func Bench_Hotels_Endpoint(b *testing.B) {

	in_date_str := "2015-04-09"
	out_date_str := "2015-04-24"
	lat_val := "38.0235"
	lon_val := "-122.095"

	//Construct request
	inDate := "?inDate=" + in_date_str
	outDate := "&outDate=" + out_date_str
	lat := "&lat=" + lat_val
	lon := "&lon=" + lon_val

	request_string := *apiURL + "/hotels" + inDate + outDate + lat + lon

	for i := 0; i < b.N; i++ {
		resp, err := http.Get(request_string)

		//Failed Request
		if err != nil {
			b.Fatalf("RPC failed: %v", err)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			b.Fatalf("Error reading body: %v", err)
		}

		//Check StatusCode
		if resp.Status != status_code_ok {
			b.Fatalf("Status code not ok, instead: %v", resp.Status)
		}

		//Unexpected Payload
		if !bytes.Equal(data, hotels_expected) {
			b.Fatalf("wrong payload returns: expecting %v; got %v", hotels_expected, data)
		}
		resp.Body.Close()
	}

	fmt.Printf("NUM_REQUESTS: %d \n", b.N)
	fmt.Printf("EXECUTION_LENGTH: %f \n", b.Elapsed().Seconds())
	fmt.Printf("THROUGHPUT: %f reqs/sec\n", float64(b.N)/b.Elapsed().Seconds())
	fmt.Printf("AVERAGE_LATENCY:%f\n", float64(b.Elapsed().Milliseconds())/float64(b.N))
	// fmt.Printf("LATENCY_STATS:\n %s", bench.Histogram(10).StringStats())

}

func Bench_Hotels_Endpoint_Histogram(b *testing.B) {

	in_date_str := "2015-04-09"
	out_date_str := "2015-04-24"
	lat_val := "38.0235"
	lon_val := "-122.095"

	//Construct request
	inDate := "?inDate=" + in_date_str
	outDate := "&outDate=" + out_date_str
	lat := "&lat=" + lat_val
	lon := "&lon=" + lon_val

	request_string := *apiURL + "/hotels" + inDate + outDate + lat + lon

	bench := hrtime.NewBenchmark(*num_requests)
	// defer bench.HistogramClamp()

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	for bench.Next() {
		resp, err := http.Get(request_string)

		//Failed Request
		if err != nil {
			b.Fatalf("RPC failed: %v", err)
		}

		// // defer resp.Body.Close()
		// data, err := io.ReadAll(resp.Body)
		resp.Body.Close() // close the response body
		// if err != nil {
		// 	b.Fatalf("Error reading body: %v", err)
		// }

		// //Check StatusCode
		// if resp.Status != status_code_ok {
		// 	b.Fatalf("Status code not ok, instead: %v", resp.Status)
		// }

		// // //Unexpected Payload
		// if !bytes.Equal(data, hotels_expected) {
		// 	b.Fatalf("wrong payload returns: expecting %v; got %v", hotels_expected, data)
		// }

	}

	// fmt.Println(bench.Histogram(10))

	//Get Average Latency
	//Get Throughput
	fmt.Printf("NUM_REQUESTS: %d \n", *num_requests)
	fmt.Printf("EXECUTION_LENGTH: %f \n", b.Elapsed().Seconds())
	fmt.Printf("THROUGHPUT: %f reqs/sec\n", float64(*num_requests)/b.Elapsed().Seconds())
	fmt.Printf("AVERAGE_LATENCY:%s\n", time.Duration(bench.Histogram(10).Average).String())
	fmt.Printf("LATENCY_STATS:\n %s", bench.Histogram(10).StringStats())

	// runs := bench.Float64s()
	// fmt.Printf("Mean: %f\n", stats.Mean(runs)*0.001)
	// fmt.Printf("StdDev: %f\n", stats.StdDev(runs)*0.001)
	// fmt.Printf("NumElements: %d\n", len(runs))
	// fmt.Printf("Time in Microsecondes: %d \n", b.Elapsed().Microseconds())
	// fmt.Printf("Time in Seconds: %f \n", b.Elapsed().Seconds())
	// fmt.Printf("Sanity Check Time Microseconds: %f \n", vec.Sum(runs)*0.001)
	// fmt.Printf("Throughput: %f \n", float64(len(runs))/b.Elapsed().Seconds())

	//PRINT OUT ALL RUNS
	// runs := bench.Float64s()
	// fmt.Printf("RUNS:%f \n", runs)
}

func Test_All(t *testing.T) {

	//Check for argument, exit if notx
	// url := flag.Arg(0)
	// if url == nil {
	// 	//No IP declared attempt
	// 	t.Fatalf("No target frontend declared")
	// }

	//Test Endpoints
	t.Run("hotels_endpoint", Test_Hotels_Endpoint)
	t.Run("recommendations_endpoint", Test_Recommendation_Endpoint)
	t.Run("user_endpoint", Test_User_Endpoint)
	t.Run("reservation_endpoint", Test_Reservation_Endpoint)

}

var (
	// TestServer server

	status_code_ok = "200 OK"

	//Expected payload for request:
	//"http://192.168.49.2:30052/hotels?inDate=2015-04-09&outDate=2015-04-24&lat=38.0235&lon=-122.095"
	hotels_expected = []byte{123, 34, 102, 101, 97, 116, 117, 114, 101, 115, 34, 58, 91, 123, 34, 103, 101, 111, 109, 101, 116, 114, 121, 34, 58, 123, 34, 99, 111, 111, 114, 100, 105, 110, 97, 116, 101, 115, 34, 58, 91, 45, 49, 50, 50, 46, 48, 57, 56, 48, 49, 44, 51, 56, 46, 48, 49, 55, 53, 93, 44, 34, 116, 121, 112, 101, 34, 58, 34, 80, 111, 105, 110, 116, 34, 125, 44, 34, 105, 100, 34, 58, 34, 51, 57, 34, 44, 34, 112, 114, 111, 112, 101, 114, 116, 105, 101, 115, 34, 58, 123, 34, 110, 97, 109, 101, 34, 58, 34, 83, 116, 46, 32, 82, 101, 103, 105, 115, 32, 83, 97, 110, 32, 70, 114, 97, 110, 99, 105, 115, 99, 111, 34, 44, 34, 112, 104, 111, 110, 101, 95, 110, 117, 109, 98, 101, 114, 34, 58, 34, 40, 52, 49, 53, 41, 32, 50, 56, 52, 45, 52, 48, 51, 57, 34, 125, 44, 34, 116, 121, 112, 101, 34, 58, 34, 70, 101, 97, 116, 117, 114, 101, 34, 125, 44, 123, 34, 103, 101, 111, 109, 101, 116, 114, 121, 34, 58, 123, 34, 99, 111, 111, 114, 100, 105, 110, 97, 116, 101, 115, 34, 58, 91, 45, 49, 50, 50, 46, 48, 55, 52, 48, 48, 53, 44, 51, 56, 46, 48, 51, 53, 53, 93, 44, 34, 116, 121, 112, 101, 34, 58, 34, 80, 111, 105, 110, 116, 34, 125, 44, 34, 105, 100, 34, 58, 34, 52, 50, 34, 44, 34, 112, 114, 111, 112, 101, 114, 116, 105, 101, 115, 34, 58, 123, 34, 110, 97, 109, 101, 34, 58, 34, 83, 116, 46, 32, 82, 101, 103, 105, 115, 32, 83, 97, 110, 32, 70, 114, 97, 110, 99, 105, 115, 99, 111, 34, 44, 34, 112, 104, 111, 110, 101, 95, 110, 117, 109, 98, 101, 114, 34, 58, 34, 40, 52, 49, 53, 41, 32, 50, 56, 52, 45, 52, 48, 52, 50, 34, 125, 44, 34, 116, 121, 112, 101, 34, 58, 34, 70, 101, 97, 116, 117, 114, 101, 34, 125, 93, 44, 34, 116, 121, 112, 101, 34, 58, 34, 70, 101, 97, 116, 117, 114, 101, 67, 111, 108, 108, 101, 99, 116, 105, 111, 110, 34, 125, 10}

	//Expected payload for request:
	//"http://192.168.49.2:30052/recommendations?require=dis&lat=38.0235&lon=-122.095"
	recommendation_expected = []byte{123, 34, 102, 101, 97, 116, 117, 114, 101, 115, 34, 58, 91, 123, 34, 103, 101, 111, 109, 101, 116, 114, 121, 34, 58, 123, 34, 99, 111, 111, 114, 100, 105, 110, 97, 116, 101, 115, 34, 58, 91, 45, 49, 50, 50, 46, 48, 57, 48, 48, 48, 52, 44, 51, 56, 46, 48, 50, 51, 53, 48, 50, 93, 44, 34, 116, 121, 112, 101, 34, 58, 34, 80, 111, 105, 110, 116, 34, 125, 44, 34, 105, 100, 34, 58, 34, 52, 48, 34, 44, 34, 112, 114, 111, 112, 101, 114, 116, 105, 101, 115, 34, 58, 123, 34, 110, 97, 109, 101, 34, 58, 34, 83, 116, 46, 32, 82, 101, 103, 105, 115, 32, 83, 97, 110, 32, 70, 114, 97, 110, 99, 105, 115, 99, 111, 34, 44, 34, 112, 104, 111, 110, 101, 95, 110, 117, 109, 98, 101, 114, 34, 58, 34, 40, 52, 49, 53, 41, 32, 50, 56, 52, 45, 52, 48, 52, 48, 34, 125, 44, 34, 116, 121, 112, 101, 34, 58, 34, 70, 101, 97, 116, 117, 114, 101, 34, 125, 93, 44, 34, 116, 121, 112, 101, 34, 58, 34, 70, 101, 97, 116, 117, 114, 101, 67, 111, 108, 108, 101, 99, 116, 105, 111, 110, 34, 125, 10}

	user_expected = []byte{123, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 34, 76, 111, 103, 105, 110, 32, 115, 117, 99, 99, 101, 115, 115, 102, 117, 108, 108, 121, 33, 34, 125, 10}

	//Expected payload for request:
	//http://192.168.49.2:30052/reservation?inDate=2015-04-09&outDate=2015-04-24&lat=38.0235&lon=-122.095&hotelId=5&customerName=Cornell_311&username=Cornell_311&password=311311311311311311311311311311&number=1
	reservation_expected = []byte{123, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 34, 82, 101, 115, 101, 114, 118, 101, 32, 115, 117, 99, 99, 101, 115, 115, 102, 117, 108, 108, 121, 33, 34, 125, 10}
)

func Test_Hotels_Endpoint(t *testing.T) {

	in_date_str := "2015-04-09"
	out_date_str := "2015-04-24"
	lat_val := "38.0235"
	lon_val := "-122.095"

	//Construct request
	inDate := "?inDate=" + in_date_str
	outDate := "&outDate=" + out_date_str
	lat := "&lat=" + lat_val
	lon := "&lon=" + lon_val

	request_string := *apiURL + "/hotels" + inDate + outDate + lat + lon

	resp, err := http.Get(request_string)

	//Failed Request
	if err != nil {
		t.Fatalf("RPC failed: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}

	//Check StatusCode
	if resp.Status != status_code_ok {
		t.Fatalf("Status code not ok, instead: %v", resp.Status)
	}

	//Unexpected Payload
	if !bytes.Equal(data, hotels_expected) {
		t.Fatalf("wrong payload returns: expecting %v; got %v", hotels_expected, data)
	}

	//Verify Headers
	//Verify Trailers

}

func Test_Recommendation_Endpoint(t *testing.T) {

	param := "dis"
	lat_val := "38.0235"
	lon_val := "-122.095"

	//Construct request
	req_param := "?require=" + param
	lat := "&lat=" + lat_val
	lon := "&lon=" + lon_val

	// 	local method = "GET"
	//   local path = url .. "/recommendations?require=" .. req_param ..
	//     "&lat=" .. tostring(lat) .. "&lon=" .. tostring(lon)

	request_string := *apiURL + "/recommendations" + req_param + lat + lon

	resp, err := http.Get(request_string)

	//Failed Request
	if err != nil {
		t.Fatalf("RPC failed: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}

	//Check StatusCode
	if resp.Status != status_code_ok {
		t.Fatalf("Status code not ok, instead: %v", resp.Status)
	}

	//Unexpected Payload
	if !bytes.Equal(data, recommendation_expected) {
		t.Fatalf("wrong payload returns: expecting %v; got %v", recommendation_expected, data)
	}
}

func Test_User_Endpoint(t *testing.T) {

	//Get username and password
	suffix, _ := rand.Int(rand.Reader, big.NewInt(500))
	user := "Cornell_" + suffix.String()
	pass := ""
	for j := 0; j < 10; j++ {
		pass += suffix.String()
	}

	user_name := "?username=" + user
	password := "&password=" + pass

	request_string := *apiURL + "/user" + user_name + password

	resp, err := http.Get(request_string)

	//Failed Request
	if err != nil {
		t.Fatalf("RPC failed: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}

	//Check StatusCode
	if resp.Status != status_code_ok {
		t.Fatalf("Status code not ok, instead: %v", resp.Status)
	}

	//Unexpected Payload
	if !bytes.Equal(data, user_expected) {
		t.Fatalf("wrong payload returns: expecting %v; got %v", user_expected, data)
	}

	//Verify Headers
	//Verify Trailers
}

func Test_Reservation_Endpoint(t *testing.T) {

	//Get username and password
	suffix, _ := rand.Int(rand.Reader, big.NewInt(500))
	user := "Cornell_" + suffix.String()
	pass := ""
	for j := 0; j < 10; j++ {
		pass += suffix.String()
	}

	in_date_str := "2015-04-09"
	out_date_str := "2015-04-24"
	lat_val := "38.0235"
	lon_val := "-122.095"
	id, _ := rand.Int(rand.Reader, big.NewInt(80))
	hotel_id_val := id.String()
	customer := user
	num := "1"

	//Construct request
	inDate := "?inDate=" + in_date_str
	outDate := "&outDate=" + out_date_str
	lat := "&lat=" + lat_val
	lon := "&lon=" + lon_val
	hotelId := "&hotelId=" + hotel_id_val
	customerName := "&customerName=" + customer
	user_name := "&username=" + user
	password := "&password=" + pass
	number := "&number=" + num

	request_string := *apiURL + "/reservation" + inDate + outDate + lat + lon + hotelId + customerName + user_name + password + number

	resp, err := http.Get(request_string)

	//Failed Request
	if err != nil {
		t.Fatalf("RPC failed: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}

	//Check StatusCode
	if resp.Status != status_code_ok {
		t.Fatalf("Status code not ok, instead: %v", resp.Status)
	}

	//Unexpected Payload
	if !bytes.Equal(data, reservation_expected) {
		t.Fatalf("wrong payload returns: expecting %v; got %v", reservation_expected, data)
	}

	//Verify Headers
	//Verify Trailers
}

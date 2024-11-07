package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "example.com/m/proto"
)

const address = "0.0.0.0:50050"

// Server struct to implement WeatherService
type server struct {
	pb.UnimplementedWeatherServiceServer
}

// Function to fetch weather data (simulated)
func fetchWeatherSim(city string, scale pb.TemperatureScale) float64 {
	// Simulated temperature data (Celsius)
	temperatureCelsius := 20 + float64(len(city)%10)

	// Convert temperature to Fahrenheit if required
	if scale == pb.TemperatureScale_FAHRENHEIT {
		return (temperatureCelsius * 9 / 5) + 32
	}

	return temperatureCelsius
}

// Bi-directional streaming method to get weather data
func (s *server) GetWeatherStream(stream pb.WeatherService_GetWeatherStreamServer) error {
	log.Println("GetWeatherStream was invoked")

	// Loop to receive requests and send responses
	for {
		// Receive request from client
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("error receiving request: %v\n", err)
			return err
		}

		log.Printf("received request for city: %s, scale: %s\n", req.City, req.Scale)

		// Fetch weather data (simulated)
		temperature := fetchWeatherSim(req.City, req.Scale)

		// Send response to client
		err = stream.Send(&pb.WeatherResponse{
			City:        req.City,
			Temperature: temperature,
			Scale:       req.Scale,
		})
		if err != nil {
			log.Printf("error sending response: %v\n", err)
			return err
		}

		log.Printf("sent response for city: %s, temperature: %f, scale: %s\n", req.City, temperature, req.Scale)
	}
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on %v: %v", address, err)
	}

	log.Printf("listening on %s\n", address)

	s := grpc.NewServer()
	pb.RegisterWeatherServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

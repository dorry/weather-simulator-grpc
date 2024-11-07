package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "example.com/m/proto"
)

const address = "localhost:50050"

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewWeatherServiceClient(conn)

	// Create a context for the stream
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a stream to send requests and receive responses
	stream, err := client.GetWeatherStream(ctx)
	if err != nil {
		log.Fatalf("failed to create stream: %v", err)
	}

	// Cities to request weather data for
	cities := []string{"Cairo", "New York", "London", "Tokyo"}

	// Temperature scales to request
	scales := []pb.TemperatureScale{pb.TemperatureScale_CELSIUS, pb.TemperatureScale_FAHRENHEIT}

	// Loop to send requests and receive responses
	for _, city := range cities {
		for _, scale := range scales {
			// Send request to server
			err := stream.Send(&pb.WeatherRequest{
				City:  city,
				Scale: scale,
			})
			if err != nil {
				log.Printf("error sending request: %v\n", err)
				return
			}

			log.Printf("sent request for city: %s, scale: %s\n", city, scale)

			// Receive response from server
			resp, err := stream.Recv()
			if err != nil {
				log.Printf("error receiving response: %v\n", err)
				return
			}

			log.Printf("received response for city: %s, temperature: %f, scale: %s\n", resp.City, resp.Temperature, resp.Scale)
		}
	}

	stream.CloseSend()

}

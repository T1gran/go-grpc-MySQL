package main

import (
	"context"
	"google.golang.org/grpc"
	pb "intern/protos"
	"io"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMovieClient(conn)
	runGetMovies(client)
	runGetMovie(client, "1")
}

func runGetMovies(client pb.MovieClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetMovies(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovies(_) = _, %v", client, err)
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetMovies(_) = _, %v", client, err)
		}
		log.Printf("MovieInfo: %v", row)
	}
}

func runGetMovie(client pb.MovieClient, movieid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: movieid}
	res, err := client.GetMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovie(_) = _, %v", client, err)
	}
	log.Printf("MovieInfo: %v", res)
}

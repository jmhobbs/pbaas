package main

import (
	"context"
	"log"

	"github.com/jmhobbs/pbaas/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProgressBarServiceClient(conn)

	r, err := c.NewProgressBar(context.Background(), &pb.NewProgressBarRequest{5})
	if err != nil {
		log.Fatalf("could not create progress bar: %v", err)
	}

	log.Println("New Progress Bar:", r)

	// Update progress
	res, err := c.UpdateProgressBar(context.Background(), &pb.UpdateProgressBarRequest{
		Token: r.Token,
		Id: r.Id,
		NewProgressValue: 20,
	})
	if err != nil {
		log.Fatalf("could not update progress bar: %v", err)
	}

	log.Println("New Progress Bar Value:", res)

}
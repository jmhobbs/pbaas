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

	r1, err := c.NewProgressBar(context.Background(), &pb.NewProgressBarRequest{5})
	if err != nil {
		log.Fatalf("could not create progress bar: %v", err)
	}

	log.Println("New Progress Bar:", r1)

	// Get progress
	r2, err := c.GetProgressBarStatus(context.Background(), &pb.ProgressBarStatusRequest{
		Ids: []string{r1.GetId()},
	})
	if err != nil {
		log.Fatalf("could not get progress bar: %v", err)
	}

	log.Println("Got Progress Bar:", r2)


	// Update progress
	r3, err := c.UpdateProgressBar(context.Background(), &pb.UpdateProgressBarRequest{
		Token: r1.Token,
		Id: r1.Id,
		NewProgressValue: 20,
	})
	if err != nil {
		log.Fatalf("could not update progress bar: %v", err)
	}

	log.Println("New Progress Bar Value:", r3)

	// Get progress
	r4, _ := c.GetProgressBarStatus(context.Background(), &pb.ProgressBarStatusRequest{
		Ids: []string{r1.GetId()},
	})

	log.Println("Got Progress Bar:", r4)

	// Delete progress
	r5, _ := c.DeleteProgressBar(context.Background(), &pb.ProgressBarStatusRequest{
		Ids: []string{r1.GetId()},
	})

	log.Println("Deleted Progress Bar:", r5)

}
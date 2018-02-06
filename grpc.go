package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/jmhobbs/pbaas/pb"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type progressBarServiceServer struct {
	db ProgressDB
}

func (ps progressBarServiceServer) NewProgressBar(ctx context.Context, pr *pb.NewProgressBarRequest) (*pb.NewProgressBarResponse, error) {
	id := uuid.Must(uuid.NewV4()).String()
	token := newToken()
	ps.db.Create(id, token, pr.GetStartingProgress())
	return &pb.NewProgressBarResponse{Id: id, Token: token}, nil
}

func (ps progressBarServiceServer) GetProgressBarStatus(ctx context.Context, pr *pb.ProgressBarStatusRequest) (*pb.ProgressBarStatusResponse, error) {
	pbs := []*pb.ProgressBar{}
	for _, id := range pr.GetIds() {
		pbs = append(pbs, &pb.ProgressBar{id, ps.db.Get(id)})
	}
	return &pb.ProgressBarStatusResponse{ProgressBars: pbs}, nil
}

func (ps progressBarServiceServer) Serve(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProgressBarServiceServer(s, ps)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

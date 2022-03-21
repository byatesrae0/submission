package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"google.golang.org/grpc"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// racingCli is shared among tests & used to talk to the racing grpc server started in testMain().
var racingCli racing.RacingClient

func TestMain(m *testing.M) {
	code, err := testMain(m)
	if err != nil {
		log.Printf("[ERR] Unexpected error encounted during test execution: %v", err)
	}
	os.Exit(code)
}

func testMain(m *testing.M) (int, error) {
	const errorExitCode = 1

	// Start server
	go func() {
		if err := run(); err != nil {
			log.Printf("[ERR] Run: %v", err)
		}
	}()

	// Connect to server
	gRPCOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(*grpcEndpoint, gRPCOpts...)
	if err != nil {
		return errorExitCode, fmt.Errorf("failed to dial grpc : %s %s", *grpcEndpoint, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("[ERR] Failed to close connection: %v", err)
		}
	}()

	racingCli = racing.NewRacingClient(conn)

	return m.Run(), nil
}

package tests

import (
	"context"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
)

var mappedPort nat.Port

func TestMain(m *testing.M) {
	var err error
	rabbitmqContainer := createRabbitmqContainer()

	mappedPort, err = rabbitmqContainer.MappedPort(context.Background(), "5672/tcp")
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()

	ctx := context.Background()

	err = rabbitmqContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

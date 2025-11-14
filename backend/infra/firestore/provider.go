package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	ProjectID    string
	EmulatorAddr string
)

func ProvideEmulatorClient(ctx context.Context, projectID ProjectID, emulatorAddr EmulatorAddr) (*firestore.Client, error) {
	conn, err := grpc.NewClient(string(emulatorAddr),
		grpc.WithPerRPCCredentials(emulatorCreds{}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return firestore.NewClient(ctx, string(projectID), option.WithGRPCConn(conn))
}

// from cloud.google.com/go/firestore.emulatorCreds
type emulatorCreds struct{}

var _ credentials.PerRPCCredentials = emulatorCreds{}

func (emulatorCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer owner"}, nil
}

func (emulatorCreds) RequireTransportSecurity() bool { return false }

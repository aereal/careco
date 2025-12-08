package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/google/wire"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	ProjectID    string
	EmulatorAddr string
	DatabaseID   string
)

var (
	ClientProvider = wire.NewSet(
		ProvideClient,
	)
)

func ProvideClient(ctx context.Context, dbID DatabaseID, projectID ProjectID, emulatorAddr EmulatorAddr, tp trace.TracerProvider) (*firestore.Client, error) {
	conn, err := grpc.NewClient(string(emulatorAddr),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(otelgrpc.WithPropagators(propagation.TraceContext{}), otelgrpc.WithTracerProvider(tp))),
		grpc.WithStatsHandler(firestoreAttrGetter{}),
		grpc.WithPerRPCCredentials(emulatorCreds{}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return firestore.NewClientWithDatabase(ctx, string(projectID), string(dbID), option.WithGRPCConn(conn))
}

// from cloud.google.com/go/firestore.emulatorCreds
type emulatorCreds struct{}

var _ credentials.PerRPCCredentials = emulatorCreds{}

func (emulatorCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer owner"}, nil
}

func (emulatorCreds) RequireTransportSecurity() bool { return false }

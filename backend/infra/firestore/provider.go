package firestore

import (
	"context"
	"fmt"
	"iter"
	"slices"

	"careco/backend/infra/gcp"

	"cloud.google.com/go/firestore"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	EmulatorAddr string
	DatabaseID   string
)

func ProvideClient(ctx context.Context, dbID DatabaseID, projectID gcp.ProjectID, tp trace.TracerProvider) (*firestore.Client, error) {
	opts := slices.Collect(asClientOptions(commonDialOptions(tp)))
	return firestore.NewClientWithDatabase(ctx, string(projectID), string(dbID), opts...)
}

func ProvideEmulatorClient(ctx context.Context, dbID DatabaseID, projectID gcp.ProjectID, emulatorAddr EmulatorAddr, tp trace.TracerProvider) (*firestore.Client, error) {
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(emulatorCreds{}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(string(emulatorAddr), slices.AppendSeq(opts, commonDialOptions(tp))...)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return firestore.NewClientWithDatabase(ctx, string(projectID), string(dbID), option.WithGRPCConn(conn))
}

func commonDialOptions(tp trace.TracerProvider) iter.Seq[grpc.DialOption] {
	return func(yield func(grpc.DialOption) bool) {
		if !yield(grpc.WithStatsHandler(otelgrpc.NewClientHandler(otelgrpc.WithPropagators(propagation.TraceContext{}), otelgrpc.WithTracerProvider(tp)))) {
			return
		}
		if !yield(grpc.WithStatsHandler(firestoreAttrGetter{})) {
			return
		}
	}
}

func asClientOptions(s iter.Seq[grpc.DialOption]) iter.Seq[option.ClientOption] {
	return func(yield func(option.ClientOption) bool) {
		for o := range s {
			if !yield(option.WithGRPCDialOption(o)) {
				return
			}
		}
	}
}

// from cloud.google.com/go/firestore.emulatorCreds
type emulatorCreds struct{}

var _ credentials.PerRPCCredentials = emulatorCreds{}

func (emulatorCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer owner"}, nil
}

func (emulatorCreds) RequireTransportSecurity() bool { return false }

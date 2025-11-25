package test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"careco/backend/infra/firestore"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func provideTracerProvider() trace.TracerProvider {
	return sdktrace.NewTracerProvider(sdktrace.WithBatcher(tracetest.NewInMemoryExporter()))
}

func provideContext(t *testing.T) context.Context {
	t.Helper()
	return t.Context()
}

func provideDatabaseID(t *testing.T) firestore.DatabaseID {
	t.Helper()
	h := sha256.New()
	h.Write([]byte(t.Name()))
	return firestore.DatabaseID("test_" + hex.EncodeToString(h.Sum(nil)))
}

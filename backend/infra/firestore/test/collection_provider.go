package test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"careco/backend/infra/firestore"

	sdk "cloud.google.com/go/firestore"
)

type tLike interface {
	Context() context.Context
	Name() string
}

type ctxKeySuffix struct{}

func Context(t tLike) context.Context {
	ctx := t.Context()
	sum := sha256.New().Sum([]byte(t.Name()))
	suffix := "_" + hex.EncodeToString(sum)
	return context.WithValue(ctx, ctxKeySuffix{}, suffix)
}

func suffixFromContext(ctx context.Context) string {
	if s, ok := ctx.Value(ctxKeySuffix{}).(string); ok {
		return s
	}
	return ""
}

func provideTestCollectionProvider(client *sdk.Client) firestore.CollectionProvider {
	return &testCollectionProvider{client: client}
}

type testCollectionProvider struct {
	client *sdk.Client
}

var _ firestore.CollectionProvider = (*testCollectionProvider)(nil)

func (p *testCollectionProvider) DrivingRecords(ctx context.Context) *sdk.CollectionRef {
	suffix := suffixFromContext(ctx)
	return p.client.Collection("driving_records" + suffix)
}

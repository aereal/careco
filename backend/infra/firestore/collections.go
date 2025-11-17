package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type CollectionProvider interface {
	DrivingRecords(ctx context.Context) *firestore.CollectionRef
}

func ProvideProductionCollectionProvider(client *firestore.Client) CollectionProvider {
	return &productionCollectionProvider{client: client}
}

type productionCollectionProvider struct {
	client *firestore.Client
}

var _ CollectionProvider = (*productionCollectionProvider)(nil)

func (p *productionCollectionProvider) DrivingRecords(_ context.Context) *firestore.CollectionRef {
	return p.client.Collection("driving_records")
}

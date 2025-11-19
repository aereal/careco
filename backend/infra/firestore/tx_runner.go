package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type TransactionRunner interface {
	RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) error
}

var _ TransactionRunner = (*firestore.Client)(nil)

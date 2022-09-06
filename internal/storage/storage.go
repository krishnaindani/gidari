package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpine-hodler/gidari/pkg/proto"
)

const (
	// MongoType is the byte representation of a mongo database.
	MongoType uint8 = iota

	// PostgresType is the byte representation of a postgres database.
	PostgresType
)

// Errors
var (
	ErrTransactionInProgress = fmt.Errorf("transaction is already in progress")
)

// Storage is an interface that defines the methods that a storage device should implement.
type Storage interface {
	Close()
	Read(context.Context, *proto.ReadRequest, *proto.ReadResponse) error

	// StartTx will start a transaction and return a "Tx" object that can be used to put operations on a channel,
	// commit the result of all operations sent to the transaction, or rollback the result of all operations sent
	// to the transaction.
	StartTx(context.Context) Tx
	TruncateTables(context.Context, *proto.TruncateTablesRequest) error
	Upsert(context.Context, *proto.UpsertRequest, *proto.CreateResponse) error
	Type() uint8
}

// Scheme takes a byte and returns the associated DNS root database resource.
func Scheme(t uint8) string {
	switch t {
	case MongoType:
		return "mongodb"
	case PostgresType:
		return "postgresql"
	default:
		return "unknown"
	}
}

// New will attempt to return a generic storage object given a DNS.
func New(ctx context.Context, dns string) (Storage, error) {
	if strings.Contains(dns, Scheme(MongoType)) {
		return NewMongo(ctx, dns)
	}

	if strings.Contains(dns, Scheme(PostgresType)) {
		return NewPostgres(ctx, dns)
	}

	return nil, fmt.Errorf("databse for dns %q is not supported", dns)
}

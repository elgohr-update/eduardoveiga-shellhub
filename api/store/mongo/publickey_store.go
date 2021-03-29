package mongo

import (
	"context"
	"strings"

	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Store) GetPublicKey(ctx context.Context, fingerprint, tenant string) (*models.PublicKey, error) {
	pubKey := new(models.PublicKey)
	if tenant != "" {
		if err := s.db.Collection("public_keys").FindOne(ctx, bson.M{"fingerprint": fingerprint, "tenant_id": tenant}).Decode(&pubKey); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, store.ErrRecordNotFound
			}

			return nil, err
		}
	} else {
		if err := s.db.Collection("public_keys").FindOne(ctx, bson.M{"fingerprint": fingerprint}).Decode(&pubKey); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, store.ErrRecordNotFound
			}

			return nil, err
		}
	}

	return pubKey, nil
}

func (s *Store) ListPublicKeys(ctx context.Context, pagination paginator.Query) ([]models.PublicKey, int, error) {
	query := []bson.M{
		{
			"$sort": bson.M{
				"created_at": 1,
			},
		},
	}

	// Only match for the respective tenant if requested
	if tenant := apicontext.TenantFromContext(ctx); tenant != nil {
		query = append(query, bson.M{
			"$match": bson.M{
				"tenant_id": tenant.ID,
			},
		})
	}

	queryCount := append(query, bson.M{"$count": "count"})
	count, err := aggregateCount(ctx, s.db.Collection("public_keys"), queryCount)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, buildPaginationQuery(pagination)...)

	list := make([]models.PublicKey, 0)
	cursor, err := s.db.Collection("public_keys").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		key := new(models.PublicKey)
		err = cursor.Decode(&key)
		if err != nil {
			return list, count, err
		}

		list = append(list, *key)
	}

	return list, count, err
}

func (s *Store) CreatePublicKey(ctx context.Context, key *models.PublicKey) error {
	if err := key.Validate(); err != nil {
		return err
	}

	_, err := s.db.Collection("public_keys").InsertOne(ctx, key)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			return store.ErrDuplicateFingerprint
		}
	}

	return err
}

func (s *Store) UpdatePublicKey(ctx context.Context, fingerprint, tenant string, key *models.PublicKeyUpdate) (*models.PublicKey, error) {
	if err := key.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.db.Collection("public_keys").UpdateOne(ctx, bson.M{"fingerprint": fingerprint}, bson.M{"$set": key}); err != nil {
		if err != nil {
			if strings.Contains(err.Error(), "public key not found") {
				return nil, store.ErrRecordNotFound
			}
		}

		return nil, err
	}

	return s.GetPublicKey(ctx, fingerprint, tenant)
}

func (s *Store) DeletePublicKey(ctx context.Context, fingerprint, tenant string) error {
	_, err := s.db.Collection("public_keys").DeleteOne(ctx, bson.M{"fingerprint": fingerprint, "tenant_id": tenant})
	return err
}

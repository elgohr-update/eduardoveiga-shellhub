package migrations

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var migration22 = migrate.Migration{
	Version:     22,
	Description: "Insert the user on the members group for the namespace",
	Up: func(db *mongo.Database) error {
		logrus.WithFields(logrus.Fields{
			"component": "migration",
			"version":   22,
			"action":    "Up",
		}).Info("Applying migration")
		cursor, err := db.Collection("namespaces").Find(context.TODO(), bson.D{})
		if err != nil {
			return err
		}
		for cursor.Next(context.TODO()) {
			namespace := new(models.Namespace)
			err = cursor.Decode(&namespace)
			if err != nil {
				return err
			}
			for _, memberID := range namespace.Members {
				user := new(models.User)
				objID, err := primitive.ObjectIDFromHex(memberID.(string))
				if err != nil {
					return err
				}
				if err := db.Collection("users").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user); err != nil {
					if _, err := db.Collection("namespaces").UpdateOne(context.TODO(), bson.M{"tenant_id": namespace.TenantID}, bson.M{"$pull": bson.M{"members": memberID}}); err != nil {
						return err
					}
				}
			}
		}

		return nil
	},
	Down: func(db *mongo.Database) error {
		logrus.WithFields(logrus.Fields{
			"component": "migration",
			"version":   22,
			"action":    "Down",
		}).Info("Applying migration")

		return nil
	},
}

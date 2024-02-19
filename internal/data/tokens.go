package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	GUID          string
	RefreshToken  string
	RefreshExpiry time.Time
}

type TokenModel struct {
	DB *mongo.Client
}

func (m TokenModel) Insert(token *Token) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	collection := m.DB.Database("medods").Collection("tokens")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"guid":          token.GUID,
		"hashedToken":   hashedToken,
		"refreshExpiry": token.RefreshExpiry,
	}

	update := bson.M{
		"$set": doc,
	}

	opts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(ctx, bson.M{"guid": token.GUID}, update, opts)
	if err != nil {
		return err
	}

	return nil
}

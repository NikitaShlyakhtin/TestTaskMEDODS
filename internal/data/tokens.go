package data

import (
	"context"
	"medods/internal/validator"
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
		"hashedToken":   string(hashedToken),
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

func (m TokenModel) Find(refreshToken string) (string, error) {
	collection := m.DB.Database("medods").Collection("tokens")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		GUID          string    `bson:"guid"`
		HashedToken   string    `bson:"hashedToken"`
		RefreshExpiry time.Time `bson:"refreshExpiry"`
	}

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return "", err
	}
	defer cur.Close(ctx)

	tokenFound := false

	for cur.Next(ctx) {
		err := cur.Decode(&result)
		if err != nil {
			return "", err
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.HashedToken), []byte(refreshToken))
		if err == nil {
			tokenFound = true
			break
		}
	}

	if err := cur.Err(); err != nil {
		return "", nil
	}

	if time.Now().After(result.RefreshExpiry) {
		return "", ErrRefreshTokenExpired
	}

	if tokenFound {
		return result.GUID, nil
	}

	return "", ErrRecordNotFound
}

func ValidateGUID(v *validator.Validator, guid string) {
	v.Check(guid != "", "guid", "must be provided")
	v.Check(validator.Matches(guid, validator.GUIDRX), "guid", "must be in the correct format")
}

func ValidateRefreshToken(v *validator.Validator, token string) {
	v.Check(token != "", "refresh_token", "must be provided")
	v.Check(len(token) == 44, "refresh_token", "must be 44 bytes long")
}

package data

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Tokens TokenModel
}

func NewModels(db *mongo.Client) Models {
	return Models{
		Tokens: TokenModel{DB: db},
	}
}

package model

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrNotFound        = errors.New("link not found")
	ErrIdAlreadyExists = errors.New("identifier already exists")
	ErrDataBaseErr     = errors.New("database error")
)

type Shortened struct {
	Id         string    `json:"id"`
	Url        string    `json:"url"`
	CreatedBy  string    `json:"createdBy"`
	VisitCount uint64    `json:"visitCount"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type MongoShortenUrl struct {
	Id         string    `bson:"_id"`
	CreatedBy  string    `bson:"createdBy"`
	Url        string    `bson:"ulr"`
	VisitCount uint64    `bson:"visitCount"`
	CreatedAt  time.Time `bson:"createdAt"`
	UpdatedAt  time.Time `bson:"updatedAt"`
}

type ShortenedReq struct {
	Url       string
	Id        string
	CreatedBy string
}

type GhInfo struct {
	GHAccessKey string `json:"acsessKey,omitempty"`
	GHLogin     string `json:"login,omitempty"`
}

type User struct {
	IsActive  bool      `json:"isAuthorized,omitempty"`
	GhInfo    GhInfo    `json:"githubInfo,omitempty"`
	CreatedAt time.Time `bson:"createdAt"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	User `json:"user_data"`
}

func EntToMongo(shortened Shortened) MongoShortenUrl {
	return MongoShortenUrl{
		Id:         shortened.Id,
		CreatedBy:  shortened.CreatedBy,
		Url:        shortened.Url,
		VisitCount: shortened.VisitCount,
		CreatedAt:  shortened.CreatedAt,
		UpdatedAt:  shortened.UpdatedAt,
	}
}

func MongoToEnt(shortened MongoShortenUrl) *Shortened {
	return &Shortened{
		Id:         shortened.Id,
		Url:        shortened.Url,
		CreatedBy:  shortened.CreatedBy,
		VisitCount: shortened.VisitCount,
		CreatedAt:  shortened.CreatedAt,
		UpdatedAt:  shortened.UpdatedAt,
	}
}

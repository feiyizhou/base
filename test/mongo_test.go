package test

import (
	"context"
	"testing"

	"github.com/feiyizhou/base/clients"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func initMongoDB() {
	db = clients.NewMongoDB(context.Background(), clients.MongoConf{
		Host:           "43.154.139.231:27017",
		Database:       "agent_platform",
		Username:       "feiyizhou",
		Password:       "feiyizhou7816",
		ConnectTimeout: 15,
		ReplSetName:    "rs0",
	})
}

func Test_Mongo(t *testing.T) {
	initMongoDB()
}

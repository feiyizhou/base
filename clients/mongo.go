package clients

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConf struct {
	Host           string `json:"host" mapstructure:"host"`
	Database       string `json:"database" mapstructure:"database"`
	Username       string `json:"username" mapstructure:"username"`
	Password       string `json:"password" mapstructure:"password"`
	ConnectTimeout int    `json:"connectTimeout" mapstructure:"connectTimeout"`
	ReplSetName    string `json:"replSetName" mapstructure:"replSetName"`
}

type MongoClient struct {
	ctx context.Context
	DB  *mongo.Database
}

func NewMongoClient(ctx context.Context, conf MongoConf) (*MongoClient, error) {
	clientOpts := options.Client().
		SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    conf.Database,
			Username:      conf.Username,
			Password:      conf.Password,
		}).
		SetConnectTimeout(time.Duration(conf.ConnectTimeout) * time.Second).
		SetHosts([]string{conf.Host}).
		SetMaxPoolSize(uint64(4 * runtime.NumCPU())).
		SetMinPoolSize(uint64(runtime.NumCPU())).
		SetReadPreference(readpref.Primary()).
		SetDirect(true)
	if len(conf.ReplSetName) != 0 {
		clientOpts.SetReplicaSet(conf.ReplSetName)
	}

	cli, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connect to mongo failed, err: %v", err))
	}
	err = cli.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Ping mongo failed, err: %v", err))
	}
	return &MongoClient{ctx: ctx, DB: cli.Database(conf.Database)}, err
}

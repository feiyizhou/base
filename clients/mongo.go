package clients

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/event"
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

func NewMongoDB(ctx context.Context, conf MongoConf, monitor bool) *mongo.Database {
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
	if monitor {
		clientOpts.SetMonitor(&event.CommandMonitor{
			Started: func(ctx context.Context, e *event.CommandStartedEvent) {
				log.Printf("MongoDB operate command: \ncommand type: %s \nfull command: %s \ndatabase: %s\n",
					e.CommandName,
					e.Command.String(),
					e.DatabaseName,
				)
			},
			Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {
				log.Printf("operate success, time consuming: %d milliseconds\n", e.Duration.Milliseconds())
			},
			Failed: func(ctx context.Context, e *event.CommandFailedEvent) {
				log.Printf("operate failed, err: %s\n", e.Failure)
			},
		})
	}
	if len(conf.ReplSetName) != 0 {
		clientOpts.SetReplicaSet(conf.ReplSetName)
	}

	cli, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(fmt.Errorf("connect to mongo failed, err: %v", err))
	}
	err = cli.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Errorf("ping mongo failed, err: %v", err))
	}
	return cli.Database(conf.Database)
}

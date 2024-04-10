package clients

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

type ESConf struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
}

type ESClient struct {
	ctx    context.Context
	Client *elastic.Client
}

func NewESClient(ctx context.Context, conf ESConf) (*ESClient, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetHealthcheckInterval(10 * time.Second),
		elastic.SetSniff(false),
		elastic.SetURL(fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)),
		elastic.SetGzip(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetBasicAuth(conf.User, conf.Password),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &ESClient{ctx: ctx, Client: client}, err
}

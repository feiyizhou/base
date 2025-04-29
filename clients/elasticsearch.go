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

func NewESClient(conf ESConf) *elastic.Client {
	url := fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)
	opts := []elastic.ClientOptionFunc{
		elastic.SetHealthcheckInterval(10 * time.Second),
		elastic.SetSniff(false),
		elastic.SetURL(url),
		elastic.SetGzip(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetBasicAuth(conf.User, conf.Password),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(url).Do(context.Background())
	if err != nil {
		panic(fmt.Errorf("error pinging Elasticsearch: %v", err))
	}
	fmt.Printf("elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return client
}

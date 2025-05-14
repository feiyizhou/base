package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/feiyizhou/base/clients"
	"github.com/feiyizhou/base/utils"
	"github.com/olivere/elastic/v7"
)

var (
	mapping   string
	indexName = "my_index"
	es        *elastic.Client
)

// Document 文档结构
type Document struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Authors  []string `json:"authors"`
	Year     int      `json:"year"`
	Abstract string   `json:"abstract"`
}

// IndexMapping 索引结构
type IndexMapping struct {
	Settings IndexSettings `json:"settings"`
	Mappings `json:"mappings"`
}

type IndexSettings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type Mappings struct {
	Properties map[string]Properties `json:"properties"`
}

type Properties struct {
	Type     string `json:"type"`
	Analyzer string `json:"analyzer,omitempty"`
}

func init() {
	es = clients.NewESClient(clients.ESConf{
		Host:     "172.16.200.18",
		Port:     31172,
		User:     "admin",
		Password: "Sobey123",
	})
	mappingBytes, _ := json.Marshal(IndexMapping{
		Settings: IndexSettings{
			NumberOfShards:   1,
			NumberOfReplicas: 0,
		},
		Mappings: Mappings{
			Properties: map[string]Properties{
				"id":       {Type: "integer"},
				"title":    {Type: "text", Analyzer: "english"},
				"authors":  {Type: "keyword"},
				"year":     {Type: "integer"},
				"abstract": {Type: "text", Analyzer: "english"},
			},
		},
	})
	mapping = string(mappingBytes)
}

func Test_ES_Create_Idx(t *testing.T) {
	exists, err := es.IndexExists(indexName).Do(context.Background())
	if err != nil {
		fmt.Printf("Error checking if index exists: %v", err)
		return
	}
	if exists {
		fmt.Println("Index verified to exist")
		return
	} else {
		createIndex, err := es.CreateIndex(indexName).Body(mapping).Do(context.Background())
		if err != nil {
			fmt.Printf("Error creating index: %v", err)
			return
		}
		if !createIndex.Acknowledged {
			fmt.Printf("Create index not acknowledged")
			return
		}
		fmt.Printf("Successfully created index: %s\n", indexName)
	}
}

func Test_ES_Delete_Idx(t *testing.T) {
	exists, err := es.IndexExists(indexName).Do(context.Background())
	if err != nil {
		fmt.Printf("Error checking if index exists: %v", err)
		return
	}
	if exists {
		_, err := es.DeleteIndex(indexName).Do(context.Background())
		if err != nil {
			fmt.Printf("Error deleting existing index: %v", err)
			return
		}
		fmt.Printf("Deleted existing index: %s", indexName)
	}
}

func Test_ES_All_Idx(t *testing.T) {
	indices, err := es.CatIndices().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting indices: %v", err)
		return
	}
	for _, index := range indices {
		fmt.Printf("Index: %s, Docs: %d, Size: %s\n",
			index.Index, index.DocsCount, index.StoreSize)
	}
}

func Test_ES_Create_Data(t *testing.T) {
	for i := range 100 {
		resp, err := es.Index().Index(indexName).BodyJson(Document{
			ID:       i,
			Title:    fmt.Sprintf("wangmazi-%d", i),
			Authors:  []string{utils.RandUUIDStr()},
			Year:     i,
			Abstract: fmt.Sprintf("张三李四王麻子-%d", i),
		}).Do(context.Background())
		if err != nil {
			fmt.Printf("error indexing document: %v", err)
			return
		}
		fmt.Println(resp.Id)
		time.Sleep(10 * time.Second)
	}
}

func Test_ES_Query_By_ID(t *testing.T) {
	getResp, err := es.Get().
		Index(indexName).
		Id("aF76epYBet1VMiAEs7bV").
		Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			fmt.Println("document not found")
			return
		}
		fmt.Printf("error getting document: %v", err)
		return
	}
	var document Document
	if err := json.Unmarshal(getResp.Source, &document); err != nil {
		fmt.Printf("error parsing document: %v", err)
		return
	}
	fmt.Println(document)
}

func Test_ES_List(t *testing.T) {
	query := elastic.NewBoolQuery()
	// query.Must(elastic.NewMatchQuery("title", "zhangsan"))
	searchResult, err := es.Search().
		Index(indexName).
		Query(query).
		Sort("year", true).
		From(0).Size(10).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		fmt.Printf("error executing search: %v", err)
		return
	}
	var documents []Document
	for _, hit := range searchResult.Hits.Hits {
		var d Document
		fmt.Printf("es id: %s\n", hit.Id)
		if err := json.Unmarshal(hit.Source, &d); err != nil {
			fmt.Printf("error parsing product %s: %v", hit.Id, err)
			continue
		}
		documents = append(documents, d)
	}
	fmt.Println(documents)
}

func Test_ES_Delete_By_Id(t *testing.T) {
	resp, err := es.Delete().
		Index(indexName).
		Id("Zl75epYBet1VMiAEY7aR").
		Do(context.Background())
	if err != nil {
		fmt.Printf("error deleting document: %v", err)
		return
	}
	fmt.Println(resp.Id)
}

func Test_ES_Delete_By_Query(t *testing.T) {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewMatchQuery("title", "zhangsan"))
	resp, err := es.DeleteByQuery().
		Index(indexName).
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		fmt.Printf("error deleting by query: %v", err)
		return
	}
	fmt.Println(resp.Deleted)
}

package elasticsearch_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	elasticsearchInit           sync.Once
	elasticSearchClientInstance *ElasticSearchClient
)

type ElasticSearchClient struct {

	// es client
	config *elasticsearch.Config
	client *elasticsearch.Client
}

type Client interface {
	BulkToEs(records [][]string)
	SearchToEs(data interface{})
}

/*
*
bulk insert into es
*/
func BulkToEs(records [][]string, index string, workers int) {
	client := elasticSearchClientInstance.client
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         index,            // The default index name
		Client:        client,           // The Elasticsearch client
		NumWorkers:    workers,          // The number of worker goroutines
		FlushBytes:    int(5e+6),        // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		panic(err)
	}
	// loop in save
	poiDocument := parsingRecords(records)
	CreateIndex(index)
	//bulk insert into es
	bulkInsertInfo(poiDocument, bulkIndexer)

}

func parsingRecords(records [][]string) []*PoiDocument {
	var poiDocuments []*PoiDocument
	for i, record := range records {
		if len(record) < 6 {
			log.Printf("Skipping invalid record at line %d: %v", i+1, record)
			continue
		}
		if i == 0 {
			continue
		}
		lat, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			fmt.Printf("Error converting latitude to float: %v", err)
			continue
		}

		lon, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			fmt.Printf("Error converting latitude to float: %v", err)
			continue
		}
		lon, lat = ConvertTMToWGS84(lat, lon)
		poiDocuments = append(poiDocuments, &PoiDocument{
			record[2],
			record[3],
			lat,
			lon,
			[]float64{lon, lat},
			record[7],
		})
	}
	return poiDocuments
}

func CreateIndex(index string) {
	indexReq := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(getRawMapping()),
	}
	res, err := indexReq.Do(context.Background(), elasticSearchClientInstance.client)
	if err != nil {
		fmt.Errorf("exception while creating index : %s", err.Error())
	}
	fmt.Printf("Created index %s\n %s", index, res)
	res.Body.Close()
}

func getRawMapping() string {
	mapping := `{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
      },
      "mappings": {
        "properties": {
          "coordinates": {
            "type": "geo_point"
          }
        }
}}`
	return mapping
}

func bulkInsertInfo(documents []*PoiDocument, bi esutil.BulkIndexer) {
	var countSuccessful uint64
	for _, document := range documents {
		data, err := json.Marshal(document)
		log.Printf("origin : %v === data : %s", document, data)
		if err != nil {
			fmt.Errorf("Error converting documents to json: %v", err)
		}
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: document.Id,
				Body:       bytes.NewReader(data),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
	}
	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// Close the indexer
	//
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}

// elasticsearch client 생성 메서드
func NewElasticSearchClient() *ElasticSearchClient {

	elasticsearchInit.Do(func() {
		elasticSearchClientInstance = &ElasticSearchClient{}
	})
	config := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		},
	}
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		fmt.Errorf("Error creating elasticsearch client: %v", err)
	}
	elasticSearchClientInstance.client = client
	return elasticSearchClientInstance
}

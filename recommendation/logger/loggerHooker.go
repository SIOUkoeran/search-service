package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	ElasticsearchHookerInit     sync.Once
	elasticsearchHookerInstance ElasticsearchHooker
)

type (
	ElasticsearchHooker struct {
		client    *elasticsearch.Client
		indexName string
	}
)

// @Summary 로깅 후킹 객체 생성 메서드
func NewElasticsearchHook(client *elasticsearch.Client, index string) ElasticsearchHooker {
	ElasticsearchHookerInit.Do(func() {
		elasticsearchHookerInstance = ElasticsearchHooker{
			client:    client,
			indexName: index,
		}
	})
	return elasticsearchHookerInstance
}

// Fire 메서드는 로그 엔트리가 발생할 때마다 호출되어 Elasticsearch에 적재
func (h ElasticsearchHooker) Fire(entry *logrus.Entry) error {
	// 로그 엔트리를 JSON으로 직렬화

	// 문서에 기본 필드 추가(타임스탬프, 레벨, 메시지 등)
	doc := map[string]interface{}{
		"@timestamp": entry.Time.UTC().Format(time.RFC3339),
		"level":      entry.Level.String(),
		"message":    entry.Message,
		"fields":     entry.Data,
	}
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Elasticsearch에 인덱싱 요청
	req := esapi.IndexRequest{
		Index:      elasticsearchHookerInstance.indexName,
		DocumentID: "", // ES가 자동으로 ID 생성
		Body:       bytes.NewReader(docBytes),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), elasticsearchHookerInstance.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing log entry: %s", res.String())
	}

	return nil
}

// Levels 메서드로 어떤 레벨의 로그를 Elasticsearch에 보낼지 결정
func (h ElasticsearchHooker) Levels() []logrus.Level {
	return logrus.AllLevels
}

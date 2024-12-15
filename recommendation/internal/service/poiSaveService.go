package service

import (
	elasticsearch_client "recommendation/internal/elasticsearch-client"
	"sync"
)

var (
	PoiSaveServiceInit     sync.Once
	PoiSaveServiceInstance *PoiSaveService
)

type PoiSaveService struct {
	elasticRepository *elasticsearch_client.ElasticSearchClient
}

func NewPoiSaveService(elasticRepository *elasticsearch_client.ElasticSearchClient) *PoiSaveService {
	PoiSaveServiceInit.Do(func() {
		PoiSaveServiceInstance = &PoiSaveService{
			elasticRepository: elasticRepository,
		}
	})
	return PoiSaveServiceInstance
}

func SavePoiIntoEs() {

}

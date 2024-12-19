package api

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	v1 "recommendation/api/v1"
	"recommendation/infrastructure"
	"recommendation/repository"
	"recommendation/service"
	"recommendation/setting"
)

func SetRouters(r *gin.Engine, cfg *setting.Configuration, client *elasticsearch.Client) {
	poi := createPoiDomain(cfg, client)
	apiv1 := r.Group("/api/v1")
	{
		// Create artile with tags
		// Get aritle detail by id
		apiv1.POST("/poi", poi.SearchPoi)
		// Update poi detail by id
		apiv1.PUT("/poi/:id", poi.UpdatePoi)
		// Delete poi by id
		apiv1.DELETE("/poi/:id", poi.DeletePoi)

		// Add other router if necessary
	}
}

// @Summary POI 전체 도메인 생성 내부 메서드
func createPoiDomain(cfg *setting.Configuration, client *elasticsearch.Client) v1.PoiController {
	modelServer := infrastructure.NewModelServerCaller()
	repository := repository.NewPoiRepository(client)
	service := service.NewSearchService(repository, modelServer)
	poiRouter := v1.NewPoiRouter(service, cfg)
	return poiRouter
}

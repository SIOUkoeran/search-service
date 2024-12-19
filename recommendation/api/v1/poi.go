package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"recommendation/domain"
	"recommendation/dto"
	"recommendation/logger"
	"recommendation/service"
	"recommendation/setting"
	"sync"
)

var (
	poiRouterInit     sync.Once
	poiRouterInstance *poiController
)

type (
	PoiController interface {
		SearchPoi(ctx *gin.Context)
		CreatePoi(ctx *gin.Context)
		UpdatePoi(ctx *gin.Context)
		DeletePoi(ctx *gin.Context)
	}

	poiController struct {
		searchService service.SearchService
		conf          *setting.Configuration
	}
)

func (p poiController) SearchPoi(ctx *gin.Context) {
	var searchRequest dto.TitleSearchRequest
	if err := ctx.ShouldBindJSON(&searchRequest); err != nil {
		res := domain.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	req := ctx.Request.Context()
	logger.WithTrace(ctx).Info("request In : ", searchRequest)
	response := p.searchService.SearchTitle(req, searchRequest.Title)
	if len(response) == 0 {
		res := domain.BuildResponseSuccess("success", response)
		ctx.JSON(http.StatusOK, res)
		return
	}
	errRes := domain.BuildResponseFailed("failed", "failed", response)
	ctx.JSON(http.StatusOK, errRes)
}

func (p poiController) CreatePoi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p poiController) UpdatePoi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p poiController) DeletePoi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewPoiRouter(service service.SearchService, conf *setting.Configuration) PoiController {
	poiRouterInit.Do(func() {
		poiRouterInstance = &poiController{
			searchService: service,
			conf:          conf,
		}
	})
	return poiRouterInstance
}

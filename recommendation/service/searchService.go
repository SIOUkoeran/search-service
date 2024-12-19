package service

import (
	"context"
	"recommendation/dto"
	"recommendation/infrastructure"
	"recommendation/repository"
	"sync"
)

var (
	searchServiceInit     sync.Once
	searchServiceInstance *searchService
)

type (
	SearchService interface {
		SearchTitle(ctx context.Context, title string) []dto.PoiEntity
	}

	searchService struct {
		repository      repository.PoiRepository
		embeddingServer infrastructure.ModelServiceCaller
	}
)

/*
*
SearachService 생성 메서드
*/
func NewSearchService(repository repository.PoiRepository, embeddingServer infrastructure.ModelServiceCaller) SearchService {
	searchServiceInit.Do(func() {
		searchServiceInstance = &searchService{
			repository:      repository,
			embeddingServer: embeddingServer,
		}
	})
	return searchServiceInstance
}

/*
*
 */
func (s *searchService) SearchTitle(c context.Context, title string) []dto.PoiEntity {
	embedResponse := s.getVector(c, title)
	return s.repository.SearchPoiByTitle(c, embedResponse.Vector, "vector_poi")
}

func (s *searchService) getVector(c context.Context, title string) infrastructure.ResponseEmbedQuery {
	return s.embeddingServer.CallEmbedQuery(
		&infrastructure.RequestEmbedQuery{
			Query: title,
		})
}

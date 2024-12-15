package service

import (
	"context"
	"recommendation/dto"
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
		repository repository.PoiRepository
	}
)

/*
*
SearachService 생성 메서드
*/
func NewSearchService(repository repository.PoiRepository) SearchService {
	searchServiceInit.Do(func() {
		searchServiceInstance = &searchService{
			repository: repository,
		}
	})
	return searchServiceInstance
}

/*
*
 */
func (s *searchService) SearchTitle(c context.Context, title string) []dto.PoiEntity {
	return s.repository.SearchPoiByTitle(c, title, "test_poi")
}

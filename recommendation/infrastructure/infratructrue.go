package infrastructure

import (
	"recommendation/setting"
	"sync"
)

var (
	infrastructureInit sync.Once
)

// @Summary 외부 호출 설정 생성 메서드
func SetInfrastructure(cfg *setting.Configuration) {
	infrastructureInit.Do(func() {
		// model server caller init
		NewModelServerCaller()
	})
}

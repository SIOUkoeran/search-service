package elasticsearch_client

import "math"

const (
	a    = 6378137.0               // 장축 반경 (GRS80)
	f    = 1 / 298.257222101       // 편평률
	lat0 = 38.0 * math.Pi / 180.0  // 기준 위도 (중부원점: 38도)
	lon0 = 127.0 * math.Pi / 180.0 // 기준 경도 (중부원점: 127도)
	k0   = 1.0                     // 축척 계수
	x0   = 200000.0                // X 좌표 원점
	y0   = 500000.0                // Y 좌표 원점
)

func ConvertTMToWGS84(x, y float64) (float64, float64) {

	_ = math.Sqrt(2*f - f*f)
	n := f / (2 - f)
	A := a / (1 + n) * (1 + n*n/4 + n*n*n*n/64)

	// (x, y)에서 원점 보정
	x -= x0
	y -= y0

	// 위도 경도 계산
	lat := lat0 // 초기값 설정
	for i := 0; i < 5; i++ {
		lat = (y / (k0 * A)) + lat0
	}
	lon := lon0 + (x / (k0 * A * math.Cos(lat)))

	// 라디안을 도로 변환
	latDeg := lat * 180.0 / math.Pi
	lonDeg := lon * 180.0 / math.Pi

	return lonDeg, latDeg
}

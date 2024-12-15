package service

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"recommendation/internal/elasticsearch-client"
	"recommendation/internal/redis-client"
	"strconv"
)

func SaveRecordsToRedis(records [][]string, rdb *redis.Client) {
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
			fmt.Printf("Error converting longitude to float: %v", err)
			continue
		}
		lon, lat = elasticsearch_client.ConvertTMToWGS84(lat, lon)
		err = rdb.GeoAdd(redis_client.Ctx, "seoul", &redis.GeoLocation{
			Name:      record[7],
			Longitude: lon,
			Latitude:  lat,
		}).Err()
		if err != nil {
			log.Printf("error lat : %f, lon : %f", lat, lon)
			log.Printf("Failed to save record at line %d: %v", i+1, record)
			log.Printf("error :%s", err)
		} else {
			log.Printf("Successfully to save record at line %d: %v", i+1, record)
		}
	}
}

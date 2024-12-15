package main

import (
	"context"
	"log"
	"recommendation/cmd"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

var ctx = context.Background()

func main() {

	//records, err := csvreader.ReadCsv("./data/poi_data.csv")
	//if err != nil {
	//	log.Fatalf("Failed to read CSV file: %v", err)
	//}
	//log.Printf("Read %d records", len(records))
	//
	//// Redis 클라이언트 설정
	//elasticsearch_client.NewElasticSearchClient()
	//elasticsearch_client.BulkToEs(records, "test_poi", 5)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

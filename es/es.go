package es

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/olivere/elastic.v5"
)

func CreateNewElasticSearchClient(ctx context.Context, url *string, sniff *bool) *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(*url), elastic.SetSniff(*sniff), elastic.SetTraceLog(log.New(os.Stdout, "es: ", 0)))
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	info, code, err := client.Ping(*url).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion(*url)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	return client
}

func InitializeNewIndex(client *elastic.Client, ctx context.Context, indexName string, mapping string) {
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	if exists {
		_, err := client.DeleteIndex(indexName).Do(ctx)
		if err != nil {
			panic(err)
		}
	}

	createIndex, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !createIndex.Acknowledged {
		panic(errors.New(fmt.Sprintf("%s Index Not Created", indexName)))
	}
}

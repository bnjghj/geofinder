package es

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/olivere/elastic.v5"
)

func CreateNewElasticSearchClient(ctx context.Context, url *string, sniff *bool, trace *bool) *elastic.Client {
	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(*url))
	options = append(options, elastic.SetSniff(*sniff))
	if *trace {
		options = append(options, elastic.SetTraceLog(log.New(os.Stdout, "es: ", 0)))
	}
	client, err := elastic.NewClient(options...)
	if err != nil {
		fmt.Printf("Elasticsearch Url %s\n", *url)
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

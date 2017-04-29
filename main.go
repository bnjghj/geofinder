package main

import (
	"flag"
	"geofinder/es"
	"geofinder/handler"
	"geofinder/model"
	"log"
	"os"

	"golang.org/x/net/context"

	"github.com/valyala/fasthttp"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetOutput(os.Stdout)
}

func main() {
	ctx := context.Background()

	var (
		url      = flag.String("url", "http://localhost:9200", "Elasticsearch URL")
		sniff    = flag.Bool("sniff", false, "Enable or disable sniffing")
		addr     = flag.String("addr", ":8080", "TCP address to listen to")
		compress = flag.Bool("compress", true, "Whether to enable transparent response compression")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}

	client := es.CreateNewElasticSearchClient(ctx, url, sniff)

	es.InitializeNewIndex(client, ctx, model.PolygonIndexConfiguration.IndexName, model.PolygonIndexConfiguration.Mapping)
	es.LoadPolygonIndex(client, ctx)
	es.TryPolygonIndex(client, ctx)

	myHandler := &handler.GeoQueryHandler{
		ElasticClient: client,
		Context:       ctx,
	}

	h := myHandler.HandleFastHTTP
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
		panic(err)
	}
}

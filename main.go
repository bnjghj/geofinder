package main

import (
	"flag"
	"geofinder/es"
	"geofinder/handler"
	"geofinder/model"
	"log"
	"os"

	"golang.org/x/net/context"

	"fmt"

	"strconv"

	"github.com/caarlos0/env"
	"github.com/valyala/fasthttp"
)

type mainConfiguration struct {
	URL      string `env:"ELASTICSEARCH_URL" envDefault:"http://localhost:9200"`
	Sniff    bool   `env:"ELASTICSEARCH_SNIFF" envDefault:"false"`
	Port     int    `env:"PORT" envDefault:"8080"`
	Compress bool   `env:"COMPRESS" envDefault:"true"`
	Trace    bool   `env:"ELASTICSEARCH_TRACE" envDefault:"false"`
}

var cfg mainConfiguration

func init() {
	cfg = mainConfiguration{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	log.SetOutput(os.Stderr)
	log.SetOutput(os.Stdout)
}

func main() {
	ctx := context.Background()

	var (
		url      = flag.String("url", cfg.URL, "Elasticsearch URL")
		sniff    = flag.Bool("sniff", cfg.Sniff, "Enable or disable sniffing")
		addr     = flag.String("addr", fmt.Sprintf(":%s", strconv.Itoa(cfg.Port)), "TCP address to listen to")
		compress = flag.Bool("compress", cfg.Compress, "Whether to enable transparent response compression")
		trace    = flag.Bool("trace", cfg.Trace, "Enable or disable elastic search tracing")
	)
	flag.Parse()
	log.SetFlags(0)

	if *url == "" {
		*url = "http://127.0.0.1:9200"
	}

	client := es.CreateNewElasticSearchClient(ctx, url, sniff, trace)

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

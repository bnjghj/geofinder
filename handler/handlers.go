package handler

import (
	"context"

	"github.com/valyala/fasthttp"
	elastic "gopkg.in/olivere/elastic.v5"
)

type GeoQueryHandler struct {
	ElasticClient *elastic.Client
	Context       context.Context
}

// request handler in net/http style, i.e. method bound to MyCustomHandler struct.
func (h *GeoQueryHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain; charset=utf8")
	switch string(ctx.Path()) {
	case "/shape":
		shapeQueryHandlerFunc(ctx, h.ElasticClient, h.Context)
	case "/point":
		pointQueryHandlerFunc(ctx, h.ElasticClient, h.Context)
	case "/circle":
		circleQueryHandlerFunc(ctx, h.ElasticClient, h.Context)
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}

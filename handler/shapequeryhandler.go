package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"geofinder/es"
	"geofinder/model"
	"reflect"

	"github.com/valyala/fasthttp"
	elastic "gopkg.in/olivere/elastic.v5"
)

func shapeQueryHandlerFunc(ctx *fasthttp.RequestCtx, client *elastic.Client, context context.Context) {
	method := string(ctx.Method())
	switch method {
	case "POST":
		ctx.Response.Header.Add("Content-Type", "application/json")
		var shapeQuery es.ESCustomShapeQuery
		err := json.Unmarshal(ctx.Request.Body(), &shapeQuery)
		if err != nil {
			ctx.Error(fmt.Sprintf("Poly Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		var rawStringQuery = elastic.NewRawStringQuery(string(ctx.Request.Body()))
		searchResult, err := client.Search().Index(model.PolygonIndexConfiguration.IndexName).Type(model.PolygonIndexConfiguration.TypeName).Query(rawStringQuery).Do(context)
		if err != nil {
			ctx.Error(fmt.Sprintf("Poly Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		var result []model.PolyRegion
		var ptyp model.PolyRegion
		for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
			if t, ok := item.(model.PolyRegion); ok {
				result = append(result, t)
			}
		}
		if len(result) <= 0 {
			ctx.Error("Not Found", fasthttp.StatusNotFound)
			return
		}
		byteJSON, err := json.Marshal(result)
		if err != nil {
			ctx.Error(fmt.Sprintf("Not Found, Err %s", err), fasthttp.StatusNotFound)
			return
		}
		ctx.SetBody(byteJSON)
		return
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
		return
	}
}

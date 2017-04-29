package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"geofinder/es"
	"geofinder/model"
	"reflect"

	gj "github.com/kpawlik/geojson"
	"github.com/valyala/fasthttp"
	elastic "gopkg.in/olivere/elastic.v5"
)

type PointQueryRequest struct {
	Coordinates gj.Coordinate `json:"coordinates"`
}

type PointQueryResult struct {
	Name string `json:"name"`
}

func pointQueryHandlerFunc(ctx *fasthttp.RequestCtx, client *elastic.Client, context context.Context) {
	method := string(ctx.Method())
	switch method {
	case "POST":
		ctx.Response.Header.Add("Content-Type", "application/json")
		var pointQuery PointQueryRequest
		err := json.Unmarshal(ctx.Request.Body(), &pointQuery)
		if err != nil {
			ctx.Error(fmt.Sprintf("Point Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		bytePointQuery, err := json.Marshal(es.CreateESPointQuery(pointQuery.Coordinates))
		if err != nil {
			ctx.Error(fmt.Sprintf("Point Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		var rawStringQuery = elastic.NewRawStringQuery(string(bytePointQuery))
		searchResult, err := client.Search().Index(model.PolygonIndexConfiguration.IndexName).Type(model.PolygonIndexConfiguration.TypeName).Query(rawStringQuery).Do(context)
		if err != nil {
			ctx.Error(fmt.Sprintf("Point Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		var listPQR []PointQueryResult
		var ptyp model.PolyRegion
		for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
			if t, ok := item.(model.PolyRegion); ok {
				listPQR = append(listPQR, PointQueryResult{Name: t.Name})
			}
		}
		if len(listPQR) <= 0 {
			ctx.Error("Not Found", fasthttp.StatusNotFound)
			return
		}
		byteJSON, err := json.Marshal(listPQR)
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

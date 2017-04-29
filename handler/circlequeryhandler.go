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

type CircleQueryRequest struct {
	Radius      string        `json:"radius"`
	Coordinates gj.Coordinate `json:"coordinates"`
}

type CircleQueryResult struct {
	Name string `json:"name"`
}

func circleQueryHandlerFunc(ctx *fasthttp.RequestCtx, client *elastic.Client, context context.Context) {
	method := string(ctx.Method())
	switch method {
	case "POST":
		ctx.Response.Header.Add("Content-Type", "application/json")
		var circleQuery CircleQueryRequest
		err := json.Unmarshal(ctx.Request.Body(), &circleQuery)
		if err != nil {
			ctx.Error(fmt.Sprintf("Circle Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		byteCircleQuery, err := json.Marshal(es.CreateESCircleQuery(circleQuery.Radius, circleQuery.Coordinates))
		if err != nil {
			ctx.Error(fmt.Sprintf("Circle Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		var rawStringQuery = elastic.NewRawStringQuery(string(byteCircleQuery))
		searchResult, err := client.Search().Index(model.PolygonIndexConfiguration.IndexName).Type(model.PolygonIndexConfiguration.TypeName).Query(rawStringQuery).Do(context)
		if err != nil {
			ctx.Error(fmt.Sprintf("Circle Query Exception %s", err), fasthttp.StatusBadRequest)
			return
		}
		var listCQR []CircleQueryResult
		var ptyp model.PolyRegion
		for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
			if t, ok := item.(model.PolyRegion); ok {
				listCQR = append(listCQR, CircleQueryResult{Name: t.Name})
			}
		}
		if len(listCQR) <= 0 {
			ctx.Error("Not Found", fasthttp.StatusNotFound)
			return
		}
		byteJSON, err := json.Marshal(listCQR)
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

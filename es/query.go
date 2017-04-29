package es

import (
	gj "github.com/kpawlik/geojson"
)

const (
	typeNameCircleQuery = "circle"
	typeNamePointQuery  = "point"
)

type ESCustomShapeQuery struct {
	GeoShape ESGeoShapeQuery `json:"geo_shape"`
}

type ESGeoShapeQuery struct {
	Location ESLocationQuery `json:"location"`
}

type ESLocationQuery struct {
	Shape interface{} `json:"shape"`
}

type ESPointQuery struct {
	Coordinates gj.Coordinate `json:"coordinates"`
	Type        string        `json:"type"`
}

type ESCircleQuery struct {
	Coordinates gj.Coordinate `json:"coordinates"`
	Type        string        `json:"type"`
	Radius      string        `json:"radius"`
}

func CreateESCircleQuery(radius string, coordinate gj.Coordinate) ESCustomShapeQuery {
	return ESCustomShapeQuery{GeoShape: ESGeoShapeQuery{Location: ESLocationQuery{Shape: ESCircleQuery{Type: typeNameCircleQuery, Radius: radius, Coordinates: coordinate}}}}
}

func CreateESPointQuery(coordinate gj.Coordinate) ESCustomShapeQuery {
	return ESCustomShapeQuery{GeoShape: ESGeoShapeQuery{Location: ESLocationQuery{Shape: ESPointQuery{Type: typeNamePointQuery, Coordinates: coordinate}}}}
}

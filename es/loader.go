package es

import (
	"context"
	"encoding/json"
	"fmt"
	"geofinder/model"
	"reflect"

	gj "github.com/kpawlik/geojson"
	"gopkg.in/olivere/elastic.v5"
)

var polygonIndexName = model.PolygonIndexConfiguration.IndexName
var polygonTypeName = model.PolygonIndexConfiguration.TypeName

func CreatePolygonSearchService(client *elastic.Client) *elastic.SearchService {
	return client.Search().Index(polygonIndexName).Type(polygonTypeName)
}
func CreatePolygonIndexService(client *elastic.Client) *elastic.IndexService {
	return client.Index().Index(polygonIndexName).Type(polygonTypeName)
}

func LoadPolygonIndex(client *elastic.Client, ctx context.Context) {
	polyRegionBandirma := model.PolyRegion{Name: "Bandırma, Balıkesir", Location: model.Location{Type: "polygon", Coordinates: gj.MultiLine{gj.Coordinates{
		{40.33358, 27.97416},
		{40.33653, 27.95067},
		{40.34105, 27.932},
		{40.35088, 27.93103},
		{40.35876, 27.93528},
		{40.35045, 27.96369},
		{40.35869, 27.97055},
		{40.35503, 27.99699},
		{40.34215, 28.01047},
		{40.33083, 27.9963},
		{40.33358, 27.97416},
	}}}}
	putBandirma, err := client.Index().Index(polygonIndexName).Type(polygonTypeName).Id("1").BodyJson(polyRegionBandirma).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed polyregion %s to index %s, type %s\n", putBandirma.Id, putBandirma.Index, putBandirma.Type)

	polyRegionFatih := model.PolyRegion{Name: "Fatih, İstanbul", Location: model.Location{Type: "polygon", Coordinates: gj.MultiLine{gj.Coordinates{
		{41.00308, 28.95538},
		{40.99978, 28.93429},
		{40.98871, 28.91527},
		{41.02021, 28.91225},
		{41.02801, 28.91649},
		{41.04232, 28.93872},
		{41.02794, 28.95177},
		{41.0181, 28.97615},
		{41.01648, 28.98619},
		{41.00699, 28.98718},
		{41.00165, 28.97924},
		{41.00308, 28.95538},
	}}}}
	putFatih, err := client.Index().Index(polygonIndexName).Type(polygonTypeName).BodyJson(polyRegionFatih).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed polyregion %s to index %s, type %s\n", putFatih.Id, putFatih.Index, putFatih.Type)
}

func TryPolygonIndex(client *elastic.Client, ctx context.Context) {
	getFirst, err := client.Get().Index(polygonIndexName).Type(polygonTypeName).Id("1").Do(ctx)
	if err != nil {
		panic(err)
	}
	if getFirst.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", getFirst.Id, getFirst.Version, getFirst.Index, getFirst.Type)
	}

	// 	searchResult, err := client.Search().Index(indexName).Type(typeName).Query(elastic.NewRawStringQuery(`{
	//     "geo_shape": {
	//       "location": {
	//         "shape": {
	//           "type":   "point",
	//           "coordinates": [
	//             40.342211,
	//             27.970383
	//           ]
	//         }
	//       }
	//     }
	// }`)).Do(ctx)

	//shapeQuery := ShapeQuery{GeoShape{Location{Shape{Coordinates: {40.342211, 27.970383}, Type: "point"}}}}
	var pointQuery = CreateESPointQuery(gj.Coordinate{40.342211, 27.970383})
	bytePointQuery, err := json.Marshal(pointQuery)
	if err != nil {
		panic(err)
	}
	searchResult, err := CreatePolygonSearchService(client).Query(elastic.NewRawStringQuery(string(bytePointQuery))).Do(ctx)

	// 	searchResult, err := client.Search().Index(indexName).Type(typeName).Query(elastic.NewRawStringQuery(`{
	//     "geo_shape": {
	//       "location": {
	//         "shape": {
	//           "type":   "point",
	//           "coordinates": [
	//             40.342211,
	//             27.970383
	//           ]
	//         }
	//       }
	//     }
	// }`)).Do(ctx)

	if err != nil {
		panic(err)
	}

	var ptyp model.PolyRegion
	for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
		if t, ok := item.(model.PolyRegion); ok {
			strJSON, err := json.Marshal(t.Location)
			if err != nil {
				panic(err)
			}
			fmt.Printf("PolyRegion by %s: %s\n", t.Name, string(strJSON))
		}
	}

	circleQuery := CreateESCircleQuery("0.001km", gj.Coordinate{40.342211, 27.970383})
	byteCircleQuery, err := json.Marshal(circleQuery)
	if err != nil {
		panic(err)
	}
	searchResult, err = client.Search().Index(polygonIndexName).Type(polygonTypeName).Query(elastic.NewRawStringQuery(string(byteCircleQuery))).Do(ctx)

	for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
		if t, ok := item.(model.PolyRegion); ok {
			strJSON, err := json.Marshal(t.Location)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Circle Query -> PolyRegion by %s: %s\n", t.Name, string(strJSON))
		}
	}
}

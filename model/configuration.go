package model

const (
	indexName = "polyregion"
	typeName  = "poly"
	mapping   = `
{
	"settings": {
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings": {
    	"poly": {
      		"properties": {
        		"name": {
          			"type": "string"
        		},
        		"location": {
          			"type": "geo_shape"
        		}
      		}
    	}
  	}
}
`
)

type ESPolygonIndexConfiguration struct {
	IndexName string
	TypeName  string
	Mapping   string
}

var PolygonIndexConfiguration = ESPolygonIndexConfiguration{IndexName: indexName, TypeName: typeName, Mapping: mapping}

# GeoFinder Service GeoShape Query Examples

GeoFinder Written in GOLANG

GeoShape Query Examples For ElasticSearch

# Run Service
```sh
go run main.go
```

# Geo Point Query
```sh
POST: http://localhost:8080/point

Request :
{
	"coordinates":[40.342211,27.970383]
}

Response:
[
  {
    "name": "Bandırma, Balıkesir"
  }
]
```

# Geo Circle Query
```sh
POST: http://localhost:8080/circle

Request:
{
	"radius":"0.001km",
	"coordinates":[40.342211,27.970383]
}

Response:
[
  {
    "name": "Bandırma, Balıkesir"
  }
]
```

```sh
POST: http://localhost:8080/circle

Request:
{
	"radius":"1000km",
	"coordinates":[40.342211,27.970383]
}

Response:
[
  {
    "name": "Bandırma, Balıkesir"
  },
  {
    "name": "Fatih, İstanbul"
  }
]
```

# Raw Geo Shape Query
```sh
POST: http://localhost:8080/shape

Request:
{
	"geo_shape": {
		"location": {
			"shape": {
				"coordinates":[40.342211,27.970383],
				"radius":"0.001km",
				"type":"circle"
			}
		}
	}
}

Response:
[
  {
    "name": "Bandırma, Balıkesir",
    "location": {
      "type": "polygon",
      "coordinates": [
        [
          [
            40.33358,
            27.97416
          ],
          [
            40.33653,
            27.95067
          ],
          [
            40.34105,
            27.932
          ],
          [
            40.35088,
            27.93103
          ],
          [
            40.35876,
            27.93528
          ],
          [
            40.35045,
            27.96369
          ],
          [
            40.35869,
            27.97055
          ],
          [
            40.35503,
            27.99699
          ],
          [
            40.34215,
            28.01047
          ],
          [
            40.33083,
            27.9963
          ],
          [
            40.33358,
            27.97416
          ]
        ]
      ]
    }
  }
]
```
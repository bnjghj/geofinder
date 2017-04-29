package model

import (
	gj "github.com/kpawlik/geojson"
)

type Location struct {
	Type        string       `json:"type"`
	Coordinates gj.MultiLine `json:"coordinates"`
}

type PolyRegion struct {
	Name     string   `json:"name"`
	Location Location `json:"location,omitempty"`
}

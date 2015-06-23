package main

import (
	"encoding/json"
)

// ElasticSearchResponse is the main elasticsearch response object.
type ElasticSearchResponse struct {
	Took     int64
	TimedOut bool
	Shards   struct {
		Total      int64
		Successful int64
		Failed     int64
	}
	Hits struct {
		Total    int64
		MaxScore float64
		Hits     []ElasticSearchHit
	}
}

// ElasticSearchHit is the main elasticsearch single hit object.
type ElasticSearchHit struct {
	Index   string          `json:"_index"`
	Type    string          `json:"_type"`
	ID      string          `json:"_id"`
	Version int64           `json:"_version"`
	Score   float64         `json:"_score"`
	Found   bool            `json:"found"`
	Source  json.RawMessage `json:"_source"`
}

// ElasticMapping is
type ElasticMapping struct {
	Mappings map[string]interface{}
}

// ElasticBulkResponse {"took":7,"items":[{"create":{"_index":"test","_type":"type1","_id":"1","_version":1}}]}
type ElasticBulkResponse struct {
	Took  int64
	Items []struct {
		ID     string `json:"_id"`
		Index  string `json:"_index"`
		Type   string `json:"_type"`
		Status int    `json:"status"`
	} `json:"index"`
}

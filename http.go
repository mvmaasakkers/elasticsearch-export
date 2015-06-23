package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Get Executes the ES query and returns results in an ElasticSearchResponse struct
func Get(amount, from int64) (*ElasticSearchResponse, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%s/%s/_search?size=%d&from=%d", *sourceHostname, *sourcePort, *sourceIndex, amount, from)

	req, errNR := http.NewRequest("GET", url, nil)
	if errNR != nil {
		return nil, errNR
	}
	resp, errDO := client.Do(req)
	if errDO != nil {
		return nil, errDO
	}

	//
	body, errRA := ioutil.ReadAll(resp.Body)
	if errRA != nil {
		return nil, errRA
	}
	resp.Body.Close()
	elasticSearchResponse := &ElasticSearchResponse{}

	errJSON := json.Unmarshal(body, elasticSearchResponse)
	if errJSON != nil {
		return nil, errJSON
	}

	return elasticSearchResponse, nil
}

// GetMapping is
func GetMapping() (map[string]ElasticMapping, error) {
	em := make(map[string]ElasticMapping)
	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%s/%s/_mapping", *sourceHostname, *sourcePort, *sourceIndex)
	log.Println("Get mapping", *sourceIndex)
	req, errNR := http.NewRequest("GET", url, nil)
	if errNR != nil {
		return em, errNR
	}
	resp, errDO := client.Do(req)
	if errDO != nil {
		return em, errDO
	}

	//
	body, errRA := ioutil.ReadAll(resp.Body)
	if errRA != nil {
		return em, errRA
	}
	resp.Body.Close()

	json.Unmarshal(body, &em)

	return em, nil
}

// PutMapping is
func PutMapping(indexType string, mapping interface{}) error {
	data, _ := json.Marshal(mapping)

	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%s/%s/_mapping/%s", *destinationHostname, *destinationPort, *destinationIndex, indexType)
	log.Println("Put mapping", *destinationIndex, "for type", indexType)
	req, errNR := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if errNR != nil {
		return errNR
	}
	resp, errDO := client.Do(req)
	if errDO != nil {
		return errDO
	}

	//
	_, errRA := ioutil.ReadAll(resp.Body)
	if errRA != nil {
		return errRA
	}
	resp.Body.Close()

	return nil
}

// Put is
func Put(data []byte) (*ElasticBulkResponse, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%s/_bulk", *destinationHostname, *destinationPort)
	// log.Printf("Putting bulk data\n")
	req, errNR := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if errNR != nil {
		return nil, errNR
	}
	resp, errDO := client.Do(req)
	if errDO != nil {
		return nil, errDO
	}

	body, errRA := ioutil.ReadAll(resp.Body)
	if errRA != nil {
		return nil, errRA
	}
	resp.Body.Close()

	ebr := &ElasticBulkResponse{}
	json.Unmarshal(body, ebr)

	return ebr, nil
}

// ElasticIndex is
type ElasticIndex struct {
	Index struct {
		ID    string `json:"_id"`
		Index string `json:"_index"`
		Type  string `json:"_type"`
	} `json:"index"`
}

// Bulk creates bulk data update
func (e *ElasticSearchResponse) Bulk() []byte {
	var data []byte
	var endLine = []byte("\n")
	// Index
	for _, hit := range e.Hits.Hits {
		ei := ElasticIndex{}
		ei.Index.Index = *destinationIndex
		ei.Index.ID = hit.ID
		ei.Index.Type = hit.Type

		hitID, _ := json.Marshal(ei)
		hitID = append(hitID, endLine...)

		data = append(data, hitID...)
		data = append(data, []byte(hit.Source)...)
		data = append(data, endLine...)
	}

	return data
}

// CheckIndex is
func CheckIndex() error {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%s/%s/_status", *destinationHostname, *destinationPort, *destinationIndex)

	log.Printf("Checking if index %s exists\n", *destinationIndex)

	req, errNR := http.NewRequest("GET", url, nil)
	if errNR != nil {
		return errNR
	}
	resp, errDO := client.Do(req)
	if errDO != nil {
		return errDO
	}

	if resp.StatusCode == 404 {
		url := fmt.Sprintf("http://%s:%s/%s/", *destinationHostname, *destinationPort, *destinationIndex)
		log.Printf("Creating index %s\n", *destinationIndex)
		req, errNR := http.NewRequest("PUT", url, nil)
		if errNR != nil {
			return errNR
		}
		_, errDO := client.Do(req)
		if errDO != nil {
			return errDO
		}
	}

	return nil
}

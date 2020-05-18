package lema

import (
	"encoding/json"
	"net/http"
)

// Result is main lema output
type Result struct {
	Ending string `json:"ending"`
	Mi     []Mi   `json:"mi"`
	Suffix string `json:"suffix"`
	Word   string `json:"word"`
}

// Mi is lema result mi information
type Mi struct {
	MF    string `json:"mf"`
	Mi    string `json:"mi"`
	MiVdu string `json:"mi_vdu"`
	MIS   string `json:"mis"`
}

// Analyze analyze word in lema
func Analyze(w string) (*Result, error) {
	response, err := http.Get("http://localhost:7020/analyze/" + w)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	d := json.NewDecoder(response.Body)
	var l Result
	err = d.Decode(&l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

const readingsEndpointFmt = "https://engage.efergy.com/mobile_proxy/getCurrentValuesSummary?token=%s"

type EfergyClient struct {
	token string
	log   *logrus.Entry
}

type Reading struct {
	Cid       string `json:"cid"`
	LastValue int
	Data      []map[string]int `json:"data"`
	Sid       string           `json:"sid"`
	Units     string           `json:"units"`
	Age       int              `json:"age"`
}

func (reading *Reading) ParseLastValue() {
	// Find the newest reading key.
	if len(reading.Data) > 0 {
		for key := range reading.Data[0] {
			reading.LastValue = reading.Data[0][key]
			return
		}
	}
}

func NewEfergyClient(token string, log *logrus.Entry) *EfergyClient {
	return &EfergyClient{
		token: token,
		log:   log,
	}
}

func (client *EfergyClient) GetLatestReadings() ([]*Reading, error) {
	resp, err := http.Get(fmt.Sprintf(readingsEndpointFmt, client.token))
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	client.log.Println(resp.Body)

	var data []*Reading
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	// Parse the readings.
	for _, reading := range data {
		reading.ParseLastValue()
		client.log.Println(reading)
	}

	return data, nil
}

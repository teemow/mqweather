package wunderground

import (
	"fmt"
	"net/http"

	"encoding/json"
	"log"

	"github.com/teemow/mqweather/wunderground/config"
	"github.com/teemow/mqweather/wunderground/features"
)

const api = "http://api.wunderground.com"

type Wunderground struct {
	apiKey  string
	apiHost string
	client  *http.Client
}

func New(cnf *config.Config) *Wunderground {

	err := cnf.Validate()
	if err != nil {
		panic(err)
	}

	return &Wunderground{
		apiKey:  cnf.ApiKey,
		apiHost: cnf.ApiHost,
		client:  http.DefaultClient,
	}
}

func DefaultConfig(apiKey string) *config.Config {
	return &config.Config{
		ApiKey:  apiKey,
		ApiHost: api,
	}
}

func (w *Wunderground) Conditions(station string) (*features.ConditionsResponse, error) {
	url := fmt.Sprintf("%s/api/%s/conditions/q/%s.json", w.ApiHost(), w.ApiKey(), station)

	resp, err := w.makeRequest(url)
	defer resp.Body.Close()
	if err != nil {
		log.Panic("Error making request", err)
	}

	cond := &features.ConditionsResponse{}
	err = json.NewDecoder(resp.Body).Decode(cond)
	if err != nil {
		log.Panic("Error unmarshalling condition body ", err)
	}

	err = cond.Response.HasError()
	if err != nil {
		return nil, err
	}

	return cond, nil
}

func (w *Wunderground) ApiKey() string {
	return w.apiKey
}

func (w *Wunderground) ApiHost() string {
	return w.apiHost
}

func (w *Wunderground) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

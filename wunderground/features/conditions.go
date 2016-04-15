package features

type ConditionsResponse struct {
	Response  *Response   `json:"response"`
	Condition *Conditions `json:"current_observation,omitempty"`
}

type Conditions struct {
	TempC            float64 `json:"temp_c"`
	RelativeHumidity string  `json:"relative_humidity"`
	WindKPH          float64 `json:"wind_kph"`
	PressureMB       string  `json:"pressure_mb"`
	DewpointC        float64 `json:"dewpoint_c"`
}

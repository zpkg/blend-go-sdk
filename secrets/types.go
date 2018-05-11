package secrets

// Values is a bag of values.
type Values = map[string]string

// StorageResponse is the response to a vault storage request
type StorageResponse struct {
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int64       `json:"lease_duration"`
	Data          StorageData `json:"data"`
}

// StorageData is the data of a response
type StorageData struct {
	Value Values `json:"value"`
}

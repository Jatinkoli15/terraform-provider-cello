package models

type Integration struct {
	IntegrationDetails interface{} `json:"integration_details"`
	IntegrationType    string      `json:"integration_type"`
	Name               string      `json:"name"`
}

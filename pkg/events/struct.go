package events

// SecretsManagerRequest contains data coming from the Secrets Manager
type SecretsManagerRequest struct {
	Detail Detail `json:"detail"`
}

// Detail is the detail section of the request
type Detail struct {
	EventName         string            `json:"eventName"`
	RequestParameters RequestParameters `json:"requestParameters"`
}

// RequestParameters is the requestParameters section of the detail
type RequestParameters struct {
	Name     string `json:"name"`
	SecretID string `json:"secretId"`
}

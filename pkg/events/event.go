package events

// GetEventNameFromEventStruct returns the eventName from the detail section of the event.
// This function is created to encapsulate the event structure.
func GetEventNameFromEventStruct(event *SecretsManagerRequest) string {
	return event.Detail.EventName
}

// GetRequestParametersName returns the name (the name of the secret in Secrets Manager)
// from the requestParameters section inside the detail section of the event.
// This function is created to encapsulate the event structure.
func GetRequestParametersName(event *SecretsManagerRequest) string {
	res := event.Detail.RequestParameters.Name
	// This is a WTF moment which cost me at least half an hour of debugging in lambda/cloudwatch.
	// There is no consistency regarding AWS event.
	// When the action is create/delete, you get the name of the secret under "name";
	// but when the action is update, you fucking get it under "secretId".
	if res == "" {
		res = event.Detail.RequestParameters.SecretID
	}
	return res
}

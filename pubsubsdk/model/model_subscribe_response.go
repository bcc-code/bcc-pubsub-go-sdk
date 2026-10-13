package model

// SubscribeResponse is returned when creating or deleting a subscription.
// Orgs and Scopes are comma-separated lists.
type SubscribeResponse struct {
	SubscriptionId string `json:"subscriptionId"`
	Type           string `json:"type"`
	Endpoint       string `json:"endpoint"`
	ClientId       string `json:"clientId"`
	Orgs           string `json:"orgs"`
	Scopes         string `json:"scopes"`
}

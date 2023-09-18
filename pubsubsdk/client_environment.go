package pubsubsdk

type ClientCredentialsEnv struct {
	TokenUrl string
	Audience string
}

const UrlSandbox = "https://sandbox-api.bcc.no/pubsub"
const UrlProd = "https://api.bcc.no/pubsub"

var CredEnvSandbox = ClientCredentialsEnv{
	TokenUrl: "https://bcc-sso-sandbox.eu.auth0.com/oauth/token",
	Audience: "sandbox-api.bcc.no",
}

var CredEnvProd = ClientCredentialsEnv{
	TokenUrl: "https://bcc-sso.eu.auth0.com/oauth/token",
	Audience: "api.bcc.no",
}

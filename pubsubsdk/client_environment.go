package pubsubsdk

type ClientCredentialsEnv struct {
	TokenUrl string
	Audience string
}

const UrlSandbox = "https://pubsub.api.sandbox.bcc.no"
const UrlProd = "https://pubsub.api.bcc.no"

var CredEnvSandbox = ClientCredentialsEnv{
	TokenUrl: "https://bcc-sso-sandbox.eu.auth0.com/oauth/token",
	Audience: "sandbox-api.bcc.no",
}

var CredEnvProd = ClientCredentialsEnv{
	TokenUrl: "https://login.bcc.no/oauth/token",
	Audience: "api.bcc.no",
}

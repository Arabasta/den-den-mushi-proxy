package config

type JwtIssuer struct {
	Issuer        string
	Secret        string
	ExpirySeconds int
}

type JwtAudience struct {
	ExpectedIssuer   string
	ExpectedAudience string
	ExpectedTyp      string
	Secret           string
}

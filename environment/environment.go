package environment

type Environment string

const (
	LOCAL       Environment = "local"
	DEVELOPMENT Environment = "development"
	TESTING     Environment = "testing"
	PRODUCTION  Environment = "production"
)

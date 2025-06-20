package config

type Config struct {
	App struct {
		Name        string
		Environment string
		Version     string
		Port        int
		BaseURL     string
	}

	Development struct {
		SshTestKeyPath string
	}

	Logging struct {
		Level  string
		Format string
	}

	Puppet struct {
		Endpoints          []string
		TaskTimeoutSeconds int
		RetryAttempts      int
	}

	Token struct {
		Secret string
		Expiry int
	}
}

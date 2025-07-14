package config

type Config struct {
	App struct {
		Name        string
		Environment string
		Version     string
		Port        int
	}

	Development struct {
		UseInMemoryRepository bool
	}

	Logging struct {
		Level    string
		Format   string
		Output   string
		FilePath string
	}

	Token struct {
		Issuer string
		Secret string
		Ttl    int
	}
}

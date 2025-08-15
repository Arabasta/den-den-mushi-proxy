package config

type App struct {
	Name        string
	Environment string
	Version     string
	Port        int
	IsLocalHost bool
}

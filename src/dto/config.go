package dto

type Config struct {
	API APIConfig `yaml:"api"`
}

type APIConfig struct {
	Type string `yaml:"type"`
	Port string `yaml:"port"`
}

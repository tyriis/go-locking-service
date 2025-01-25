package domain

type Config struct {
	Redis struct {
		Host   string `yaml:"host"`
		Port   string `yaml:"port"`
		Prefix string `yaml:"keyPrefix"`
	} `yaml:"redis"`
	Api struct {
		Port string `yaml:"port"`
	} `yaml:"api"`
}

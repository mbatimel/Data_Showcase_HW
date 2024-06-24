package config

type Redis struct {
	Host string		`yaml:"host"`
	Port string		`yaml:"port"`
}

type Cache struct {
	Cap int			`yaml:"capacity"`
	InMemory bool	`yaml:"inMemory"`
	Redis Redis		`yaml:"redis"`
}

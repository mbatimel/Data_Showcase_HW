package config

type Config struct {
	Server Server	`yaml:"server"`
	Cache Cache		`yaml:"cache"`
}
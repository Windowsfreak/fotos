package domain

type ConfigStruct struct {
	Threads        int    `yaml:"Threads"`
	InPrefix       string `yaml:"InPrefix"`
	OutPrefix      string `yaml:"OutPrefix"`
	InFolder       string `yaml:"InFolder"`
	OutFolder      string `yaml:"OutFolder"`
	Nonce          string `yaml:"Nonce"`
	Plausibility   bool   `yaml:"Plausibility"`
	ServerAddr     string `yaml:"ServerAddr"`
	SecretKey      string `yaml:"SecretKey"`
	LastCompletion int64  `yaml:"LastCompletion"`
}

package types

type InfraConfig struct {
	Servers struct {
		Gitlab  ServerConfig `yaml:"gitlab"`
		Monitor ServerConfig `yaml:"monitor"`
		Apps    ServerConfig `yaml:"apps"`
	} `yaml:"servers"`
	Domain struct {
		Base     string `yaml:"base"`
		Gitlab   string `yaml:"gitlab"`
		Registry string `yaml:"registry"`
		Monitor  string `yaml:"monitor"`
		Apps     string `yaml:"apps"`
	} `yaml:"domain"`
	Auth struct {
		Password string `yaml:"password"`
		SSHKey   string `yaml:"ssh_key"`
	} `yaml:"auth"`
}

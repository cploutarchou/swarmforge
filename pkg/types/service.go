package types

type AppType string
type Language string

const (
	APIApp        AppType = "api"
	StandaloneApp AppType = "standalone"

	Go     Language = "go"
	Python Language = "python"
	Node   Language = "node"
	Java   Language = "java"
)

type ServiceConfig struct {
	ServiceName string
	AppType     AppType
	Language    Language
	ImageName   string
	Version     string
	Replicas    int
	CPU         string
	Memory      string
	Port        int
	Environment []string
	Domain      string
	Subdomain   string
	UseTraefik  bool
}

// Service represents a Docker service
type Service struct {
	Name        string
	Image       string
	Replicas    int
	Labels      map[string]string
	Environment map[string]string
}

// Volume represents a Docker volume
type Volume struct {
	Name       string
	Driver     string
	Mountpoint string
}

// Config represents a Docker config
type Config struct {
	Name     string
	Data     []byte
	Encoding string
}

// Secret represents a Docker secret
type Secret struct {
	Name     string
	Data     []byte
	Encoding string
}

// DeploymentConfig represents the configuration for a service deployment
type DeploymentConfig struct {
	ServiceName string
	ImageName   string
	Port        int
	Replicas    int
	Domain      string
	Subdomain   string
	UseTraefik  bool
	Labels      map[string]string
	Environment map[string]string
	Email       string
}

func ValidAppTypes() []AppType {
	return []AppType{APIApp, StandaloneApp}
}

func ValidLanguages() []Language {
	return []Language{Go, Python, Node, Java}
}

func (a AppType) String() string {
	return string(a)
}

func (l Language) String() string {
	return string(l)
}

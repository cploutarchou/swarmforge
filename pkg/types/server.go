package types

type ServerRole string

const (
	ManagerServer ServerRole = "manager"
	GitlabServer ServerRole = "gitlab"
	MonitorServer ServerRole = "monitor"
	AppsServer    ServerRole = "apps"
)

type ServerConfig struct {
	IP       string
	Username string
	Password string
	Role     ServerRole
	Labels   map[string]string
}

func ValidServerRoles() []ServerRole {
	return []ServerRole{
		ManagerServer,
		GitlabServer,
		MonitorServer,
		AppsServer,
	}
}

func (r ServerRole) String() string {
	return string(r)
}

func IsValidServerRole(role string) bool {
	for _, r := range ValidServerRoles() {
		if r.String() == role {
			return true
		}
	}
	return false
}

// GetServerLabels returns the appropriate Docker labels for each server role
func GetServerLabels(role ServerRole) map[string]string {
	labels := map[string]string{
		"role": role.String(),
	}

	switch role {
	case GitlabServer:
		labels["service.type"] = "gitlab"
		labels["backup.enabled"] = "true"
	case MonitorServer:
		labels["service.type"] = "monitor"
		labels["metrics.enabled"] = "true"
	case AppsServer:
		labels["service.type"] = "application"
		labels["app.deployment"] = "enabled"
	case ManagerServer:
		labels["node.role"] = "manager"
		labels["traefik.enabled"] = "true"
	}

	return labels
}

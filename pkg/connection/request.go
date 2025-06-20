package connection

type ConnectionPurpose string

const (
	Change      ConnectionPurpose = "change"
	Healthcheck ConnectionPurpose = "healthcheck"
)

type ConnectionType string

const (
	/* For development purposes only */
	LocalShell ConnectionType = "local_shell"
	SshTestKey ConnectionType = "ssh_test_key"
	/* For development purposes only End*/

	SshOrchestratorKey ConnectionType = "ssh_orchestrator_key"
	SshPassword        ConnectionType = "ssh_password"
)

type Connection struct {
	Purpose          ConnectionPurpose `json:"purpose"`
	Type             ConnectionType    `json:"type"`
	ServerIP         string            `json:"server_ip"`
	Port             string            `json:"port"`
	OSUser           string            `json:"os_user"`
	ChangeID         string            `json:"change_id,omitempty"`
	HealthcheckGroup string            `json:"healthcheck_group,omitempty"`
}

package types

type ConnectionMethod string

const (
	/* For development purposes only */
	LocalShell ConnectionMethod = "local_shell"
	SshTestKey ConnectionMethod = "ssh_test_key"
	/* For development purposes only End*/

	/* New Connection */
	SshOrchestratorKey ConnectionMethod = "ssh_orchestrator_key"
	SshPassword        ConnectionMethod = "ssh_password"
)

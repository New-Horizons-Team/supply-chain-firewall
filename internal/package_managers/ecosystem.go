package package_managers

type Ecosystem string

const (
	EcosystemNpm  Ecosystem = "npm"
	EcosystemPypi Ecosystem = "PyPI"
	EcosystemGo   Ecosystem = "Go"
)

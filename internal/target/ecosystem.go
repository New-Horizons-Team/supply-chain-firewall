package target

type Ecosystem string

const (
	EcosystemNpm  Ecosystem = "npm"
	EcosystemPypi Ecosystem = "PyPI"
	EcosystemGo   Ecosystem = "Go"
)

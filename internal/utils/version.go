package utils

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	buildby = "unknown"
)

func GenerateVersionSignature() string {
	return fmt.Sprintf("azure-dyndns2 %s, commit %s, built at %s, build by %s", version, commit, date, buildby)
}

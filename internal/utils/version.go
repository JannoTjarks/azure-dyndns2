package utils

import (
	"encoding/json"
	"fmt"
)

type versionStruct struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Buildby string `json:"buildby"`
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	buildby = "unknown"
)

func GenerateVersionSignature() string {
	return fmt.Sprintf("azure-dyndns2 %s, commit %s, built at %s, build by %s", version, commit, date, buildby)
}

func GenerateVersionJson() string {
	versionStruct := versionStruct{
		Version: version,
		Commit: commit,
		Date: date,
		Buildby: buildby,
	}

    jsonBytes, _ := json.Marshal(versionStruct)
    return string(jsonBytes)
}

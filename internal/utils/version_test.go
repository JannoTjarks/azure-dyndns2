package utils_test

import (
	"github.com/JannoTjarks/azure-dyndns2/internal/utils"
	"testing"
)

func TestGenerateVersionSignature(t *testing.T) {
	want := "azure-dyndns2 dev, commit none, built at unknown, build by unknown"

	msg := utils.GenerateVersionSignature()

	if msg != want {
		t.Errorf(`utils.GenerateVersionSignature() = %q, want match for %#q, nil`, msg, want)
	}
}

func TestGenerateVersionJson(t *testing.T) {
	want := `{"version":"dev","commit":"none","date":"unknown","buildby":"unknown"}`

	msg := utils.GenerateVersionJson()

	if msg != want {
		t.Errorf(`utils.GenerateVersionSignature() = %q, want match for %#q, nil`, msg, want)
	}
}

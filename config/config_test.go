package config

import (
	"os"
	"testing"
)

func TestGetEnvVariable(t *testing.T) {
	want := "testvalue"
	_ = os.Setenv("TESTENVVAR", want)

	defer func() {
		_ = os.Unsetenv("TESTENVVAR")
	}()

	got := getEnvVariable("TESTENVVAR", "defaultvalue")
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}

	want = "defaultvalue"
	got = getEnvVariable("NONEXISTENTENVVAR", want)
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
}

func TestGetEnvVariableAsInt(t *testing.T) {
	want := 42
	_ = os.Setenv("TESTINTENVVAR", "42")

	defer func() {
		_ = os.Unsetenv("TESTINTENVVAR")
	}()

	got := getEnvVariableAsInt("TESTINTENVVAR", 0)
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}

	want = 100
	got = getEnvVariableAsInt("NONEXISTENTENVVAR", want)
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}

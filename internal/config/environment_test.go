package config

import "testing"

func EnvToIntValid(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

func EnvToIntNotInt(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

func GetFromMapValid(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

func GetFromMapNoKey(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

func ReadEnvValid(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

func ReadEnvAllAreEmpty(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

func RadEnvSomeAreEmpty(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

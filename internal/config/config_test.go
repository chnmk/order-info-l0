package config

import (
	"testing"
)

func NewConfigValid(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

func NewConfigTwice(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

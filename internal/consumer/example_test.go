package consumer

import "testing"

func GoFakeValid(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

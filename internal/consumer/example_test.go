package consumer

import "testing"

func goFakeValid(t *testing.T) {
	var err error

	if err != nil {
		t.Error(err)
	}
}

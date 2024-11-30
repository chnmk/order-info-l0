package config

import (
	"os"
	"testing"

	"github.com/chnmk/order-info-l0/internal/models"
)

func TestEmptyConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on empty map access")
		}
	}()

	os.Setenv("SERVER_PORT", "3000")

	var cfg1 models.Config

	example := cfg1.Get("SERVER_PORT")
	_ = example
}

func TestNewConfig(t *testing.T) {
	cfg1 := NewConfig()
	if cfg1 == nil {
		t.Fatalf("created config shouldn't be nil")
	}

	cfg2 := NewConfig()
	if cfg2 != cfg1 {
		t.Fatalf("config should only be created once")
	}
}
func TestReadEnv(t *testing.T) {
	os.Setenv("SERVER_PORT", "3000")

	cfg1 := NewConfig()

	example := cfg1.Get("SERVER_PORT")
	if example != "3000" {
		t.Fatalf("ReadEnv() didn't read default env variable")
	}

	//

	os.Setenv("SERVER_PORT", "101010101010101")

	cfg1.ReadEnv()

	example_new := cfg1.Get("SERVER_PORT")
	if example_new == example || example_new != "101010101010101" {
		t.Fatalf("ReadEnv() didn't read new env variable")
	}
}

func TestEnvToInt(t *testing.T) {
	os.Setenv("SERVER_PORT", "3000")

	cfg1 := NewConfig()
	cfg1.ReadEnv()

	port_int := envToInt("SERVER_PORT")
	if port_int != 3000 {
		t.Fatalf("unexpedted envToInt result: expected 3000, found %d", port_int)
	}

	envToInt("KAFKA_PROTOCOL")

	select {
	case <-ExitCtx.Done():
	default:
		t.Fatalf("expected service shutdown on error")
	}

	Exit()
}

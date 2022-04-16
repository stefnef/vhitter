package config

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func loadDefaultConfig() {
	NewConfigFile("./testConfig.yml")
}

func Test_should_throw_error_if_file_does_not_exist(t *testing.T) {
	if gotConfig, err := NewConfigFile("./i-do-not-exist.yml"); gotConfig != nil {
		t.Fatalf("could load config from non-existing file")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected error: got %s", err)
	}
}

func Test_Should_read_config_file(t *testing.T) {
	expConfig := &Config{Twitter: &TwitterConfig{UserId: "1234", Bearer: "<ACCESS_TOKEN>"}}

	gotConfig, err := NewConfigFile("./testConfig.yml")
	if err != nil {
		t.Errorf("Failed to load config. Reason: %s", err)
	}
	if gotConfig == nil {
		t.Error("Failed to load config")
	} else if !reflect.DeepEqual(expConfig, gotConfig) {
		t.Fatalf(`want: %v, got: %v`, expConfig, gotConfig)
	}
}

func Test_should_return_twitter_configuration(t *testing.T) {
	expConfig := &Config{Twitter: &TwitterConfig{UserId: "1234", Bearer: "<ACCESS_TOKEN>"}}
	loadDefaultConfig()

	gotTwitter := GetTwitterConfig()
	if !reflect.DeepEqual(*expConfig.Twitter, gotTwitter) {
		t.Errorf("wrong twitter configuration received. Want: %v. Got: %v", expConfig.Twitter, gotTwitter)
	}
}

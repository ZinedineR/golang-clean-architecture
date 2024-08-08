package config_test

import (
	"boiler-plate-clean/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitAppConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("APP_ENV", "development")
	os.Setenv("APP_DEBUG", "True")
	os.Setenv("APP_VERSION", "1.0")
	os.Setenv("APP_NAME", "TestApp")

	// Call the function to initialize config
	config := config.InitAppConfig()

	//Verify the values
	if config.AppEnv != "development" {
		t.Errorf("Expected AppEnv to be 'development', but got %s", config.AppEnv)
	}

	if config.AppDebug != "True" {
		t.Errorf("Expected AppDebug to be 'True', but got %s", config.AppDebug)
	}

	if config.AppVersion != "1.0" {
		t.Errorf("Expected AppVersion to be '1.0', but got %s", config.AppVersion)
	}

	if config.AppName != "TestApp" {
		t.Errorf("Expected AppName to be 'TestApp', but got %s", config.AppName)
	}
	assert.NotNil(t, config.AppEnv)
	assert.NotNil(t, config.AppDebug)
	assert.NotNil(t, config.AppVersion)
	assert.NotNil(t, config.AppName)

	// You can write more assertions for other fields if needed
}

func TestIsStaging(t *testing.T) {
	c := config.Config{AppEnv: "development"}
	if !c.IsStaging() {
		t.Error("Expected IsStaging to return true, but it returned false")
	}

	c.AppEnv = "production"
	if c.IsStaging() {
		t.Error("Expected IsStaging to return false, but it returned true")
	}
}

func TestIsProd(t *testing.T) {
	c := config.Config{AppEnv: "production"}
	if !c.IsProd() {
		t.Error("Expected IsProd to return true, but it returned false")
	}

	c.AppEnv = "development"
	if c.IsProd() {
		t.Error("Expected IsProd to return false, but it returned true")
	}
}

func TestIsDebug(t *testing.T) {
	c := config.Config{AppDebug: "True"}
	if !c.IsDebug() {
		t.Error("Expected IsDebug to return true, but it returned false")
	}

	c.AppDebug = "False"
	if c.IsDebug() {
		t.Error("Expected IsDebug to return false, but it returned true")
	}
}

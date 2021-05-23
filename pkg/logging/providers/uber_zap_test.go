package providers

import (
	"github.com/ereb-or-od/kenobi/pkg/logging/options"
	"testing"
)

func TestNewUberZapLoggerShouldReturnDefaultUberZapLoggerWhenConfigurationDoesNotSelected(t *testing.T) {
	uberZapLogger, err := NewUberZapLogger()
	if err != nil {
		t.Error("error does not expected")
	}

	if uberZapLogger == nil {
		t.Error("default-logger must be initialized")
	}
}

func TestNewUberZapLoggerWithOptionsShouldReturnUberZapLoggerWhenConfigurationSelected(t *testing.T) {
	uberZapLogger, err := NewUberZapLoggerWithOptions(options.NewDefaultLoggerOptions())
	if err != nil {
		t.Error("error does not expected")
	}

	if uberZapLogger == nil {
		t.Error("default-logger must be initialized")
	}
}

func TestNewUberZapLoggerWithOptionsShouldReturnErrorWhenConfigurationIsNil(t *testing.T) {
	uberZapLogger, err := newUberZapLogger(nil)
	if err == nil {
		t.Errorf("error should be expected as %s", configurationMustBeSpecifiedError)
	}

	if uberZapLogger != nil {
		t.Error("default-logger should not be initialized")
	}
}

func TestNewZapConfigShouldReturnDefaultZapConfigWhenConfigurationSelectedAsDefault(t *testing.T) {
	config := newZapConfig(options.NewDefaultLoggerOptions())
	if config == nil {
		t.Error("config should be initialized")
	}
}

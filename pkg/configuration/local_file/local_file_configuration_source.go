package local_file

import (
	"github.com/ereb-or-od/kenobi/pkg/configuration/interfaces"
	"github.com/spf13/viper"
	"time"
)

type localConfigurationSource struct {
	config *viper.Viper
}

func (l localConfigurationSource) GetValueByKey(key string) interface{} {
	return l.config.Get(key)
}

func (l localConfigurationSource) GetIntArrayValueByKey(key string) []int {
	return l.config.GetIntSlice(key)
}

func (l localConfigurationSource) GetDurationValueByKey(key string) time.Duration {
	return l.config.GetDuration(key)
}

func (l localConfigurationSource) GetStringArrayValueByKey(key string) []string {
	return l.config.GetStringSlice(key)
}

func (l localConfigurationSource) GetStringValueByKey(key string) string {
	return l.config.GetString(key)
}

func (l localConfigurationSource) GetIntValueByKey(key string) int {
	return l.config.GetInt(key)
}

func (l localConfigurationSource) GetInt64ValueByKey(key string) int64 {
	return l.config.GetInt64(key)
}

func (l localConfigurationSource) GetFloatValueByKey(key string) float64 {
	return l.config.GetFloat64(key)
}

func (l localConfigurationSource) GetBooleanValueByKey(key string) bool {
	return l.config.GetBool(key)
}

func (l localConfigurationSource) GetTimeValueByKey(key string) time.Time {
	return l.config.GetTime(key)
}

func NewWithOptions(fileName string, fileType string, filePath string) (interfaces.ConfigurationSource, error) {
	if options, err := New(fileName, fileType, filePath); err != nil {
		return nil, err
	} else {
		config := viper.New()
		config.SetConfigName(options.GetFileName())
		config.SetConfigType(options.GetFileType())
		config.AddConfigPath(options.GetFilePath())
		err := config.ReadInConfig()
		if err != nil {
			return nil, err
		}
		return &localConfigurationSource{
			config: config,
		}, nil
	}
}

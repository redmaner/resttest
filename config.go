package resttest

import "github.com/spf13/viper"

type Config struct {
	BaseUrl string       `mapstructure:"base_url"`
	Headers []HttpHeader `mapstructure:"http_headers"`
	Tests   []HttpTest   `mapstructure:"tests"`
}

type HttpHeader struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

type HttpTest struct {
	Method string          `mapstructure:"method"`
	Body   string          `mapstructure:"body"`
	Path   string          `mapstructure:"path"`
	Expect TestExpectation `mapstructure:"expect"`
}

type TestExpectation struct {
	StatusCode int          `mapstructure:"status_code"`
	Headers    []HttpHeader `mapstructure:"http_headers"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := new(Config)
	return cfg, viper.Unmarshal(cfg)
}

package models

type ServerConfig struct {
	Port           string         `mapstructure:"port" yaml:"port"`
	Host           string         `mapstructure:"host" yaml:"host"`
	LLMConfig      LLMConfig      `mapstructure:"llm" yaml:"llm"`
	DocumentConfig DocumentConfig `mapstructure:"document" yaml:"document"`
}

type LLMConfig struct {
	BASE_URL string `mapstructure:"base_url" yaml:"base_url"`
	MODEL    string `mapstructure:"model" yaml:"model"`
	API_KEY  string `mapstructure:"api_key" yaml:"api_key"`
}

type DocumentConfig struct {
	Addr    string `mapstructure:"addr" yaml:"addr"`
	API_KEY string `mapstructure:"api_key" yaml:"api_key"`
	Model   string `mapstructure:"model" yaml:"model"`
}

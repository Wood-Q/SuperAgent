package config

type LLMConfig struct {
	BASE_URL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"`
	MODEL    string `mapstructure:"model" json:"model" yaml:"model"`
}

type ServerConfig struct {
	HOST      string    `mapstructure:"host" json:"host" yaml:"host"`
	PORT      string    `mapstructure:"port" json:"port" yaml:"port"`
	LLMConfig LLMConfig `mapstructure:"llm" json:"llm" yaml:"llm"`
}

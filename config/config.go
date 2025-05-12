package config

type LLMConfig struct {
	BASE_URL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"`
	MODEL    string `mapstructure:"model" json:"model" yaml:"model"`
}

type ServerConfig struct {
	HOST           string         `mapstructure:"host" json:"host" yaml:"host"`
	PORT           string         `mapstructure:"port" json:"port" yaml:"port"`
	LLMConfig      LLMConfig      `mapstructure:"llm" json:"llm" yaml:"llm"`
	DocumentConfig DocumentConfig `mapstructure:"document" json:"document" yaml:"document"`
}

type DocumentConfig struct {
	Addr      string `mapstructure:"addr" json:"addr" yaml:"addr"`
	ArkApiKey string `mapstructure:"ark_api_key" json:"ark_api_key" yaml:"ark_api_key"`
	ArkModel  string `mapstructure:"ark_model" json:"ark_model" yaml:"ark_model"`
}

package config

type tomlConfig struct {
	Http   http   `toml:"http"`
	Log    log    `toml:"log"`
	Verify verify `toml:"verify"`
}

type http struct {
	Port int16
}

type log struct {
	Level string
	Path  string
}

type verify struct {
	Url string
}

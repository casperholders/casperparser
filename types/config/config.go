package config

var ConfigParsed Config

type ContractType struct {
	Entrypoints []Entrypoint
	NamedKeys   []string
}

type Entrypoint struct {
	Name string
	Args []string
}

type ModuleByte struct {
	StrictArgs bool
	Args       []string
	Events     []string `mapstructure:",omitempty"`
}

type Config struct {
	ContractTypes map[string]ContractType
	ModuleBytes   map[string]ModuleByte
}

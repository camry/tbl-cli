package model

type IgnoreConfig struct {
    Ignores []*IgnoreTable `yaml:"ignores"`
}

type IgnoreTable struct {
    Table   string   `yaml:"table"`
    Columns []string `yaml:"columns"`
}

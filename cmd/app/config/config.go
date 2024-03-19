package config

type Config struct {
	Default  *DefaultConfig `yaml:"default"`
	DBConfig *DBConfig      `yaml:"dbConfig"`
	SQLs     []*SQLConfig   `yaml:"sqlConfig"`
}

type DefaultConfig struct {
	Mode string    `yaml:"mode"`
	App  AppConfig `yaml:"app"`
}

type AppConfig struct {
	RpcPort  int32 `yaml:"rpcPort"`
	HttpPort int32 `yaml:"httpPort"`
}

type DBConfig struct {
	DSN           string `yaml:"dsn"`
	MaxOpenConn   int    `yaml:"maxOpenConn"`
	MinIdleConn   int    `yaml:"maxIdleConn"`
	MaxLifeSecond int    `yaml:"maxLifeTime"`
}

type SQLConfig struct {
	Name  string `yaml:"name"`
	Table string `yaml:"table"`
	Sql   string `yaml:"sql"`
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) FindSQLByName(name string) *SQLConfig {
	for _, sql := range c.SQLs {
		if sql.Name == name {
			return sql
		}
	}
	return nil
}

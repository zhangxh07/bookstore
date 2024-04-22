package conf

var (
	config *Config
)

// 把全局对象保护起来
func C() *Config {
	if config == nil {
		panic("请加载程序配置, LoadConfigFromToml/LoadConfigFromEnv")
	}
	return config
}

func LoadConfigFromToml() (*Config, error) {
	conf := DefaultConfig()
	//_, err := toml.DecodeFile(path, conf)
	//if err != nil {
	//	return nil, err
	//}

	// 赋值给全局变量
	config = conf
	return conf, nil
}

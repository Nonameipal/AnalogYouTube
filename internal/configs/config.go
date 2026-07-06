package configs

func GetConfig() (*Configs, error) {
	if err := Load(); err != nil {
		return nil, err
	}
	return &AppSettings, nil
}

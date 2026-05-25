package configs

func GetConfig() (*Configs, error) {
	if err := ReadSettings(); err != nil {
		return nil, err
	}
	return &AppSettings, nil
}

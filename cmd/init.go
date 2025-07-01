package cmd

import "github.com/jackchuka/gh-oss-watch/services"

func HandleInit(configService services.ConfigService, output services.Output) error {
	config, err := configService.Load()
	if err != nil {
		return err
	}

	configPath, err := configService.GetConfigPath()
	if err != nil {
		return err
	}

	err = configService.Save(config)
	if err != nil {
		return err
	}

	output.Printf("Initialized config file at %s\n", configPath)
	output.Println("Use 'gh oss-watch add <repo>' to start watching repositories")
	return nil
}

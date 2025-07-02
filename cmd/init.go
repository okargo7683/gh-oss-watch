package cmd

func (c *CLI) handleInit() error {
	config, err := c.configService.Load()
	if err != nil {
		return err
	}

	configPath, err := c.configService.GetConfigPath()
	if err != nil {
		return err
	}

	err = c.configService.Save(config)
	if err != nil {
		return err
	}

	c.output.Printf("Initialized config file at %s\n", configPath)
	c.output.Println("Use 'gh oss-watch add <repo>' to start watching repositories")
	return nil
}

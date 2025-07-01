package services

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type CacheServiceImpl struct{}

func NewCacheService() CacheService {
	return &CacheServiceImpl{}
}

func (c *CacheServiceImpl) Load() (*CacheData, error) {
	cachePath, err := c.getCachePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return &CacheData{
			LastCheck: time.Time{},
			Repos:     make(map[string]RepoState),
		}, nil
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache CacheData
	err = yaml.Unmarshal(data, &cache)
	if err != nil {
		return nil, err
	}

	if cache.Repos == nil {
		cache.Repos = make(map[string]RepoState)
	}

	return &cache, nil
}

func (c *CacheServiceImpl) Save(cache *CacheData) error {
	configDir, err := c.getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	cachePath, err := c.getCachePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func (c *CacheServiceImpl) getCachePath() (string, error) {
	configDir, err := c.getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "cache.yaml"), nil
}

func (c *CacheServiceImpl) getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".gh-oss-watch"), nil
}

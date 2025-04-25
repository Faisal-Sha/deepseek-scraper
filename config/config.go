package config

import (
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    Kafka struct {
        Brokers []string `yaml:"brokers"`
        Topic   string   `yaml:"topic"`
        GroupID string   `yaml:"group_id"`
    } `yaml:"kafka"`
    Scraper struct {
        BaseURL        string `yaml:"base_url"`
        MaxDepth       int    `yaml:"max_depth"`
        DelaySeconds   int    `yaml:"delay_seconds"`
        UserAgent      string `yaml:"user_agent"`
        OutputFormat   string `yaml:"output_format"` // "db" or "json"
        JSONOutputPath string `yaml:"json_output_path"`
    } `yaml:"scraper"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    config := &Config{}
    if err := yaml.Unmarshal(data, config); err != nil {
        return nil, err
    }

    return config, nil
}
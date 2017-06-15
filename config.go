package main

import (
  "fmt"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

type Config struct {
  Host string
  Port int
}

const configName string = "config.yml"

func (config Config) hasRequiredFields() bool {
  return (config.Host != "") && (config.Port != 0)
}

func readConfig() (Config, error) {
  config := Config{}

  fmt.Println("Config: Reading file...")
  file, err := ioutil.ReadFile(configName)
  if err != nil {
    return config, err
  }

  fmt.Println("Config: File read successful, parsing yaml...")
  err = yaml.Unmarshal(file, &config)
  if err != nil {
    return config, err
  }

  if !config.hasRequiredFields() {
    // TODO: The error should say what fields are missing
    return config, fmt.Errorf("Missing required fields in config")
  }

  fmt.Println("Config: All good!")

  return config, nil
}

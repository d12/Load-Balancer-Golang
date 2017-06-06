package main

import (
  "fmt"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

type Config struct {
  HostOrigin string
}

const configName string = "config.yml"

func (config Config) hasRequiredFields() bool {
  return (config.HostOrigin != "")
}

func readConfig() (Config, error) {
  config := Config{}

  file, err := ioutil.ReadFile(configName)
  if err != nil {
    return config, err
  }

  err = yaml.Unmarshal(file, &config)
  if err != nil {
    return config, err
  }

  if !config.hasRequiredFields() {
    // TODO: The error should say what fields are missing
    return config, fmt.Errorf("Missing required fields in config")
  }

  return config, nil
}
// TODO: Pass in host and port, not full name

package main

import (
  "fmt"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

const configName string = "config.yml"

func (proxy Proxy) hasRequiredFields() bool {
  return (proxy.Host != "") && (proxy.Port != 0)
}

func readConfig() (Proxy, error) {
  proxy := Proxy{}

  file, err := ioutil.ReadFile(configName)
  if err != nil {
    return proxy, err
  }

  err = yaml.Unmarshal(file, &proxy)
  if err != nil {
    return proxy, err
  }

  if !proxy.hasRequiredFields() {
    // TODO: The error should say what fields are missing
    return proxy, fmt.Errorf("Missing required fields in config")
  }

  return proxy, nil
}

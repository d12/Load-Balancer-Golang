package main

import (
  "fmt"
  "strings"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

const configName string = "config.yml"

func (proxy Proxy) validateHost() (bool, string) {
  if proxy.Host == "" {
    return true, "the 'host' field cannot be blank!"
  } else {
    return false, ""
  }
}

func (proxy Proxy) validatePort() (bool, string) {
  if proxy.Port == 0 {
    return true, "the 'port' field cannot be blank!"
  } else {
    return false, ""
  }
}

func (proxy Proxy) validateServers() (bool, string) {
  if len(proxy.Servers) == 0 {
    return true, "the config must specify at least 1 server"
  } else {
    return false, ""
  }
}

func (proxy Proxy) validateFields() error {
  var errors = []string{}
  var validations = [](func() (bool, string)){
    proxy.validateHost,
    proxy.validatePort,
    proxy.validateServers,
  }

  for _, validation := range validations {
    has_error, error_message := validation()
    if has_error {
      errors = append(errors, error_message)
    }
  }

  if(len(errors) == 0) {
    return nil
  } else {
    return fmt.Errorf(strings.Join(errors, ", "))
  }
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

  err = proxy.validateFields()
  if err != nil {
    return proxy, err
  }

  return proxy, nil
}

package main

import (
  "fmt"
  "strings"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

const configName string = "config.yml"

func validation(condition bool, errorMessage string) string {
  if condition {
    return errorMessage
  } else {
    return ""
  }
}

func removeEmpty(errors []string) []string {
  var filtered = []string{}
  for _,e := range errors {
    if e != "" {
      filtered = append(filtered, e)
    }
  }

  return filtered
}

func generateValidationErrors(proxy Proxy) []string {
  return removeEmpty([]string{
    validation(
      proxy.Host == "",
      "the 'host' field cannot be blank",
    ),
    validation(
      proxy.Port == 0,
      "the 'port' field cannot be blank",
    ),
    validation(
      len(proxy.Servers) == 0,
      "the config must specify at least 1 server",
    ),
  })
}

func validateFields(proxy Proxy) error {
  var errors = generateValidationErrors(proxy)

  if(len(errors) == 0) {
    return nil
  } else {
    return fmt.Errorf(strings.Join(errors, ", "))
  }
}

func ReadConfig() (Proxy, error) {
  proxy := Proxy{}

  file, err := ioutil.ReadFile(configName)
  if err != nil {
    return proxy, err
  }

  err = yaml.Unmarshal(file, &proxy)
  if err != nil {
    return proxy, err
  }

  err = validateFields(proxy)
  if err != nil {
    return proxy, err
  }

  return proxy, nil
}

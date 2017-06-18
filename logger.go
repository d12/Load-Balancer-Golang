package main

import (
  "fmt"
  "time"
  "os"
)

// HH/mm/SS DD/MM/YYYY
const formatString string = "15:04:05 01/02/2006"

func LogInfo(msg string) {
  printTime()
  fmt.Println(msg)
}

func LogWarn(msg string) {
  printTime()
  fmt.Print("[Warning] ")
  fmt.Println(msg)
}

func LogErr(msg string) {
  printTime()
  fmt.Print("[ERROR] ")
  fmt.Println(msg)
}

func LogErrAndCrash(msg string) {
  LogErr(msg)
  os.Exit(1)
}

func printTime() {
  fmt.Print("[" + getFormattedTime() + "] ")
}

func getFormattedTime() string {
  return time.Now().Format(formatString)
}

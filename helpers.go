package main

import (
  "os"
)

func LookupStrEnv(name string, dflt string) (string) {
  strEnvStr, ok := os.LookupEnv(name)
  if ok == false {
    return dflt
  }

  return strEnvStr
}

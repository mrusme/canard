package main

import (
  "os"
  "time"
  "bytes"
  "io"
  "strings"
  "log"
  "net/http"
  "mime/multipart"
  "encoding/json"
)

func LookupStrEnv(name string, dflt string) (string) {
  strEnvStr, ok := os.LookupEnv(name)
  if ok == false {
    return dflt
  }

  return strEnvStr
}

func call(apiKey string, urlPath string) (ApiResponse, error) {
  client := &http.Client{
    Timeout: time.Second * 10,
  }
  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  fw, err := writer.CreateFormField("api_key")
  if err != nil {
  }
  _, err = io.Copy(fw, strings.NewReader(apiKey))
  if err != nil {
      return ApiResponse{}, err
  }
  writer.Close()
  req, err := http.NewRequest("POST", urlPath, bytes.NewReader(body.Bytes()))
  if err != nil {
      return ApiResponse{}, err
  }
  req.Header.Set("Content-Type", writer.FormDataContentType())
  resp, _ := client.Do(req)
  if resp.StatusCode != http.StatusOK {
    log.Printf("Request failed with response code: %d", resp.StatusCode)
  }

  defer resp.Body.Close()

  var response ApiResponse
  if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
    log.Fatal(err)
  }

  return response, nil
}

package gosberbank

import (
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"
)

func sendPost(values map[string][]string, uri string) (body []byte, err error) {
  var req *http.Request
  var resp *http.Response
  if err != nil {
    return body, err
  }
  client := &http.Client{}
  data := url.Values(values)
  req, err = http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  if err != nil {
    return body, err
  }
  resp, err = client.Do(req)
  if err != nil {
    return body, err
  }
  body, err = ioutil.ReadAll(resp.Body)
  return body, err
}

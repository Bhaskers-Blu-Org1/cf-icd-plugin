/**
 * Copyright IBM Corporation 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 package webhook

import (
    "bytes"
    "net/http"
    "syscall"
    "os"
    "errors"
    "io/ioutil"
)

func Request(url string, method string, buf *bytes.Buffer) (string) {
    client := &http.Client{}
    req, err := http.NewRequest(method, url, buf)
    check(err)
    req.Header.Add("x-create-connection", "true")
    req.Header.Add("content-type", "application/json")
    resp, err := client.Do(req)
    check(err)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    check(err)
    return string(body)
}

/**
* Register new webhook
*/
func Register(webhook_url string) error {
    var err error
    if webhook_url[:5] != "https" {
        err = errors.New("Error: https required")
        return err
    }
    var file = ConfigFile()
    err = (*file).Truncate(0)
    _, err = (*file).WriteString(webhook_url)
    err = (*file).Close()
    return err
}

/**
* TODO: This needs to secure the webhook before storing or move
* the function to a vault
*/
func ConfigFile() (*os.File) {
    var webhookConfigFile = os.TempDir() + "webhook"
    var file *os.File
    var mode = os.FileMode(int(0600))
    if _, err := os.Stat(webhookConfigFile); os.IsNotExist(err) {
       file, err = os.Create(webhookConfigFile)
       check(err)
       err = (*file).Chmod(mode)
       check(err)
    } else {
       file, err = os.OpenFile(webhookConfigFile, syscall.O_RDWR, mode)
       check(err)
    }
    return file
}

func Config() (string, error) {
    var webhookConfigFile = os.TempDir() + "webhook"
    dat, err := ioutil.ReadFile(webhookConfigFile)
    return string(dat), err
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

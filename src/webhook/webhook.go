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

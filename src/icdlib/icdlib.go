package icdlib

import (
    "bytes"
    "net/http"
    "syscall"
    "os"
    "io/ioutil"
)

func Request(url string, method string, buf *bytes.Buffer) (string) {
    client := &http.Client{}
    req, err := http.NewRequest(method, url, buf)
    Check(err)
    resp, err := client.Do(req)
    Check(err)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    Check(err)
    return string(body)
}

func Check(e error) {
    if e != nil {
        panic(e)
    }
}

func WebhookConfigFile() (*os.File) {
    var webhookConfigFile = os.TempDir() + "webhook"
    var file *os.File
    var mode = os.FileMode(int(0600))
    if _, err := os.Stat(webhookConfigFile); os.IsNotExist(err) {
       file, err = os.Create(webhookConfigFile)
       Check(err)
       err = (*file).Chmod(mode)
       Check(err)
    } else {
       file, err = os.OpenFile(webhookConfigFile, syscall.O_RDWR, mode)
       Check(err)
    }
    return file
}

func WebhookConfig() (string) {
    var webhookConfigFile = os.TempDir() + "webhook"
    dat, err := ioutil.ReadFile(webhookConfigFile)
    Check(err)
    return string(dat)
}

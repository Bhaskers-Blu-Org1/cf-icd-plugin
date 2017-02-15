package main

import (
    "testing"
    "webhook"
)

var TEST_PORT = "6490"

func TestGitInfo(t *testing.T) {
    t.Run("Test gitinfo", func(t *testing.T) {
        //test git config
        gitinf, err := GitInfo()
        if err != nil {
           t.Errorf("Expected successful git info: %s", err)
        }
        if gitinf.GitURL != "https://github.com/" {
           t.Errorf("Unexpected git info URL: %s", gitinf.GitURL)
        }
        if gitinf.GitCommitID != "123" {
           t.Errorf("Unexpected git info commit ID: %s", gitinf.GitCommitID)
        }
        if gitinf.GitBranch != "master" {
           t.Errorf("Unexpected git info branch: %s", gitinf.GitBranch)
        }
    })
}

func TestRegisterWebhook(t *testing.T) {
    t.Run("Test --register-webhook", func(t *testing.T) {
        //make sure we catch the required https error
        err := webhook.Register("http://shouldfail")
        if err == nil || err.Error() != "Error: https required" {
           t.Errorf("Expected https required error")
        }

        //test for catching other protocols not supported
        err = webhook.Register("ftp://shouldfail")
        if err == nil || err.Error() != "Error: https required" {
           t.Errorf("Expected https required error")
        }

        //test for valid registration
        var registered_webhook = "https://localhost:" + TEST_PORT
        err = webhook.Register(registered_webhook)
        if err != nil {
           t.Errorf("Expected successful registration: %s", err)
        }

        //test previous webhook registered
        webhook, err := webhook.Config()
        if err != nil {
           t.Errorf("Expected successful query: %s", err)
        }
        if webhook != registered_webhook {
           t.Errorf("Expected webhook: %s. Actual: %s", registered_webhook, webhook)
        }
    })
}

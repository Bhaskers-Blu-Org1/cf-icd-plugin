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

 package webhook_test

import (
    "testing"
    "webhook"
)

var TEST_PORT = "6490"

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

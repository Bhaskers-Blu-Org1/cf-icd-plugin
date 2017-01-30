package main_test

import (
    "testing"
    "fmt"
)

func TestFoo(t *testing.T) {
    t.Run("Test foo", func(t *testing.T) {
        fmt.Println("Testing foo");
    })
}

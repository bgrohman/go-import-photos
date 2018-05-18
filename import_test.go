package main

import "testing"

func TestExample(t *testing.T) {
    a := 1
    b := 2
    if a + b != 3 {
        t.Error("a + b is not equal to 3")
    }
}

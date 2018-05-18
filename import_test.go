package main

import (
    "path/filepath"
    "os"
    "reflect"
    "testing"
)

func TestExample(t *testing.T) {
    a := 1
    b := 2
    if a + b != 3 {
        t.Error("a + b is not equal to 3")
    }
}

const SOURCE = "./test/source"
const DESTINATION = "./test/destination"

func TestImport(t *testing.T) {
    err := os.RemoveAll(DESTINATION)
    if err != nil {
        t.Error(err)
        t.FailNow()
    }

    operations, err := Import(SOURCE, DESTINATION)
    if err != nil {
        t.Error(err)
    }

    expectedOperations := map[string]string {
        "test/source/test1.jpg": "./test/destination",
        "test/source/test2.NEF": "./test/destination",
        "test/source/test3.jpg": "./test/destination",
        "test/source/test4.jpeg": "./test/destination",
    }

    if reflect.DeepEqual(expectedOperations, operations) != true {
        t.Error("Wrong operations list, %s", operations)
    }

    expectedFiles := []string {filepath.Join(DESTINATION, "2011", "2011-10-03", "test1.jpg"),
                               filepath.Join(DESTINATION, "2010", "2010-09-19", "test2.NEF"),
                               filepath.Join(DESTINATION, "2012", "2012-08-04", "test3.jpg"),
                               filepath.Join(DESTINATION, "2011", "2011-10-05", "test4.jpeg")}

    for _, f := range expectedFiles {
        _, err := os.Stat(f)
        if err != nil {
            t.Error("Missing destination file %s", f)
        }
    }
}

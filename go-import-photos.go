package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "path/filepath"
    "os"
    "strconv"
    "github.com/rwcarlsen/goexif/exif"
)

func help() {
    fmt.Println("Usage:")
    fmt.Println("go-import-photos <source-path> <destination-path>");
}

func copyFile(source, destination string) error {
    in, err := os.Open(source)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(destination)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }

    return out.Close()
}

func importFile(sourceFilePath string, destinationPath string, collisions map[string]int) error {
    file, err := os.Open(sourceFilePath)
    if err != nil {
        log.Fatal(err)
        return err
    }

    exifData, err := exif.Decode(file)
    if err != nil {
        log.Fatal(err)
        return err
    }

    dateTime, err := exifData.DateTime()
    if err != nil {
        log.Fatal(err)
        return err
    }

    // Time.Format reference time: Mon Jan 2 15:04:05 -0700 MST 2006
    year := dateTime.Format("2006")
    month := dateTime.Format("01")
    day := dateTime.Format("02")
    destinationDirectory := filepath.Join(destinationPath, year, year + "-" + month + "-" + day)
    sourceBase := filepath.Base(sourceFilePath)
    destinationFilePath := filepath.Join(destinationDirectory, sourceBase)

    fmt.Println("Copying " + sourceBase + " to " + destinationFilePath)
    err = os.MkdirAll(destinationDirectory, 0777)
    if err != nil {
        log.Fatal(err)
        return err
    }

    if _, err := os.Stat(destinationFilePath); err == nil {
        next := collisions[destinationFilePath] + 1
        collisions[destinationFilePath] = next
        destinationFilePath = destinationFilePath + "_" + strconv.Itoa(next)
        fmt.Println("Found duplicate for " + sourceBase + ", renaming to " + filepath.Base(destinationFilePath))
    }

    err = copyFile(sourceFilePath, destinationFilePath)
    if err != nil {
        log.Fatal(err)
        return err
    }

    return nil
}

func main() {
    args := os.Args[1:]
    if len(args) != 2 {
        help()
        return
    }

    sourcePath, err := filepath.Abs(args[0])
    if err != nil {
        log.Fatal(err)
        return
    }

    destinationPath, err := filepath.Abs(args[1])
    if err != nil {
        log.Fatal(err)
        return
    }

    files, err := ioutil.ReadDir(sourcePath)
    if err != nil {
        log.Fatal(err)
        return
    }

    collisions := make(map[string]int)
    unimportedFiles := make([]string, 0, len(files))

    for _, f := range files {
        sourceFilePath := filepath.Join(sourcePath, f.Name())
        err = importFile(sourceFilePath, destinationPath, collisions)
        if err != nil {
            unimportedFiles = append(unimportedFiles, sourceFilePath)
        }
    }

    if len(unimportedFiles) > 0 {
        fmt.Println("The following files were not imported:")

        for _, f := range unimportedFiles {
            fmt.Println(f)
        }
    }
}

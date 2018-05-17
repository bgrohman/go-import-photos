package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "path/filepath"
    "os"
    "strings"
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

func getPathWithoutCollision(path string, originalPath string, collisions map[string]int) string {
    if _, err := os.Stat(path); err == nil {
        if originalPath == "" {
            originalPath = path
        }
        next := collisions[originalPath] + 1
        collisions[originalPath] = next
        ext := filepath.Ext(originalPath)
        withoutExt := strings.TrimSuffix(originalPath, ext)
        newPath := withoutExt + "_" + strconv.Itoa(next) + ext
        safePath := getPathWithoutCollision(newPath, originalPath, collisions)
        fmt.Println("Found duplicate for " + filepath.Base(originalPath) + ", renaming to " + filepath.Base(safePath))
        return safePath
    }

    return path
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

    destinationFilePath = getPathWithoutCollision(destinationFilePath, "", collisions)

    err = copyFile(sourceFilePath, destinationFilePath)
    if err != nil {
        log.Fatal(err)
        return err
    }

    return nil
}

func checkFatalError(err error) {
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
}

func main() {
    args := os.Args[1:]
    if len(args) != 2 {
        help()
        return
    }

    sourcePath, err := filepath.Abs(args[0])
    checkFatalError(err)

    destinationPath, err := filepath.Abs(args[1])
    checkFatalError(err)

    files, err := ioutil.ReadDir(sourcePath)
    checkFatalError(err)

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

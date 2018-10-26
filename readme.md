# go-import-photos

A command-line utility for importing photos.

## Usage

    go-import-photos <source-path> <destination-path>

Where `source-path` is the path to a directory containing the photos to import
and `destination-path` is the path for storing the photos. Photos are copied,
not moved.

Photos will be stored by year and date using EXIF data from the files according
to the following pattern within the destination directory:

    <destination>/<year>/<year>-<month>-<day/

Example:

    ~/photos/2018/2018-05-17/photo.jpg


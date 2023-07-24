# Fotos

The goal of this project was to publish a 3+ Terabyte large photo archive
with a million JPG files on an online server with just 25 GB drive space.
To achieve this, all images were compressed to 2048x2048 pixels and 400x200
pixels and with a compression that achieves file sizes of &lt;150kb and
&lt;2kb respectively. The folder structure shall remain to allow easy
updating and amending of new albums to the existing structure.

On an average computer, resizing an image takes a fraction of a second.
Measurements have been taken to speed this process up and achieve the
second goal: Accessing galleries should only consume small data quotas
and very fast to load and use.

Thumbnails are saved in JPEG-XL, which requires a browser that supports
this new file format. At the moment, [Thorium](https://thorium.rocks) is
the only widely available chromium-based browser to support JPEG-XL,
unless you got an iPhone, of course.

## Image viewer

The image viewer reads JSON files from the folder and generates a
mobile-friendly preview page. The folder navigation uses encrypted
folder names, which are extracted from the JSON file, ensuring
enhanced security and preventing unauthorized access. With the
implementation of path and file name encryption through hashing,
as well as disabled directory listing, the system provides robust
protection against unauthorized guessing. It is recommended to use
additional security measures such as using a HTTPS-only web server
to keep the folder structure hidden.

**Dependencies (Node.js modules only)**

- minireset.css
- photoswipe

## Preview files generator

The preview files generator walks through a source directory tree,
mirrors the folder structure to the destination path, and for each
pho.to file it encounters, a pho.to.o.jxl, a pho.to.h.jxl and a
pho.to.s.jxl is generated. The generator will place an index.json
file into each directory that contains the date taken, the dimensions
and the name of each file. GPS coordinates and blurry thumbnails
are saved into the json file as well.

This script can start where it had left off and will read existing
index.json files and thumbnails and check their modify dates.

## Getting started

**On Windows,** go-fitz does not compile! Remove its reference in
converter.go by replacing the respective switch case with
`return nil, nil`, from `go.mod` if needed, and run `go mod vendor`

**Dependencies:**

- cgo (needs to be enabled)
- cjxl
- dcraw
- exiftool
- ffmpeg
- golang 1.14.3

**On Windows,** obtain cjxl, dcraw, exiftool and ffmpeg and place
their executables in the `%PATH%`. Enable cgo by changing your
environment variables or `SET CGO_ENABLED=1`

**Set up the vendor folder:**

- Run this. [Why? Read here.](https://github.com/nomad-software/vend)
  ```shell
  go get github.com/nomad-software/vend 
  $GOPATH/bin/vend
  ```
  **On Windows,** you can use the above go get command, followed by
  `go mod vendor` to fix your local folder. Then go into your goroot,
  locate vend and build it. Then jump back to the project and execute
  vend.exe by providing its full path or putting it into the `%PATH%`

**Run:**

- Run `cp config.example.yml config.yml`
- Edit your `config.yml` to suit your needs.
- Compile with `go build fotos/fotos/main`
- Run `./main`.

> ***<span style="color:red !important">This tool can delete folders and files so carefully check your config.yml file</span>***

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

## Image viewer

The image viewer reads JSON files from the folder and generates a
mobile-friendly preview page. Folder navigation will also be
extracted from the JSON file. Beware that this plain implementation
neither contains a password protection nor any dynamic authentication.
If you want a password protection, use the HTTP Basic Authentication
together with a HTTPS-only web server.

**To prevent lags or slow UI response on mobile devices:**

- Albums should not contain more than 1.500 pictures
- Full screen images should not be larger than 1024 pixels
  in any dimension

**Dependencies (Node.js modules only)**

- minireset.css
- javascript-flex-images
- photoswipe

## Preview files generator

The preview files generator walks through a source directory tree,
mirrors the folder structure to the destination path, and for each
pho.to file it encounters, a pho.to.h.jxl and a pho.to.s.jxl is
generated. The generator will place an index.json file into each
directory that contains the date taken, the dimensions and the name
of each file. GPS coordinates and averaged edge colors are generated.

This script can start where it had left off and will read existing
index.json files and thumbnails and check their modify dates.

## REST Server

This release is baked into a REST server that is listening to requests
to create and delete files programmatically based on query parameters.

### /pictures/add
**Using GET:**

Available query parameters:

<dl>
  <dt><strong>token (Required!)</strong></dt>
  <dd>A pre-shared key to authenticate a discord bot against the server</dd>
  <dt><strong>id (Required!)</strong></dt>
  <dd>64-bit unsigned integer id field for the discord user id</dd>
  <dt><strong>username (Required!)</strong></dt>
  <dd>Username to be displayed in the image viewer</dd>
  <dt><strong>discriminator (Required!)</strong></dt>
  <dd>Discriminator to be displayed in the image viewer</dd>
  <dt><strong>url (Required!)</strong></dt>
  <dd>Path to a file to be downloaded into the discord user's gallery</dd>
</dl>

For example: Calling the following URL
will add a bird into the Björn Eberhardt's gallery.

    /pictures/add?token=TOP-SECRET&id=215568977756291072&username=Bj%C3%B6rn Eberhardt&discriminator=2964&url=https://i.imgur.com/XsTeOCT.jpg

**Using POST:**

Refer to this example object to construct your request object:

```json
{
  "userId": "215568977756291072",
  "userName": "Björn Eberhardt",
  "discriminator": "2964",
  "gallery": "215568977756291072",
  "url": "https://i.imgur.com/XsTeOCT.jpg",
  "preSharedKey": "TOP-SECRET"
}
```

### /pictures/del

This command will delete a picture. If a gallery contains no
files, the entire gallery will be deleted.

### /pictures/random

This command takes no arguments and returns a JSON with data
about a randomly picked image. Example:

```json
{
  "userId": "215568977756291072/",
  "userName": "Björn Eberhardt",
  "discriminator": "2964",
  "gallery": "215568977756291072/",
  "filename": "XsTeOCT-abcde.jpg"
}
```

## Getting started

**Dependencies:**

- cgo (needs to be enabled)
- dcraw
- exiftool
- ffmpeg
- golang 1.14.3
- all the other letters in the alphabet

**Set up the vendor folder:**

- Run this. [Why? Read here.](https://github.com/nomad-software/vend)
  ```shell
  go get github.com/nomad-software/vend 
  $GOPATH/bin/vend
  ```

**Run:**

- Run `cp config.example.yml config.yml`
- Edit your `config.yml` to suit your needs.
- Compile with `go build fotos/server/main`
- Run `./main --help`.

> ***<span style="color:red !important">This tool can delete folders and files so carefully check your program arguments</span>***

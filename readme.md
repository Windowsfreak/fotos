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
pho.to file it encounters, a pho.to.h.webp and a pho.to.s.webp is
generated. The generator will place an index.json file into each
directory that contains the date taken, the dimensions and the name
of each file. GPS coordinates and averaged edge colors are generated.

This script can start where it had left off and will read existing
index.json files and thumbnails and check their modify dates.

**Dependencies:**

- cgo (needs to be enabled)
- dcraw
- exiftool
- ffmpeg
- golang 1.14.3
- all the other letters in the alphabet

**Run:**

- Compile with `./go build`
- Run `./fotos --help`.

> ***<span style="color:red !important">This tool can delete folders and files so carefully check your program arguments</span>***

# Fotos

The goal of this project was to publish a 3+ Terabyte large photo archive
with a million JPG files on an online server with just 25 GB drive space.
To achieve this, all images were compressed to 800x800 pixels and 200x200
pixels and with a compression that achieves file sizes of &lt;50kb and
&lt;5kb respectively. The folder structure shall remain to allow easy
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

**Known bugs**

- The PHP version will not show a proper error message when a
  folder doesn't contain a JSON file
- The pure HTML/JS version takes much time to show more than
  1.000 pictures

## Preview files generator

The preview files generator walks through a source directory tree,
mirrors the folder structure to the destination path, and for each
.jpg file it encounters, a .jpg.thumb.jpg and a .jpg.small.jpg is
generated. A SQL database is used to keep track of the progress,
so when an error or timeout occurs, the script can continue where
it has stopped. The generator will also place an index.json file
into each directory that contains the date taken, the dimensions
and the name of each file.

The commands executed are designed for Windows, but it is possible
to change them to work on Linux as well.

**Dependencies:**

- ImageMagick

**Known bugs:**

- PHP < 7.1 does not support Non-ASCII7-Characters very well.
  Using them, the shell scripts won't find the files
  that were retrieved using `scandir`.
- Images that use an EXIF Orientation are not properly rotated.
  Preview images and dimensions are based on the data structure

**Possible Improvements:**

- Read back data from JSON to prevent identifying files where
  enough information is already present.
- Different expiry dates for images and folders
- Generate .webp instead of .jpg files
- Multi-Threading
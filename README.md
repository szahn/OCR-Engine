# OCR Engine

# Rasterizer
Takes a document from any source (Web, PDF) and renders it into an image.

To rasterize PDFs, Ghostscript is required.

### Ghostscript Installation (Linux):

Git clone Ghostscript source into the parent folder:

`git clone git://git.ghostscript.com/ghostpdl.git ghostpdl`

Run `make` to build the source and then `make so` to build the libraries.

# OCR
Converts images to text.

###Terreract OCR Engine

Git clone

`git clone https://github.com/tesseract-ocr/tesseract.git tesseract`

# Parser
Partitions text into boundaries. The result is some kind of document with structure.

# GRPC Server

Server/client model for sending documents to be ocr-ed
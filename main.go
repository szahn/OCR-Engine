package main

import (
	"fmt"
	"log"
	"os"

	"ocr-engine/ocr"
	"ocr-engine/parser"
	"ocr-engine/rasterize"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing filename argument")
		return
	}

	filename := os.Args[1]
	doc := rasterize.NewSourceDocument(filename)
	rasterOpts := &rasterize.DocumentRasterizeOptions{rasterize.EngineGhostScriptLib}

	image, err := doc.ToImage(rasterOpts)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Rasterized %vkb\n", (len(image) / 1204))

	hocr, err := ocr.GetHOCR(image)
	if err != nil {
		log.Printf("Error OCRing file %s: %v\n", filename, err)
		return
	}

	blocks := parser.Parse(hocr)

	for _, b := range blocks {
		fmt.Printf("%s\n", parser.DescribeBlock(b))
	}

}

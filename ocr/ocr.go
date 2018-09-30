package ocr

import (
	"github.com/otiai10/gosseract"
)

func GetHOCR(image []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(image)
	hocr, err := client.HOCRText()
	return hocr, err
}

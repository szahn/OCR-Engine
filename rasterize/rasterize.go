package rasterize

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"unsafe"

	"ocr-engine/random"
)

/*
#cgo CFLAGS: -I../../ghostpdl/base/
#cgo CFLAGS: -I../../ghostpdl/psi/
#cgo LDFLAGS: -L../../ghostpdl/sobin -lgs
#include "ierrors.h"
#include "iapi.h"

typedef int (GSDLLAPIPTR PFN_gsapi_new_instance)(void **pinstance, void *caller_handle);

typedef int (GSDLLAPIPTR PFN_gsapi_set_arg_encoding)(void *instance, int encoding);

typedef int (GSDLLAPIPTR PFN_gsapi_init_with_args)(void *instance, int argc, char **argv);

typedef int (GSDLLAPIPTR PFN_gsapi_exit)(void *instance);

typedef void (GSDLLAPIPTR PFN_gsapi_delete_instance)(void *instance);

int convert(char **gsargv, int gsargc)
{
	void *minst;

    int code, code1;
    code = gsapi_new_instance(&minst, NULL);
    if (code < 0)
	return 1;
    code = gsapi_set_arg_encoding(minst, 1);
    if (code == 0)
        code = gsapi_init_with_args(minst, gsargc, gsargv);
    code1 = gsapi_exit(minst);
    if ((code == 0) || (code == -101))
	code = code1;

    gsapi_delete_instance(minst);

    if ((code == 0) || (code == -101))
	return 0;
    return 1;
}
*/
import "C"

const (
	srcUnknown = 0
	srcPDF     = 1
	srcURL     = 2
)

const (
	EngineGhostScriptShell = 0
	EngineGhostScriptLib   = 1
)

type GhostScriptOptions struct {
	dpi    int
	device string
}

type SourceDocument struct {
	source     string
	sourceType int
}

type DocumentRasterizeOptions struct {
	PdfEngine int
}

func (r *SourceDocument) ToImage(rasterizeOpts *DocumentRasterizeOptions) ([]byte, error) {

	switch r.sourceType {
	case srcPDF:
		{
			switch rasterizeOpts.PdfEngine {
			case EngineGhostScriptLib:
				{
					gsOpts := GhostScriptOptions{300, "png16m"}
					ext := "png"
					destPath := "temp"
					os.Mkdir(destPath, os.ModePerm)

					destFilename := fmt.Sprintf("%s/%s.%s", destPath, random.NextRandom(), ext)
					args := []string{"ps2pdf", "-dSAFER", "-dBATCH", "-dNOPAUSE", fmt.Sprintf("-sDEVICE=%s", gsOpts.device), fmt.Sprintf("-r%v", gsOpts.dpi), fmt.Sprintf("-sOutputFile=%s", destFilename), r.source}

					gsargc := C.int(len(args))

					arr := (**C.char)(C.malloc(C.size_t(C.int(len(args)))))
					view := (*[1 << 30]*C.char)(unsafe.Pointer(arr))[0:len(args):len(args)]
					for i, x := range args {
						view[i] = C.CString(x)
					}

					result := C.convert(arr, gsargc)
					if result != 0 {
						return nil, errors.New(fmt.Sprintf("Error code %v", result))
					}

					image, readErr := ioutil.ReadFile(destFilename)
					if readErr != nil {
						return nil, readErr
					}

					removeErr := os.Remove(destFilename)
					if removeErr != nil {
						log.Printf("Error deleting temp file %s: %v", destFilename, removeErr)
					}

					return image, nil
				}
			case EngineGhostScriptShell:
				{
					gsOpts := GhostScriptOptions{300, "png16m"}
					ghostScriptPath := "gs"
					ext := "png"
					destPath := "temp"
					os.Mkdir(destPath, os.ModePerm)

					destFilename := fmt.Sprintf("%s/%s.%s", destPath, random.NextRandom(), ext)
					args := []string{"-dSAFER", "-dBATCH", "-dNOPAUSE", fmt.Sprintf("-sDEVICE=%s", gsOpts.device), fmt.Sprintf("-r%v", gsOpts.dpi), fmt.Sprintf("-sOutputFile=%s", destFilename), r.source}
					cmd := exec.Command(ghostScriptPath, args...)
					out, gsErr := cmd.CombinedOutput()
					if gsErr != nil {
						log.Println(string(out))
						return nil, gsErr
					}

					image, readErr := ioutil.ReadFile(destFilename)
					if readErr != nil {
						return nil, readErr
					}

					removeErr := os.Remove(destFilename)
					if removeErr != nil {
						log.Printf("Error deleting temp file %s: %v", destFilename, removeErr)
					}

					return image, nil
				}
			}
		}
	default:
		{
			break
		}
	}

	return nil, errors.New("unknown source type")
}

func NewSourceDocument(docSource string) *SourceDocument {
	docType := srcUnknown

	if strings.HasSuffix(strings.ToLower(docSource), ".pdf") {
		docType = srcPDF
	}

	doc := &SourceDocument{docSource, docType}
	return doc
}

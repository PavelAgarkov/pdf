package service

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PdfAdapter struct{}

func NewPdfAdapter() *PdfAdapter {
	return &PdfAdapter{}
}

func (pdfFacade *PdfAdapter) MergeFiles(inFiles []string, outFile string) error {

	//out := "./internal/service/files/new.pdf"

	//list := []string{"./internal/service/files/file3.pdf", "./internal/service/files/File0003.pdf", "./internal/service/files/File0004.pdf"}

	err := api.MergeAppendFile(inFiles, outFile, nil)
	if err != nil {
		return fmt.Errorf("can't merge in %s file: %w", outFile, err)
		//fmt.Println("own")
	}
	return nil
}

func (pdfFacade *PdfAdapter) SplitFile(inFile, outDir string) error {
	err := api.SplitFile(inFile, outDir, 0, nil)
	if err != nil {
		return fmt.Errorf("can't split file %s: %w", inFile, err)
	}
	return nil
}

func (pdfFacade *PdfAdapter) CutFile(inFile, outDir, outFile string, selectedPages []string) error {
	err := api.CutFile(inFile, outDir, outFile, selectedPages, nil, nil)
	if err != nil {
		return fmt.Errorf("can't cut file %s: %w", inFile, err)
	}
	return nil
}

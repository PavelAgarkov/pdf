package adapter

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

const (
	PdfAlias = "pdf"
)

type PdfAdapter struct{}

func NewPdfAdapter() *PdfAdapter {
	return &PdfAdapter{}
}

func (pdfAdapter *PdfAdapter) GetAlias() string {
	return PdfAlias
}

func (pdfAdapter *PdfAdapter) MergeFiles(inFiles []string, outFile string) error {
	err := api.MergeAppendFile(inFiles, outFile, nil)
	if err != nil {
		return fmt.Errorf("can't merge in %s file: %w", outFile, err)
	}
	return nil
}

func (pdfAdapter *PdfAdapter) SplitFile(inFile, outDir string) error {
	err := api.SplitFile(inFile, outDir, 0, nil)
	if err != nil {
		return fmt.Errorf("can't split file %s: %w", inFile, err)
	}
	return nil
}

func (pdfAdapter *PdfAdapter) CutFile(inFile, outDir, outFile string, selectedPages []string) error {
	err := api.CutFile(inFile, outDir, outFile, selectedPages, nil, nil)
	if err != nil {
		return fmt.Errorf("can't cut file %s: %w", inFile, err)
	}
	return nil
}

func (pdfAdapter *PdfAdapter) RemovePagesFile(inFile, outFile string, selectedPages []string) error {
	err := api.RemovePagesFile(inFile, outFile, selectedPages, nil)
	if err != nil {
		return fmt.Errorf("can't remove pages from file %s: %w", inFile, err)
	}
	return nil
}

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
	err := api.SplitFile(inFile, outDir, 1, nil)
	if err != nil {
		return fmt.Errorf("can't split file %s: %w", inFile, err)
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

func (pdfAdapter *PdfAdapter) Optimize(inFile, outFile string) error {
	err := api.OptimizeFile(inFile, outFile, nil)
	if err != nil {
		return fmt.Errorf("can't optimize file %s: %w", inFile, err)
	}
	return nil
}

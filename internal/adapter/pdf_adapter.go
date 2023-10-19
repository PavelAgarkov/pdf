package adapter

import (
	"context"
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

func (pdfAdapter *PdfAdapter) MergeFiles(ctx context.Context, inFiles []string, outFile string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := api.MergeAppendFile(inFiles, outFile, nil)
	if err != nil {
		return fmt.Errorf("can't merge in file: %w", err)
	}
	return nil
}

func (pdfAdapter *PdfAdapter) SplitFile(ctx context.Context, inFile, outDir string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := api.SplitFile(inFile, outDir, 1, nil)
	if err != nil {
		return fmt.Errorf("can't split file: %w", err)
	}
	return nil
}

func (pdfAdapter *PdfAdapter) RemovePagesFile(ctx context.Context, inFile, outFile string, selectedPages []string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := api.RemovePagesFile(inFile, outFile, selectedPages, nil)
	if err != nil {
		return fmt.Errorf("can't remove pages from file: %w", err)
	}
	return nil
}

func (pdfAdapter *PdfAdapter) Optimize(ctx context.Context, inFile, outFile string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := api.OptimizeFile(inFile, outFile, nil)
	if err != nil {
		return fmt.Errorf("can't optimize file: %w", err)
	}
	return nil
}

func (pdfAdapter *PdfAdapter) PageCount(ctx context.Context, inFile string) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	pageCount, err := api.PageCountFile(inFile)
	if err != nil {
		return 0, err
	}
	return pageCount, nil
}

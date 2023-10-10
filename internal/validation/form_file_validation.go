package validation

import (
	"errors"
	"fmt"
	"mime/multipart"
	"pdf/internal"
	"slices"
	"strconv"
	"strings"
)

func FormFileValidation(form *multipart.Form) error {
	var sumSize int64 = 0
	countFiles := 0
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			sumSize += fileHeader.Size
			countFiles++
		}
	}

	archiveFormatSlice, ok := form.Value[internal.ArchiveFormatKeyForRequest]
	archiveFormat := ""
	if len(archiveFormatSlice) > 0 {
		archiveFormat = archiveFormatSlice[0]
	}

	formats := []string{
		internal.ZipFormat,
		internal.ZipZstFormat,
		internal.TarFormat,
		internal.TarGzFormat,
	}

	if !ok || !slices.Contains(formats, archiveFormat) {
		return errors.New("archive format must be selected from the list")
	}

	if sumSize > internal.MaxSumUploadFilesSizeByte {
		return errors.New("upload files must be less 100Mb")
	}

	if countFiles > internal.MaxNumberUploadFiles {
		return errors.New("number upload files must be less 100")
	}

	return nil
}

func AlphaSymbolValidation(form *multipart.Form, key string) error {
	const alpha = "1234567890-"

	for _, interval := range form.Value[key] {
		for _, char := range interval {
			if !strings.Contains(alpha, strings.ToLower(string(char))) {
				return errors.New(fmt.Sprintf("invalid symbol, !%s!", string(char)))
			}
		}

		chunks := strings.Split(interval, "-")
		if len(chunks) > 2 {
			return errors.New("format must be some '2-5' or 5")
		}
	}
	return nil
}

func OrderIntervalValidation(form *multipart.Form, key string) error {
	_, intervals := internal.ParseIntervals(form.Value[key])
	for _, interval := range intervals {
		if len(interval) == 2 {
			if interval[0] > interval[1] {
				return errors.New("interval format 'n-n' must be written in ascending order")
			}
			if interval[0] == 0 || interval[1] == 0 {
				return errors.New("interval format 'n-n' must be written in ascending order, not zero")
			}
		}
	}
	return nil
}

func NumberFilesValidation(form *multipart.Form, must int) error {
	number := 0
	for _, fileHeaders := range form.File {
		number = len(fileHeaders)
		break
	}
	if number != must {
		return errors.New(fmt.Sprintf("for split operation must %s pdf files", strconv.Itoa(must)))
	}
	return nil
}

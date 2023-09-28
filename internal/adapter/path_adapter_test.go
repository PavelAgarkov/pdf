package adapter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_remove_pages(t *testing.T) {
	pa := NewPathAdapter()

	expected := "ServiceAgreement_template.pdf"
	resource := filepath.FromSlash("./files/ServiceAgreement_template.pdf")
	path, last, err := pa.StepBack(Path(resource))

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(last)
	assert.NotEqual(t, Path(resource), path)
	assert.Equal(t, expected, last)

	path, _, err = pa.StepBack(path)
	_, _, err = pa.StepBack(path)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func Test_add_pages(t *testing.T) {
	pa := NewPathAdapter()

	resource := "ServiceAgreement_template.pdf"
	old := filepath.FromSlash("./files/")
	expected := filepath.FromSlash("./files/ServiceAgreement_template.pdf")
	path, newPath, err := pa.StepForward(Path(old), resource)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(newPath)
	assert.NotEqual(t, Path(resource), path)
	assert.Equal(t, Path(expected), newPath)

	if err != nil {
		fmt.Println(err.Error())
	}
}

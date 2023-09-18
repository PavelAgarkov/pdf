package adapter

type FileAdapter struct{}

const (
	FileAlias = "file"
)

func NewFileAdapter() *FileAdapter {
	return &FileAdapter{}
}

func (fa *FileAdapter) GetAlias() string {
	return FileAlias
}

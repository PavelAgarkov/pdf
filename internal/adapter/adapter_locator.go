package adapter

type Adapter interface {
	GetAlias() string
}

type Locator struct {
	fileAdapter *FileAdapter
	faAlias     string
	pathAdapter *PathAdapter
	paAlias     string
	pdfAdapter  *PdfAdapter
	pdfaAlias   string
}

func NewAdapterLocator(fileAdapter *FileAdapter, pathAdapter *PathAdapter, pdfAdapter *PdfAdapter) *Locator {
	return &Locator{
		fileAdapter: fileAdapter,
		faAlias:     fileAdapter.GetAlias(),
		pathAdapter: pathAdapter,
		paAlias:     pathAdapter.GetAlias(),
		pdfAdapter:  pdfAdapter,
		pdfaAlias:   pdfAdapter.GetAlias(),
	}
}

func (l *Locator) Locate(alias string) Adapter {
	switch alias {
	case FileAlias:
		return l.fileAdapter
	case PathAlias:
		return l.pathAdapter
	case PdfAlias:
		return l.pdfAdapter
	default:
		return nil
	}
}

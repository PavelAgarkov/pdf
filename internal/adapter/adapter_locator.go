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
	rarAdapter  *RarAdapter
	rarAlias    string
}

func NewAdapterLocator(fileAdapter *FileAdapter, pathAdapter *PathAdapter, pdfAdapter *PdfAdapter, rarAdapter *RarAdapter) *Locator {
	return &Locator{
		fileAdapter: fileAdapter,
		faAlias:     fileAdapter.GetAlias(),
		pathAdapter: pathAdapter,
		paAlias:     pathAdapter.GetAlias(),
		pdfAdapter:  pdfAdapter,
		pdfaAlias:   pdfAdapter.GetAlias(),
		rarAdapter:  rarAdapter,
		rarAlias:    rarAdapter.GetAlias(),
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
	case RarAlias:
		return l.rarAdapter

	default:
		return nil
	}
}

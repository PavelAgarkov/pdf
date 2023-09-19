package adapter

type RarAdapter struct{}

const (
	RarAlias = "rar"
)

func NewRarAdapterAdapter() *RarAdapter {
	return &RarAdapter{}
}

func (fa *RarAdapter) GetAlias() string {
	return FileAlias
}

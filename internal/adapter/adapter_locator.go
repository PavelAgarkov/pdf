package adapter

type Adapter interface {
	GetAlias() string
}

type Locator struct {
	locateMap map[string]Adapter
}

func NewAdapterLocator(adapters ...Adapter) *Locator {
	l := &Locator{
		locateMap: make(map[string]Adapter),
	}

	for _, adapter := range adapters {
		l.locateMap[adapter.GetAlias()] = adapter
	}
	return l
}

func (l *Locator) setAdapter(alias string, adapter Adapter) *Locator {
	l.locateMap[alias] = adapter
	return l
}

func (l *Locator) Locate(alias string) Adapter {
	adapter, ok := l.locateMap[alias]

	if !ok {
		return nil
	}

	return adapter
}

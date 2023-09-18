package pdf_operation

type OperationConfiguration struct {
	splitIntervals       []string
	cutIntervals         []string
	removePagesIntervals []string
	mergeOrder           []string
}

func NewOperationConfiguration(
	splitIntervals,
	cutIntervals,
	removePagesIntervals,
	mergeOrder []string,
) *OperationConfiguration {
	return &OperationConfiguration{}
}

func (oc *OperationConfiguration) GetSplitIntervals() []string {
	return oc.splitIntervals
}

func (oc *OperationConfiguration) GetCutIntervals() []string {
	return oc.cutIntervals
}

func (oc *OperationConfiguration) GetRemovePagesIntervals() []string {
	return oc.removePagesIntervals
}

func (oc *OperationConfiguration) GetMergeOrder() []string {
	return oc.mergeOrder
}

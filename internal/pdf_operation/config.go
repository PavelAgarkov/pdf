package pdf_operation

type OperationConfiguration struct {
	splitIntervals       []string
	removePagesIntervals []string
}

func NewConfiguration(
	splitIntervals,
	removePagesIntervals []string,
) *OperationConfiguration {
	return &OperationConfiguration{
		splitIntervals:       splitIntervals,
		removePagesIntervals: removePagesIntervals,
	}
}

func (oc *OperationConfiguration) GetSplitIntervals() []string {
	return oc.splitIntervals
}

func (oc *OperationConfiguration) GetRemovePagesIntervals() []string {
	return oc.removePagesIntervals
}

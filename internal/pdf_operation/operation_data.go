package pdf_operation

import (
	"pdf/internal/adapter"
	"pdf/internal/entity"
)

type OperationDataInterface interface {
	CanDeleted() bool
	GetUserData() *entity.UserData
	SetStatus(status OperationStatus) *OperationData
}

const (
	StatusStarted          = "started"
	StatusProcessed        = "processed"
	StatusCompleted        = "completed"
	StatusExpired          = "expired"
	StatusCanceled         = "canceled"
	StatusAwaitingDownload = "awaiting_download"
)

type OperationData struct {
	ud            *entity.UserData
	archiveDir    adapter.ArchiveDir
	status        OperationStatus //статус операции нужен для контоля отмены токена и очистки памяти
	stoppedReason StoppedReason
}

func NewOperationData(
	ud *entity.UserData,
	archiveDir adapter.ArchiveDir,
	status OperationStatus,
	stoppedReason StoppedReason,
) *OperationData {
	return &OperationData{
		ud:            ud,
		archiveDir:    archiveDir,
		status:        status,
		stoppedReason: stoppedReason,
	}
}

func (od *OperationData) GetUserData() *entity.UserData {
	return od.ud
}

func (od *OperationData) GetArchiveDir() adapter.ArchiveDir {
	return od.archiveDir
}

func (od *OperationData) GetStatus() OperationStatus {
	return od.status
}

func (od *OperationData) GetStoppedReason() StoppedReason {
	return od.stoppedReason
}

func (od *OperationData) CanDeleted() bool {
	return od.status == StatusExpired ||
		od.status == StatusCanceled ||
		od.status == StatusCompleted
}

func (od *OperationData) SetStatus(status OperationStatus) *OperationData {
	od.status = status
	return od
}

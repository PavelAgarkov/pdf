package pdf_operation

import (
	"pdf/internal"
	"pdf/internal/entity"
)

type OperationDataInterface interface {
	CanDeleted() bool
	GetUserData() *entity.UserData
	SetStatus(status internal.OperationStatus) *OperationData
}

//		 		expired
//		   		/\
//	       		|
//
// started->processed->awaiting_download->completed
//
//	   			|
//	  			\/
//			canceled

type OperationData struct {
	ud            *entity.UserData
	archivePath   internal.ArchiveDir
	status        internal.OperationStatus //статус операции нужен для контоля отмены токена и очистки памяти
	stoppedReason internal.StoppedReason
}

func NewOperationData(
	ud *entity.UserData,
	archiveDir internal.ArchiveDir,
	status internal.OperationStatus,
	stoppedReason internal.StoppedReason,
) *OperationData {
	return &OperationData{
		ud:            ud,
		archivePath:   archiveDir,
		status:        status,
		stoppedReason: stoppedReason,
	}
}

func (od *OperationData) GetUserData() *entity.UserData {
	return od.ud
}

func (od *OperationData) GetArchivePath() internal.ArchiveDir {
	return od.archivePath
}

func (od *OperationData) GetStatus() internal.OperationStatus {
	return od.status
}

func (od *OperationData) GetStoppedReason() internal.StoppedReason {
	return od.stoppedReason
}

func (od *OperationData) CanDeleted() bool {
	return od.status == internal.StatusExpired ||
		od.status == internal.StatusCanceled ||
		od.status == internal.StatusCompleted
}

func (od *OperationData) SetStatus(status internal.OperationStatus) *OperationData {
	od.status = status
	return od
}

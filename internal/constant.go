package internal

import "time"

const (
	ZipFormat    = ".zip"
	ZipZstFormat = ".zip.zst"
	TarFormat    = ".tar"
	TarGzFormat  = ".tar.gz"

	AuthenticationKey = "X-HASH"

	Salt = "1af1dfa857bf1d8814fe1af898 3c18080019922e557f15a8a"

	LogDir      = "./log"
	PanicLog    = "./log/panic.log"
	ErrLog      = "./log/error.log"
	WarningLog  = "./log/warning.log"
	InfoLog     = "./log/info.log"
	FrontendLog = "./log/frontend.log"

	Timer5  = 5 * time.Minute
	Timer10 = 10 * time.Minute
	Timer15 = 15 * time.Minute

	StatusStarted          = "started"
	StatusProcessed        = "processed"
	StatusCompleted        = "completed"
	StatusExpired          = "expired"
	StatusCanceled         = "canceled"
	StatusAwaitingDownload = "awaiting_download"

	FrontendDist   = "./pdf-frontend/dist"
	FilesPath      = "./files/"
	FaviconFile    = "./pdf-frontend/dist/favicon.ico"
	FrontendAssets = "./pdf-frontend/dist/assets/"
)

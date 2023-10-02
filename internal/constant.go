package internal

import "time"

const (
	ZipFormat    = ".zip"
	ZipZstFormat = ".zip.zst"
	TarFormat    = ".tar"
	TarGzFormat  = ".tar.gz"

	AuthenticationHeader = "Authorization"
	Bearer               = "Bearer"
	BearerSeparator      = "__"

	Salt = "1af1dfa857bf1d8814fe1af898 3c18080019922e557f15a8a"

	Timer5  = 5
	Timer10 = 10
	Timer15 = 15

	Minute = time.Minute

	StatusStarted          = "started"
	StatusProcessed        = "processed"
	StatusCompleted        = "completed"
	StatusExpired          = "expired"
	StatusCanceled         = "canceled"
	StatusAwaitingDownload = "awaiting_download"

	// для совместимости ОС filepath.FromSlash()
	LogDir      = "./log"
	PanicLog    = "./log/panic.log"
	ErrLog      = "./log/error.log"
	WarningLog  = "./log/warning.log"
	InfoLog     = "./log/info.log"
	FrontendLog = "./log/frontend.log"

	FrontendDist   = "./pdf-frontend/dist"
	FilesPath      = "./files/"
	FaviconFile    = "./pdf-frontend/dist/favicon.ico"
	FrontendAssets = "./pdf-frontend/dist/assets/"
)

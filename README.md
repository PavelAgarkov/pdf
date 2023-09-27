Логирование стека вызовов - loggerFactory.GetLogger(logger.ErrorName).With(zap.Stack("stackTrace")).Error("errror")


cat /proc/sys/fs/file-max --- файловые дескрипторы
lsof | wc -l ---- открытые
lsof | grep pdf/ | wc -l

анализ трейса
f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()

go tool trace trace.out

Запуск проекта из директории бэкенда.

1. docker compose build
2. docker compose up

это для обновления компилированных изменений
3. docker exec -it node-local sh -----> npm run build
4. docker exec -it pdf-pdf-1 sh -----> go run main.go

GOOS=windows GOARCH=amd64 go build -buildvcs=false
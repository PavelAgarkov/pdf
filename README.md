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
4. docker exec -it pdf-pdf-1 sh -----> 1.go install, 2.go mod vendor, 3.go mod tidy, 4.go run main.go


открыть термилнал Ubuntu в Windows и пересобирать проект этой командой
docker exec node-local npm run build && docker exec pdf-pdf-1 go build -race -buildvcs=false && docker exec -it pdf-pdf-1 ./pdf

GOOS=windows GOARCH=amd64 go build -buildvcs=false
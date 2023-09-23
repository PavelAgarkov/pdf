Логирование стека вызовов - loggerFactory.GetLogger(logger.ErrorName).With(zap.Stack("stackTrace")).Error("errror")


анализ трейса
f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()

go tool trace trace.out
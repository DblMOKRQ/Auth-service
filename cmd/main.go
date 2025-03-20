package main

// stack: zap, grpc,sqlite, migration
import (
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/DblMOKRQ/auth-service/internal/config"
	"github.com/DblMOKRQ/auth-service/internal/repository"
	"github.com/DblMOKRQ/auth-service/internal/service"
	"github.com/DblMOKRQ/auth-service/internal/token"
	rout "github.com/DblMOKRQ/auth-service/internal/transport"
	"github.com/DblMOKRQ/auth-service/pkg/logger"

	"github.com/DblMOKRQ/auth-service/internal/storage/sqlite"
)

func main() {

	logger := logger.NewLogger()
	logger.Info("Starting auth-service")

	cfg := config.MustLoad()
	_ = cfg

	db, err := sqlite.NewStorage()

	if err != nil {
		logger.Fatal("Database connection error", zap.Error(err))
	}

	repo := repository.NewRepository(db)

	t, err := token.NewJWTMaker(cfg.Token.SecretKey, cfg.Token.ExpirationTime)
	if err != nil {
		logger.Fatal("Token error", zap.Error(err))
	}

	service := service.NewService(repo, logger, t)

	r := rout.NewRouter(grpc.NewServer(), service)
	addr := "0.0.0.0:50051"
	logger.Info("Server started", zap.String("address", addr))

	go func() {
		if err := r.Run(addr); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	<-quit
	logger.Info("Shutting down server")
	r.GracfulShutdown()
	logger.Info("Server stopped")
	/*
		Update Idea:
		1. Сделать фукцнию забыли пароль
		TODO:
		1.Доделать базу данных и protobuf
		2. Написать тесты
				2.1 Написать тесты на repository
		 	2.2 Написать тесты на service
				2.3 Написать тесты на token

		3. Написать документацию
		4. Сделать миграцию БД
		5. Сделать логгер через context
		7. Сделать CI/CD
		8. Сделать gracful shutdown
		Bugs:
		БД закрывается сама по себе
	*/
}

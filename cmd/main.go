package main

// stack: zap, grpc,sqlite
import (
	"github.com/DblMOKRQ/auth-service/internal/config"
	"github.com/DblMOKRQ/auth-service/internal/repository"
	"github.com/DblMOKRQ/auth-service/internal/service"
	"github.com/DblMOKRQ/auth-service/internal/token"
	rout "github.com/DblMOKRQ/auth-service/internal/transport"
	"github.com/DblMOKRQ/auth-service/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

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
	if err := r.Run(addr); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}

	// TODO:
	// 1.Доделать базу данных и protobuf
	// 2. Написать тесты
	// 3. Написать документацию
	// 4. Сделать миграцию БД
	// 5. Сделать логгер через context
	// 6. Запихнуть все в докер +
	// 7. Сделать CI/CD
	// 8. Сделать gracful shutdown
	// Bugs:
	// БД закрывается сама по себе
}

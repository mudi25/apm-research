package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"research-apm/pkg/database/gormx"
	"research-apm/pkg/database/gormx/dialector/apm/mysql"
	"research-apm/pkg/database/gormx/dialector/apm/pgsql"
	"research-apm/pkg/database/gormx/dialector/apm/sqlserver"
	"research-apm/pkg/database/mongox"
	"research-apm/pkg/database/redisx"
	"research-apm/pkg/ginx"
	"research-apm/pkg/tracer"
	"research-apm/services/api/internal/delivery"
	"research-apm/services/api/internal/repository"
	"research-apm/services/api/internal/service"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type App struct {
	Server      *http.Server
	MongoClient *mongo.Client
	DBMessage   *gorm.DB
	DBClientDO  *gorm.DB
	DBProfil    *gorm.DB
}

func (a *App) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		fmt.Println(err.Error())
	}
	mongox.Disconnect(a.MongoClient)
	gormx.Disconnect(a.DBMessage)
	gormx.Disconnect(a.DBClientDO)
	gormx.Disconnect(a.DBProfil)
}

func NewApp(ctx context.Context) (*App, error) {
	if err := tracer.InitTracer(tracer.Config{
		Env:            os.Getenv("ENV"),
		ServiceName:    os.Getenv("SERVICE_NAME"),
		Version:        os.Getenv("SERVICE_VERSION"),
		ServerUrl:      os.Getenv("APM_SERVER_URL"),
		SecretToken:    os.Getenv("APM_SECRET_TOKEN"),
		IsUsingLogging: true,
	}); err != nil {
		return nil, err
	}
	dbClient, err := mongox.NewClient(mongox.Config{
		Uri:             os.Getenv("MONGO_DB_URL"),
		MinPoolSize:     2,
		MaxPoolSize:     5,
		MaxConnIdleTime: 10 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	dbMessage, err := gormx.NewClient(gormx.Config{
		Dialector:  pgsql.NewDialector(os.Getenv("MESSAGE_DB_URL")),
		PoolConfig: nil,
		GormConfig: nil,
	})
	if err != nil {
		return nil, err
	}
	dbClientDo, err := gormx.NewClient(gormx.Config{
		Dialector:  sqlserver.NewDialector(os.Getenv("CLIENT_DO_DB_URL")),
		PoolConfig: nil,
		GormConfig: nil,
	})
	if err != nil {
		return nil, err
	}
	dbProfil, err := gormx.NewClient(gormx.Config{
		Dialector:  mysql.NewDialector(os.Getenv("PROFIL_DB_URL")),
		PoolConfig: nil,
		GormConfig: nil,
	})
	if err != nil {
		return nil, err
	}
	dbRedis, err := redisx.NewClient(ctx, redisx.Config{
		Url:    os.Getenv("REDIS_URL"),
		UseApm: true,
	})
	if err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile("./log.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	engine := ginx.NewEngine(
		ginx.WithTraceID(),
		ginx.WithLogFile(ctx, logFile, ginx.LogConfig{
			AppName:      os.Getenv("SERVICE_NAME"),
			AppSite:      "",
			AppEnv:       os.Getenv("ENV"),
			AppVersion:   os.Getenv("SERVICE_VERSION"),
			AppDBVersion: "",
		}),
		ginx.WithElasticAPM(),
	)
	server := delivery.NewDelivery(
		engine,
		service.NewService(repository.NewRepository(dbClient, dbMessage, dbClientDo, dbProfil, dbRedis)),
	)
	return &App{
		Server:      server,
		MongoClient: dbClient,
		DBMessage:   dbMessage,
		DBClientDO:  dbClientDo,
		DBProfil:    dbProfil,
	}, nil

}

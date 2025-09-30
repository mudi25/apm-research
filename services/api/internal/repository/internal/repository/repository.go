package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"research-apm/pkg/database/mongox"
	"research-apm/pkg/tracer"
	"research-apm/services/api/internal/entity"
	"research-apm/services/api/internal/repository/internal/model"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/gorm"
)

type Repository struct {
	dbUser     *mongo.Collection
	dbMessage  *gorm.DB
	dbClientDO *gorm.DB
	dbProfil   *gorm.DB
	dbRedis    *redis.Client
}

func NewRepository(mongoClient *mongo.Client, dbMessage *gorm.DB, dbClientDO *gorm.DB, dbProfil *gorm.DB, dbRedis *redis.Client) *Repository {
	return &Repository{
		dbUser:     mongoClient.Database("research_apm").Collection("user"),
		dbMessage:  dbMessage,
		dbClientDO: dbClientDO,
		dbProfil:   dbProfil,
		dbRedis:    dbRedis,
	}
}

func isTrue() bool {
	return rand.Intn(100) < 80 // 80% true, 20% false
}

// get user
func (repo *Repository) GetUser(ctx context.Context) ([]entity.User, error) {
	ctx, span := tracer.StartSpan(ctx, "repository.GetUser")
	defer span.End()
	if !isTrue() {
		time.Sleep(10 * time.Millisecond)
		err := fmt.Errorf("dummy error get user")
		return nil, err
	}
	cur, err := repo.dbUser.Find(
		ctx,
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(100),
	)
	if err != nil {
		return nil, mongox.NewError(ctx, err)
	}
	result := make([]entity.User, 0)
	for cur.Next(ctx) {
		var doc model.User
		if err := cur.Decode(&doc); err != nil {
			return nil, mongox.NewError(ctx, err)
		}
		result = append(result, doc.ToEntity())
	}
	return result, nil
}

// create user
func (repo *Repository) CreateUser(ctx context.Context, data entity.User) error {
	ctx, span := tracer.StartSpan(ctx, "repository.CreateUser")
	defer span.End()
	if !isTrue() {
		time.Sleep(10 * time.Millisecond)
		err := fmt.Errorf("dummy error create user")
		return err
	}
	if _, err := repo.dbUser.InsertOne(ctx, model.NewUser(data)); err != nil {
		return mongox.NewError(ctx, err)
	}
	return nil
}

// get message
func (repo *Repository) GetMessage(ctx context.Context) ([]entity.Message, error) {
	ctx, span := tracer.StartSpan(ctx, "repository.GetMessage")
	defer span.End()
	if !isTrue() {
		time.Sleep(10 * time.Millisecond)
		err := fmt.Errorf("dummy error get message")
		return nil, err
	}
	rows, err := repo.dbMessage.WithContext(ctx).
		Model(&model.Message{}).
		Limit(200).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]entity.Message, 0)

	for rows.Next() {
		var msg model.Message
		if err := repo.dbMessage.ScanRows(rows, &msg); err != nil {
			return nil, err
		}
		result = append(result, msg.ToEntity())
	}

	return result, nil
}

// get client do
func (repo *Repository) GetClientDO(ctx context.Context) ([]entity.ClientDo, error) {
	ctx, span := tracer.StartSpan(ctx, "repository.GetClientDO")
	defer span.End()
	// if !isTrue() {
	// 	time.Sleep(10 * time.Millisecond)
	// 	err := errors.NewRetryable(fmt.Errorf("dummy error get client do"))
	// 	tracer.CaptureError(ctx, err)
	// 	return nil, err
	// }
	vals, err := repo.dbRedis.LRange(ctx, "research_apm.client_do", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	result := make([]entity.ClientDo, 0)
	if len(vals) > 0 {
		for _, v := range vals {
			var p entity.ClientDo
			if err := json.Unmarshal([]byte(v), &p); err != nil {
				return nil, err
			}
			result = append(result, p)
		}
		return result, nil
	}
	rows, err := repo.dbClientDO.WithContext(ctx).
		Model(&model.ClientDo{}).
		Limit(200).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listCache := make([]string, 0)

	for rows.Next() {
		var msg model.ClientDo
		if err := repo.dbClientDO.ScanRows(rows, &msg); err != nil {
			return nil, err
		}
		dd := msg.ToEntity()
		result = append(result, dd)
		if v, err := json.Marshal(dd); err == nil {
			listCache = append(listCache, string(v))
		}
	}
	if len(listCache) > 0 {
		if err := repo.dbRedis.RPush(ctx, "research_apm.client_do", listCache).Err(); err != nil {
			return nil, err
		}
		if err := repo.dbRedis.Expire(ctx, "research_apm.client_do", 10*time.Second).Err(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// get profil
func (repo *Repository) GetProfil(ctx context.Context) ([]entity.Profil, error) {
	ctx, span := tracer.StartSpan(ctx, "repository.GetProfil")
	defer span.End()
	if !isTrue() {
		time.Sleep(10 * time.Millisecond)
		err := fmt.Errorf("dummy error get profil")
		return nil, err
	}
	rows, err := repo.dbProfil.WithContext(ctx).
		Model(&model.Profil{}).
		Limit(200).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]entity.Profil, 0)

	for rows.Next() {
		var msg model.Profil
		if err := repo.dbProfil.ScanRows(rows, &msg); err != nil {
			return nil, err
		}
		result = append(result, msg.ToEntity())
	}

	return result, nil
}

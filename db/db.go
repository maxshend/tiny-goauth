package db

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maxshend/tiny_goauth/models"
)

var ctx = context.Background()

// DataLayer is the interface that wraps methods to access database
type DataLayer interface {
	CreateUser(*models.User) error
	UserExistsWithField(fl validator.FieldLevel) (bool, error)
	UserByEmail(string) (*models.User, error)
	StoreCache(key string, payload interface{}, exp time.Duration) error
	DeleteCache(key string) (int64, error)
	GetCacheValue(key string) (string, error)
	DeleteUser(id int64) error
	GetRoles() ([]string, error)
	CreateRoles(names []string) error
	DeleteRoles(names []string) error
	Close()
	Migrate() error
}

type datastore struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

type dbErr string

func (e dbErr) Error() string { return string(e) }

const zeroDeleteRows = dbErr("No row found to delete")
const zeroInsertedRows = dbErr("No rows have been inserted")

// Init initializes connection to the database
func Init() (DataLayer, error) {
	db, err := pgxpool.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &datastore{db: db, rdb: rdb}, nil
}

// Close closes connections to the datastores
func (s *datastore) Close() {
	s.db.Close()
	defer s.rdb.Close()
}

// Migrate runs database migrations
func (s *datastore) Migrate() error {
	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return nil
	}
	sort.Strings(files)
	var sqlScript strings.Builder

	for _, file := range files {
		c, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		sqlScript.WriteString(string(c))
	}

	_, err = s.db.Exec(ctx, sqlScript.String())
	if err != nil {
		return err
	}

	return nil
}

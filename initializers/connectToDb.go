package initializers

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connection to database")
	}
}

const (
	maxConns            = 50
	healthCheckedPeriod = 3 * time.Minute
	maxConnIdleTime     = 1 * time.Minute
	maxConnLifeTime     = 3 * time.Minute
	minConns            = 10
	lazyConnect         = false
)

func ConnectDb() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DB")
	//dataSource := fmt.Sprintf("host=localhost user=postgres password=229847 dbname=ozinshe port=5432 sslmode=disable")
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = maxConns
	poolCfg.HealthCheckPeriod = healthCheckedPeriod
	poolCfg.MaxConnIdleTime = maxConnIdleTime
	poolCfg.MaxConnLifetime = maxConnLifeTime
	poolCfg.MinConns = minConns
	poolCfg.LazyConnect = lazyConnect

	connPool, err := pgxpool.ConnectConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.ConnectConfig")
	}
	return connPool, nil

}

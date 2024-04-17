package data

import (
	"database/sql"
	"fmt"
	"order-service/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewOrdersRepository)

type Data struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db, dbErr := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{})
	if dbErr != nil {
		errMsg := fmt.Sprintf("Error opening database connection: %s", dbErr)
		log.Errorf(errMsg)
		return nil, nil, dbErr
	}

	sqlDB, sqlDBErr := db.DB()
	if sqlDBErr != nil {
		errMsg := fmt.Sprintf("Error getting database connection: %s", sqlDBErr)
		log.Errorf(errMsg)
		return nil, nil, sqlDBErr
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SrtConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	log.Infof("Connected to database successfully!")

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		sqlDB.Close()
	}
	return &Data{db, sqlDB}, cleanup, nil
}

// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/apigee/registry/server/registry/internal/storage/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A complete list of entities used by the storage system
// and managed using gorm features.
var entities = []interface{}{
	&models.Project{},
	&models.Api{},
	&models.Version{},
	&models.Spec{},
	&models.SpecRevisionTag{},
	&models.Deployment{},
	&models.DeploymentRevisionTag{},
	&models.Artifact{},
	&models.Blob{},
}

// Client represents a connection to a storage provider.
type Client struct {
	db *gorm.DB
}

// NewClient creates a new database session using the provided driver and data source name.
// Driver must be one of [ sqlite3, postgres, cloudsqlpostgres ]. DSN format varies per database driver.
//
// PostgreSQL DSN Reference: See "Connection Strings" at https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
// SQLite DSN Reference: See "URI filename examples" at https://www.sqlite.org/c3ref/open.html
func NewClient(ctx context.Context, driver, dsn string) (*Client, error) {
	switch driver {
	case "sqlite3":
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger:      NewGormLogger(ctx),
			PrepareStmt: true,
		})
		if err != nil {
			c := &Client{db: db}
			c.close()
			return nil, grpcErrorForDBError(ctx, err)
		}
		if err := applyConnectionLimits(db); err != nil {
			c := &Client{db: db}
			c.close()
			return nil, grpcErrorForDBError(ctx, err)
		}
		return &Client{db: db}, nil
	case "postgres", "cloudsqlpostgres":
		db, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: driver,
			DSN:        dsn,
		}), &gorm.Config{
			Logger:      NewGormLogger(ctx),
			PrepareStmt: true,
		})
		if err != nil {
			c := &Client{db: db}
			c.close()
			return nil, grpcErrorForDBError(ctx, err)
		}
		if err := applyConnectionLimits(db); err != nil {
			c := &Client{db: db}
			c.close()
			return nil, grpcErrorForDBError(ctx, err)
		}
		return &Client{db: db}, nil
	default:
		return nil, fmt.Errorf("unsupported database %s", driver)
	}
}

// Applies limits to concurrent connections.
// Could be made configurable.
func applyConnectionLimits(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(60 * time.Second)
	return nil
}

// Close closes a database session.
func (c *Client) Close() {
	c.close()
}

func (c *Client) close() {
	sqlDB, _ := c.db.DB()
	sqlDB.Close()
}

func (c *Client) ensureTable(ctx context.Context, v interface{}) error {
	if !c.db.Migrator().HasTable(v) {
		if err := c.db.Migrator().CreateTable(v); err != nil {
			return grpcErrorForDBError(ctx, err)
		}
	}
	return nil
}

// EnsureTables ensures that all necessary tables exist in the database.
func (c *Client) EnsureTables(ctx context.Context) error {
	for _, entity := range entities {
		if err := c.ensureTable(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Migrate(ctx context.Context) error {
	return grpcErrorForDBError(ctx, c.db.WithContext(ctx).AutoMigrate(entities...))
}

func (c *Client) DatabaseName(ctx context.Context) string {
	return c.db.WithContext(ctx).Name()
}

func (c *Client) TableNames(ctx context.Context) ([]string, error) {
	var tableNames []string
	switch c.db.WithContext(ctx).Name() {
	case "postgres":
		if err := c.db.WithContext(ctx).Table("information_schema.tables").Where("table_schema = ?", "public").Order("table_name").Pluck("table_name", &tableNames).Error; err != nil {
			return nil, grpcErrorForDBError(ctx, err)
		}
	case "sqlite":
		if err := c.db.WithContext(ctx).Table("sqlite_schema").Where("type = 'table' AND name NOT LIKE 'sqlite_%'").Order("name").Pluck("name", &tableNames).Error; err != nil {
			return nil, grpcErrorForDBError(ctx, err)
		}
	default:
		return nil, status.Errorf(codes.Internal, "unsupported database %s", c.db.Name())
	}
	return tableNames, nil
}

func (c *Client) RowCount(ctx context.Context, tableName string) (int64, error) {
	var count int64
	err := c.db.WithContext(ctx).Table(tableName).Count(&count).Error
	return count, grpcErrorForDBError(ctx, err)
}

func (c *Client) Transaction(ctx context.Context, fn func(context.Context, *Client) error) error {
	err := c.db.Transaction(func(tx *gorm.DB) error {
		return fn(ctx, &Client{db: tx})
	})
	return grpcErrorForDBError(ctx, err)
}

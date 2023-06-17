package providers

import (
	"context"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"verifymy-golang-test/models"
)

var dbConn *gorm.DB

func Close() {
	dbConn = nil
}

func NewDBConnection(dialector gorm.Dialector) (*gorm.DB, error) {
	if dbConn == nil {
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			return nil, err
		}

		if err := db.AutoMigrate(&models.User{}); err != nil {
			return nil, err
		}

		dbConn = db
	}

	return dbConn, nil
}

func NewDBDialector(ctx context.Context) gorm.Dialector {
	driver, connString := getDriverAndConnString(ctx)

	if driver == "mysql" {
		return mysql.New(mysql.Config{
			DriverName: "mysql",
			DSN:        connString,
		})
	}

	return sqlite.Open(connString)
}

func getDriverAndConnString(ctx context.Context) (string, string) {
	env := os.Getenv("ENV")
	driver := getDriverByEnv(env)
	connString := getConnStringByDriverAndEnv(ctx, env, driver)

	return driver, connString
}

func getDriverByEnv(env string) string {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		if env != "test" {
			driver = "mysql"
		} else {
			driver = "sqlite"
		}
	}

	return driver
}

func getConnStringByDriverAndEnv(
	ctx context.Context, env string, driver string,
) string {
	connString := os.Getenv("DB_CONN_STRING")
	if env == "test" {
		connString = "db.sqlite3"
	}

	return connString
}

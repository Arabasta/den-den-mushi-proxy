package mysql

import (
	"den-den-mushi-Go/pkg/config"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Client(cfg *config.SqlDb, sslCfg *config.Ssl, log *zap.Logger) (*gorm.DB, error) {
	log.Info("Connecting to SQL database...")
	log.Info("Connection parameters",
		zap.String("User", cfg.User),
		zap.String("Host", cfg.Host),
		zap.Int("Port", cfg.Port),
		zap.String("DBName", cfg.DBName),
		zap.String("Params", cfg.Params),
	)

	if cfg.SSLEnabled {
		log.Info("SSL enabled, using CA file", zap.String("CAFile", sslCfg.CAFile))
		if err := registerMySQLTLSCA(sslCfg.CAFile); err != nil {
			return nil, err
		}
		cfg.Params += "&tls=custom"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Params,
	)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// connection pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute)

	// ping
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	log.Info("Connected to SQL database")
	return db, nil
}

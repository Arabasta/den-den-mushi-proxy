package main

import (
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/server"
	"den-den-mushi-Go/internal/control/testdata"
	"den-den-mushi-Go/internal/proxy/websocket_jwt/jti"
	"den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/dto/cyberark"
	"den-den-mushi-Go/pkg/dto/host"
	"den-den-mushi-Go/pkg/dto/iexpress"
	"den-den-mushi-Go/pkg/dto/implementor_groups"
	"den-den-mushi-Go/pkg/dto/os_adm_users"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/dto/proxy_lb"
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/dto/regex_filters"
	"den-den-mushi-Go/pkg/logger"
	"den-den-mushi-Go/pkg/mysql"
	"flag"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	_ = godotenv.Load(".env")
	cfg := config.Load(configPath())

	log := logger.Init(cfg.Logger, cfg.App)
	if log == nil {
		panic("failed to initialize logger")
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Starting Den Den Mushi Control", zap.String("version", cfg.App.Version))

	var db *gorm.DB
	var err error

	if !cfg.Development.IsUsingInvDb {
		db, err = mysql.Client(cfg.DdmDB, cfg.Ssl, log)
	} else {
		db, err = mysql.Client(cfg.InvDB, cfg.Ssl, log)
	}
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if !cfg.Development.IsSMX {
		log.Info("Running AutoMigrate for non-production environment")
		if err := db.AutoMigrate(
			&host.Model{},
			&proxy_lb.Model{},
			&proxy_host.Model{},
			&change_request.Model{},
			&implementor_groups.Model{},
			&cyberark.Model{},
			&regex_filters.Model{},
			&pty_sessions.Model{},
			&connections.Model{},
			&proxy_host.Model{},
			&jti.Model{},
			&os_adm_users.Model{},
			&iexpress.Model{},
		); err != nil {
			log.Fatal("AutoMigrate failed", zap.Error(err))
		}
		if cfg.Development.IsAutoMigrateEnabled {
			testdata.CallAll(db)
		}
	}

	if cfg.Development.IsSMX {
		log.Info("Running AutoMigrate")
		if err := db.AutoMigrate(
			&proxy_lb.Model{},
			&proxy_host.Model{},
			&regex_filters.Model{},
			&pty_sessions.Model{},
			&connections.Model{},
			&jti.Model{},
			&os_adm_users.Model{},
		); err != nil {
			log.Fatal("AutoMigrate failed", zap.Error(err))
		}

		if cfg.Development.IsAutoMigrateEnabled {
			//testdata.CreateProxyHostAndLb(db, cfg)
		}
	}
	reapChildren()
	s := server.New(db, root.Files, cfg, log)
	if err := server.Start(s, cfg, log); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}
}

// configPath usage: go run main.go -config /path/to/config.json
func configPath() string {
	configPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	var finalPath string
	if *configPath != "" {
		finalPath = *configPath
	} else {
		// default to config.json next to executable
		exe, _ := os.Executable()
		exeDir := filepath.Dir(exe)
		finalPath = filepath.Join(exeDir, "config.json")
	}

	return finalPath
}

func reapChildren() {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGCHLD)
		for range sig {
			for {
				var status syscall.WaitStatus
				pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
				if pid <= 0 || err != nil {
					break
				}
			}
		}
	}()
}

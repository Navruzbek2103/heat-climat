package app

import (
	"fmt"
	"log"

	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"gitlab.com/climate.uz/api"

	"gitlab.com/climate.uz/internal/controller/storage"
	"gitlab.com/climate.uz/config"
	"gitlab.com/climate.uz/pkg/db"
	"gitlab.com/climate.uz/pkg/logger"

	"github.com/casbin/casbin/v2"
)

func Run(cfg *config.Config) {
	var (
		casbinEnforcer *casbin.Enforcer
	)
	l := logger.New(cfg.LogLevel)

	casbinEnforcer, err := casbin.NewEnforcer(cfg.AuthConfigPath, cfg.CsvFilePath)
	if err != nil {
		l.Error("casbin enforcer error", err)
		return
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		l.Error("casbin error load policy ", err)
		return
	}
	// postgres://user:password@host:5432/database
	pgxURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PosgresDatabase,
	)
	pg, err := db.New(pgxURL, db.MaxPoolSize(cfg.PGXPoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app run postgres.New %w", err))
	}
	defer pg.Close()

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddDomainMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddDomainMatchingFunc("keyMatch3", util.KeyMatch3)

	strg := storage.NewStoragePg(pg)

	server := api.New(api.Options{
		Conf:           *cfg,
		Logger:         *l,
		Storage:        strg,
		CasbinEnforcer: casbinEnforcer,
	})
	if err := server.Run(cfg.HttpPort); err != nil {
		log.Fatal("failed to run http server", err)
		panic(err)
	}
	log.Print("Server stopped")
}

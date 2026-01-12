package main

import (
	"log"
	"os"

	"github.com/gostructure/app/internal/adapter/storage"
	"github.com/gostructure/app/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down|version]")
	}

	command := os.Args[1]

	// Load app config for timezone
	appCfg, err := config.Load()
	if err != nil {
		log.Fatalf("Load app config failed: %v", err)
	}

	// Load database config
	dbCfg, err := config.LoadDatabaseOnly()
	if err != nil {
		log.Fatalf("Load database config failed: %v", err)
	}

	// Connect using database factory (dynamic based on DB_DRIVER)
	database, err := storage.NewDatabase(dbCfg, appCfg.App.Timezone)
	if err != nil {
		log.Fatalf("Connect database failed: %v", err)
	}
	defer database.Close()

	// Create migrator (supports multiple drivers)
	migrator, err := storage.NewMigrator(database.GetDB(), dbCfg)
	if err != nil {
		log.Fatalf("Create migrator failed: %v", err)
	}
	defer migrator.Close()

	switch command {
	case "up":
		err = migrator.Up()
	case "down":
		err = migrator.Down()
	case "version":
		v, dirty, e := migrator.Version()
		if e != nil {
			log.Fatal(e)
		}
		log.Printf("version=%d dirty=%v", v, dirty)
		return
	default:
		log.Fatalf("Unknown command: %s", command)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed successfully")
}

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal().Msg("DATABASE_URL manquant")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connexion DB")
	}
	defer db.Close()

	// Créer table de suivi migrations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		filename VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		log.Fatal().Err(err).Msg("create migrations table")
	}

	// Lire les fichiers de migration
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal().Err(err).Msg("glob migrations")
	}
	sort.Strings(files)

	for _, file := range files {
		name := filepath.Base(file)

		// Vérifier si déjà appliquée
		var exists bool
		db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename=$1)", name).Scan(&exists)
		if exists {
			fmt.Printf("  ⏭  %s (déjà appliquée)\n", name)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatal().Err(err).Str("file", name).Msg("read migration")
		}

		// Ignorer les commentaires TimescaleDB si extension non disponible
		sql := string(content)
		if !strings.Contains(sql, "create_hypertable") {
			_, err = db.Exec(sql)
		} else {
			// Exécuter statement par statement pour gérer les erreurs TimescaleDB
			for _, stmt := range strings.Split(sql, ";") {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" { continue }
				_, execErr := db.Exec(stmt)
				if execErr != nil && !strings.Contains(execErr.Error(), "already") {
					fmt.Printf("  ⚠  %s: %v\n", name, execErr)
				}
			}
			err = nil
		}

		if err != nil {
			log.Fatal().Err(err).Str("file", name).Msg("apply migration")
		}

		db.Exec("INSERT INTO schema_migrations (filename) VALUES ($1)", name)
		fmt.Printf("  ✅ %s\n", name)
	}

	fmt.Println("Migrations terminées.")
}

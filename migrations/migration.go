package migrations

import "database/sql"

type Migration struct {
	Connector *sql.DB
}

func MigrateContainer(connector *sql.DB) {
	migration := &Migration{
		Connector: connector,
	}
	migration.TableUsersMigrate()
}

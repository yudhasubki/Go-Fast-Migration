package database

import "database/sql"

type Connection struct {
	connector *sql.DB
}

type Fields func(fields []map[string]string)

func (c *Connection) New(connector *sql.DB) *Connection {
	return &Connection{
		connector,
	}
}

func (c *Connection) Column(migrator Fields) {

}

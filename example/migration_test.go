package example

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	migration "github.com/yudhasubki/go-fastmigration"
	schema "github.com/yudhasubki/go-fastmigration"
)

type Migration struct {
	Connector *sql.DB
}

func TestMigration(t *testing.T) {
	connector, err := sql.Open("mysql", "root:test123@tcp(127.0.0.1:3306)/test_database")
	if err != nil {
		log.Fatalln(err.Error())
	}
	migration := &Migration{
		Connector: connector,
	}

	migration.GenderMigration()
	migration.UserMigration()
}

func (m *Migration) GenderMigration() {
	id := schema.Create().Column("id").Type("INT").Nullable(false).Length(11).PrimaryKey().AutoIncrement()
	gender := schema.Create().Column("gender").Type("enum").Enum([]string{"Men", "Women"}).NullableEnum()
	table := schema.Table{
		Columns: []schema.Schema{*id, *gender},
	}
	err := schema.Blueprint(m.Connector, "genders", table)
	if err != nil {
		return
	}
	return
}

func (m *Migration) UserMigration() (err error) {
	id := migration.Create().Column("id").Type("INT").Nullable(false).Length(11).PrimaryKey().AutoIncrement()
	name := migration.Create().Column("name").Type("VARCHAR").Nullable(true).Length(75)
	gender := migration.Create().Column("gender").Type("INT").Length(11)
	created_at := migration.Create().Column("created_at").Type("TIMESTAMP").DefaultTimestamp()
	updated_at := migration.Create().Column("updated_at").Type("TIMESTAMP").NullableTimestamp()
	constraints, err := migration.Add().ForeignKey("gender").References("id").On("genders")

	table := migration.Table{
		Columns:    []migration.Schema{*id, *name, *gender, *created_at, *updated_at},
		Constraint: *constraints,
	}
	err = migration.Blueprint(m.Connector, "users", table)
	if err != nil {
		return
	}
	return
}

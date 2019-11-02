# Go Fast Migration

*Simple database migration for Go (inspired by [laravel](https://github.com/laravel/laravel) migration). This migration powered by golang testing library*

## Installation

```bash
go get github.com/yudhasubki/go-fastmigration
```

## Usage

```go
package example

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
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
		Name:      "genders",
		Connector: m.Connector,
		Columns:   []schema.Schema{*id, *gender},
	}
	err := table.Blueprint()
	if err != nil {
		return
	}
	return
}

func (m *Migration) UserMigration() (err error) {
	id := schema.Create().Column("id").Type("INT").Nullable(false).Length(11).PrimaryKey().AutoIncrement()
	name := schema.Create().Column("name").Type("VARCHAR").Nullable(true).Length(75)
	gender := schema.Create().Column("gender").Type("INT").Length(11)
	created_at := schema.Create().Column("created_at").Type("TIMESTAMP").DefaultTimestamp()
	updated_at := schema.Create().Column("updated_at").Type("TIMESTAMP").NullableTimestamp()
	constraints, err := schema.Add().ForeignKey("gender").References("id").On("genders")

	table := schema.Table{
		Name:       "users",
		Columns:    []schema.Schema{*id, *name, *gender, *created_at, *updated_at},
		Constraint: *constraints,
		Connector:  m.Connector,
	}
	err = table.Blueprint()
	if err != nil {
		return
	}
	return
}

```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

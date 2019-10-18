package migrations

import (
	"github.com/yudhasubki/go-fastmigration/database"
)

func (m *Migration) TableUsersMigrate() (err error) {
	id := database.Create().Column("id").Type("INT").Nullable(false).Length(11).PrimaryKey().AutoIncrement()
	name := database.Create().Column("name").Type("VARCHAR").Nullable(true).Length(75)
	gender := database.Create().Column("gender").Type("enum").Enum([]string{"Men", "Women"}).NullableEnum()
	created_at := database.Create().Column("created_at").Type("TIMESTAMP").DefaultTimestamp()
	updated_at := database.Create().Column("updated_at").Type("TIMESTAMP").NullableTimestamp()
	err = database.Blueprint(m.Connector, "users", id, name, gender, created_at, updated_at)
	if err != nil {
		return
	}
	return
}

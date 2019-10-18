package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Schema struct {
	column        string
	dataType      string
	enum          string
	length        int
	null          string
	autoIncrement string
	primaryKey    string
	defaultColumn string
}

const (
	Default = "DEFAULT"
	Engine  = "InnoDB"
	Charset = "utf8"
)

func Blueprint(connector *sql.DB, tableName string, schemas ...*Schema) (err error) {
	var fields []string
	createTable := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (` + "\n"
	for _, schema := range schemas {
		v := reflect.Indirect(reflect.ValueOf(schema))
		typeOf := v.Type()

		columns := make([]string, 0)
		for i := 0; i < v.NumField(); i++ {
			if typeOf.Field(i).Name == "column" && v.Field(i).String() == "" {
				err = errors.New("Please fill a column name")
				return
			}

			if typeOf.Field(i).Name == "dataType" && v.Field(i).String() == "" {
				err = errors.New("Please fill a data type")
				return
			}
			switch v.Field(i).Kind() {
			case reflect.String:
				var column string
				column = v.Field(i).String()

				if v.Field(i).String() == "" {
					continue
				}

				if typeOf.Field(i).Name == "column" {
					column = fmt.Sprintf("`%s`", v.Field(i).String())
				}
				columns = append(columns, column)
				continue
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v.Field(i).Int() > 0 {
					columns = append(columns, fmt.Sprintf("(%s)", strconv.Itoa(int(v.Field(i).Int()))))
				}
				continue
			}
		}
		fields = append(fields, strings.Join(columns, " "))
	}

	var data string
	for idx, field := range fields {
		if len(fields)-1 == idx {
			data += field + "\n"
			continue
		}
		data += field + ", \n"
	}

	createTable += data
	createTable += fmt.Sprintf(`) ENGINE=%s CHARACTER SET=%s`, Engine, Charset)
	_, err = connector.Exec(createTable)

	if err != nil {
		log.Fatalln(fmt.Sprintf("%v", err))
		return
	}

	return
}

func Create() *Schema {
	return &Schema{}
}

func (s *Schema) Column(columnName string) *Schema {
	if columnName != "" {
		s.column = columnName
	}

	return s
}

func (s *Schema) Type(dataType string) *Schema {
	if dataType != "" {
		s.dataType = dataType
	}
	return s
}

func (s *Schema) Nullable(isNull bool) *Schema {
	s.null = "NOT NULL"

	if isNull {
		s.null = "NULL"
	}

	return s
}

func (s *Schema) PrimaryKey() *Schema {
	if s.primaryKey == "" {
		s.primaryKey = "PRIMARY KEY"
	}

	return s
}

func (s *Schema) AutoIncrement() *Schema {
	s.autoIncrement = "AUTO_INCREMENT"
	return s
}

func (s *Schema) Length(length int) *Schema {
	if length > 0 {
		s.length = length
	}

	return s
}

func (s *Schema) Enum(value []string) *Schema {
	if len(value) > 0 {
		var enum []string
		for _, val := range value {
			enum = append(enum, fmt.Sprintf("'%s'", val))
		}
		s.enum = fmt.Sprintf("(%s)", strings.Join(enum, ","))
	}
	return s
}

func (s *Schema) Default(value string) *Schema {
	if value != "" {
		s.defaultColumn = fmt.Sprintf("DEFAULT '%s'", value)
	}

	return s
}

func (s *Schema) DefaultTimestamp() *Schema {
	s.defaultColumn = "DEFAULT CURRENT_TIMESTAMP"
	return s
}

func (s *Schema) NullableTimestamp() *Schema {
	s.defaultColumn = "NULL DEFAULT NULL"
	return s
}

func (s *Schema) NullableEnum() *Schema {
	s.null = "DEFAULT NULL"
	return s
}

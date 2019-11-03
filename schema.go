package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Table struct {
	Name       string
	Columns    []Schema
	Constraint Constraint
	Engine     string
	Charset    string
	Connector  *sql.DB
}

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

type Constraint struct {
	foreignKeys []string
	references  []string
	constraints string
}

const (
	Default = "DEFAULT"
)

var Engine = struct {
	InnoDB,
	MyISAM,
	Memory,
	CSV,
	Merge,
	Archive,
	Blackhole,
	Federated string
}{
	"InnoDB",
	"MyISAM",
	"Memory",
	"CSV",
	"Merge",
	"Archive",
	"Blackhole",
	"Federated",
}

var Charset = struct {
	Big5,
	Dec8,
	Cp850,
	Hp8,
	Koi8r,
	Latin1,
	Latin2,
	Swe7,
	Ascii,
	Ujis,
	Sjis,
	Hebrew,
	Tis620,
	Euckr,
	Koi82,
	Gb2312,
	Greek,
	Cp1250,
	Gbk,
	Latin5,
	Armscii8,
	Utf8,
	Ucs2,
	Cp866,
	Keybcs2,
	Macce,
	Macroman,
	Cp852,
	Latin7,
	Utf8mb4,
	Cp1251,
	Utf16,
	Cp1256,
	Cp1257,
	Utf32,
	Binary,
	Geostd8,
	Cp932,
	eucjpms string
}{
	"big5",
	"dec8",
	"cp850",
	"hp8",
	"koi8r",
	"latin1",
	"latin2",
	"swe7",
	"ascii",
	"ujis",
	"sjis",
	"hebrew",
	"tis620",
	"euckr",
	"koi8u",
	"gb2312",
	"greek",
	"cp1250",
	"gbk",
	"latin5",
	"armscii8",
	"utf8",
	"ucs2",
	"cp866",
	"keybcs2",
	"macce",
	"macroman",
	"cp852",
	"latin7",
	"utf8mb4",
	"cp1251",
	"utf16",
	"cp1256",
	"cp1257",
	"utf32",
	"binary",
	"geostd8",
	"cp932",
	"eucjpms",
}

var ErrorMessage = struct {
	ConstraintNotMatch string
	ColumnNameEmpty    string
	DataTypeEmpty      string
}{
	"Error: Column Name is Empty",
	"Error: Constraint length not match",
	"Error: Data Type is Empty",
}

func (t *Table) Blueprint() (err error) {
	var fields []string
	createTable := `CREATE TABLE IF NOT EXISTS ` + t.Name + ` (` + "\n"
	for _, schema := range t.Columns {
		v := reflect.Indirect(reflect.ValueOf(schema))
		typeOf := v.Type()

		columns := make([]string, 0)
		for i := 0; i < v.NumField(); i++ {
			if typeOf.Field(i).Name == "column" && v.Field(i).String() == "" {
				err = makeError(ErrorMessage.ColumnNameEmpty)
				return
			}

			if typeOf.Field(i).Name == "dataType" && v.Field(i).String() == "" {
				err = makeError(ErrorMessage.DataTypeEmpty)
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
			if t.Constraint.constraints != "" {
				data += field + ",\n"
				continue
			}
			data += field + "\n"
			continue
		}
		data += field + ", \n"
	}

	createTable += data

	if t.Constraint.constraints != "" {
		createTable += t.Constraint.constraints
	}

	engine := t.Engine
	if engine == "" {
		engine = Engine.InnoDB
	}

	charset := t.Charset
	if charset == "" {
		charset = Charset.Utf8
	}

	createTable += fmt.Sprintf(`) ENGINE=%s CHARACTER SET=%s;`, engine, charset)
	_, err = t.Connector.Exec(createTable)

	if err != nil {
		log.Fatalln(fmt.Sprintf("%v", err))
		return
	}

	return
}

func Create() *Schema {
	return &Schema{}
}

func Add() *Constraint {
	return &Constraint{}
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

func (s *Schema) DefaultCurrentTimestamp() *Schema {
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

func (c *Constraint) ForeignKey(foreignKeyName ...string) *Constraint {
	fmt.Println(foreignKeyName)
	c.foreignKeys = foreignKeyName
	return c
}

func (c *Constraint) References(references ...string) *Constraint {
	fmt.Println(references)
	c.references = references
	return c
}

/*
	On - parameters is filled a table name
*/
func (c *Constraint) On(tables ...string) (*Constraint, error) {
	if len(c.foreignKeys) != len(c.references) {
		return c, makeError(ErrorMessage.ConstraintNotMatch)
	}

	var wrappingUpFk []string
	for idx, fk := range c.foreignKeys {
		foreignKey := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", fk, tables[idx], string(c.references[idx]))
		wrappingUpFk = append(wrappingUpFk, foreignKey)
	}

	var foreignKeys string
	for idx, fk := range wrappingUpFk {
		if idx == len(wrappingUpFk)-1 {
			foreignKeys += fmt.Sprintf("%s \n", fk)
			break
		}
		foreignKeys += fmt.Sprintf("%s, \n", fk)
	}

	c.constraints = foreignKeys
	return c, nil
}

func makeError(errorMessage string) error {
	return errors.New(errorMessage)
}

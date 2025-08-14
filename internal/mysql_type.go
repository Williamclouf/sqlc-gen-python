package python

import (
	"log"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func mysqlType(req *plugin.GenerateRequest, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)

	switch columnType {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		return "int"
	case "float", "double", "real":
		return "float"
	case "decimal", "numeric":
		return "decimal.Decimal"
	case "bit", "boolean", "bool":
		return "bool"
	case "json":
		return "Any"
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		return "memoryview"
	case "date":
		return "datetime.date"
	case "time":
		return "datetime.time"
	case "datetime", "timestamp":
		return "datetime.datetime"
	case "char", "varchar", "text", "tinytext", "mediumtext", "longtext":
		return "str"
	case "enum", "set":
		return "str"
	default:
		for _, schema := range req.Catalog.Schemas {
			if schema.Name == "information_schema" || schema.Name == "mysql" {
				continue
			}
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if schema.Name == req.Catalog.DefaultSchema {
						return "models." + modelName(enum.Name, req.Settings)
					}
					return "models." + modelName(schema.Name+"_"+enum.Name, req.Settings)
				}
			}
		}
		log.Printf("unknown MySQL type: %s\n", columnType)
		return "Any"
	}
}

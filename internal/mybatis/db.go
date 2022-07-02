package mybatis

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// GetTableColumns 查询所有字段
func GetTableColumns(tableName string) (columns []map[string]string) {
	var (
		queryColumns []string
	)
	rows, _ := DB.Query(fmt.Sprintf("select column_name, column_comment, data_type, column_key, extra from information_schema.columns where table_schema='%s' and table_name= '%s'", Database, tableName))
	queryColumns, _ = rows.Columns()
	values := make([]sql.RawBytes, len(queryColumns))
	scans := make([]interface{}, len(queryColumns))

	for i := range values {
		scans[i] = &values[i]
	}

	for rows.Next() {
		_ = rows.Scan(scans...)
		each := make(map[string]string)
		for i, col := range values {
			each[queryColumns[i]] = string(col)
		}
		columns = append(columns, each)
	}
	return
}

// GetTables 查询表
func GetTables() (tables []map[string]string) {
	var colSql strings.Builder
	colSql.WriteString(fmt.Sprintf("select table_name, table_comment from information_schema.tables where table_schema='%s'", Database))
	if !IsAllTables {
		tableNams := strings.Split(Table, ",")
		var inCaluse strings.Builder
		for i, v := range tableNams {
			inCaluse.WriteString(fmt.Sprintf("'%s'", v))
			if i != len(tableNams)-1 {
				inCaluse.WriteString(",")
			}
		}
		colSql.WriteString(fmt.Sprintf(" and table_name in (%s)", inCaluse.String()))
	}

	rows, _ := DB.Query(colSql.String())
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scans := make([]interface{}, len(columns))

	for i := range values {
		scans[i] = &values[i]
	}

	// 所有表
	for rows.Next() {
		_ = rows.Scan(scans...)
		each := make(map[string]string)
		for i, col := range values {
			each[columns[i]] = string(col)
		}
		tables = append(tables, each)
	}
	return
}

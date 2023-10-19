package dbhelper

import (
	"helper_go/comhelper"
	"html/template"
	"strings"
)

/*
select * from tables where table_schema = 'dbnote'

SELECT

	COLUMN_NAME,DATA_TYPE, COLUMN_COMMENT,
	COLUMN_DEFAULT,COLUMN_KEY,Comment
	FROM COLUMNS
	WHERE TABLE_NAME = 'mail'  and TABLE_SCHEMA = 'dbnote'
*/
type ModelInfo struct {
	BDName          string
	TablePrefixName string
	TableName       string
	PackageName     string
	ModelName       string
	TableSchema     *[]TABLE_SCHEMA
}

type TABLE_SCHEMA struct {
	Field   string `db:"Field" json:"Field"`
	Type    string `db:"Type" json:"Type"`
	Key     string `db:"Key" json:"Key"`
	Comment string `db:"Comment" json:"Comment"`
}

func (m *ModelInfo) ColumnNames() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {

		result = append(result, t.Field)

	}
	return result
}

func (m *ModelInfo) ColumnCount() int {
	return len(*m.TableSchema)
}

func (m *ModelInfo) PkColumnsSchema() []TABLE_SCHEMA {
	result := make([]TABLE_SCHEMA, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.Key == "PRI" {
			result = append(result, t)
		}
	}
	return result
}

func (m *ModelInfo) HavePk() bool {
	return len(m.PkColumnsSchema()) > 0
}

func (m *ModelInfo) NoPkColumnsSchema() []TABLE_SCHEMA {
	result := make([]TABLE_SCHEMA, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.Key != "PRI" {
			result = append(result, t)
		}
	}
	return result
}

func (m *ModelInfo) NoPkColumns() []string {
	noPkColumnsSchema := m.NoPkColumnsSchema()
	result := make([]string, 0, len(noPkColumnsSchema))
	for _, t := range noPkColumnsSchema {
		result = append(result, t.Field)
	}
	return result
}

func (m *ModelInfo) PkColumns() []string {
	pkColumnsSchema := m.PkColumnsSchema()
	result := make([]string, 0, len(pkColumnsSchema))
	for _, t := range pkColumnsSchema {
		result = append(result, t.Field)
	}
	return result
}

func IsUUID(str string) bool {
	return "uuid" == str
}

func FirstCharLower(str string) string {
	if len(str) > 0 {
		return strings.ToLower(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

func FirstCharUpper(str string) string {
	if len(str) > 0 {
		return strings.ToUpper(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

func Tags(columnName string) template.HTML {

	return template.HTML("`db:" + `"` + columnName + `"` +
		" json:" + `"` + columnName + "\"`")
}

func ExportColumn(columnName string) string {
	columnItems := strings.Split(columnName, "_")
	columnItems[0] = FirstCharUpper(columnItems[0])
	for i := 0; i < len(columnItems); i++ {
		item := strings.Title(columnItems[i])

		//if strings.ToUpper(item) == "ID" {
		//	item = "Id"
		//}

		columnItems[i] = item
	}

	return strings.Join(columnItems, "")

	//return strings.Title(columnName)
}

//func TypeConvert(str string) string {
//
//	switch str {
//	case "smallint", "tinyint":
//		return "int8"
//
//	case "varchar", "text", "longtext", "char":
//		return "string"
//
//	case "date":
//		return "string"
//
//	case "int":
//		return "int"
//
//	case "timestamp", "datetime":
//		return "time.Time"
//
//	case "bigint":
//		return "int64"
//
//	case "float", "double", "decimal":
//		return "float64"
//	case "uuid":
//		return "gocql.UUID"
//
//	default:
//		return str
//	}
//}

func TypeConvert(str string) string {

	sliceInt8 := []string{"smallint", "tinyint"}
	if v := comhelper.InArrayContains(sliceInt8, str); v {
		//return "int8"
		return "int"
	}

	sliceStr := []string{"varchar", "text", "longtext", "char", "date", "enum"}
	if v := comhelper.InArrayContains(sliceStr, str); v {
		return "string"
	}

	sliceDate := []string{"timestamp", "datetime"}
	if v := comhelper.InArrayContains(sliceDate, str); v {
		return "time.Time"
	}

	sliceBig := []string{"bigint"}
	if v := comhelper.InArrayContains(sliceBig, str); v {
		return "int64"
	}

	sliceFlo := []string{"float", "double", "decimal"}
	if v := comhelper.InArrayContains(sliceFlo, str); v {
		return "float64"
	}

	sliceInt := []string{"int"}
	if v := comhelper.InArrayContains(sliceInt, str); v {
		return "int"
	}

	return str
}

func Join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func ColumnAndType(table_schema []TABLE_SCHEMA) string {
	result := make([]string, 0, len(table_schema))
	for _, t := range table_schema {
		result = append(result, t.Field+" "+TypeConvert(t.Type))
	}
	return strings.Join(result, ",")
}

func ColumnWithPostfix(columns []string, Postfix, sep string) string {
	result := make([]string, 0, len(columns))
	for _, t := range columns {
		result = append(result, t+Postfix)
	}
	return strings.Join(result, sep)
}

func MakeQuestionMarkList(num int) string {
	a := strings.Repeat("?,", num)
	return a[:len(a)-1]
}

/*
func joinQuestionMarkByComma(tableSchema *[]TABLE_SCHEMA) string {
	columns := make([]string, 0, len(*tableSchema))
	for _, _ = range *tableSchema {
		columns = append(columns, "?")
	}

	return strings.Join(columns, ",")
}
*/

func (m *ModelInfo) CheckFirstTable() string {
	tableName := LoadIni("DataBase", "table_names")
	dbName := LoadIni("DataBase", "db_name")
	if tableName != "" {
		tableNameSlice := strings.Split(tableName, ",")
		return tableNameSlice[0]
	} else {
		getTablesNameSql := "show tables from " + dbName
		tablaNames, _ := NewEngineInit().QueryString(getTablesNameSql)
		return tablaNames[0]["Tables_in_"+dbName]
	}
}

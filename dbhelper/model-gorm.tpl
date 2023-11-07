{{$exportModelName := .ModelName | ExportColumn}}

package {{.PackageName}}

import (
    "errors"
	"fmt"
	"helper_go/comhelper"
	"helper_go/dbhelper"
	"log"
)

type {{$exportModelName}} struct {
{{range .TableSchema}} {{.Field | ExportColumn}} {{.Type | TypeConvert}} {{.Field | Tags}} // {{.Comment}}
{{end}}}

var Default{{$exportModelName}} = &{{$exportModelName}}{}
{{if eq .CheckFirstTable .ModelName}}
type Arr map[string]interface{}
{{end}}

// Query Executes native sql, returning the defined struct type
func (m *{{$exportModelName}}) Query(args ...interface{}) ([]{{$exportModelName}}, error) {
	// sql
	s := ""
	switch val := args[0].(type) {
	case string:
		s = val
	}
	// Parameter mapping
	var params []interface{}
	if len(args) > 1 {
		switch val := args[1].(type) {
		case []interface{}:
			params = val
		}
	}
	// order
	if len(args) > 2 {
		switch val := args[2].(type) {
		case string:
			s += " " + val
		}
	}
	// limit
	if len(args) > 3 {
		switch val := args[3].(type) {
		case []int:
			s += " limit " + comhelper.IntToString(val[0]) + "," + comhelper.IntToString(val[1])
		}
	}

	var result []{{$exportModelName}}
	err := dbhelper.GormConnection().Raw(s, params...).Find(&result).Error
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " QueryStructure")
		return []{{$exportModelName}}{}, err
	}
	return result, nil
}

// QueryByMap Assemble the SQL query through the parameter, return the defined structure type
func (m *{{$exportModelName}}) QueryByMap(args ...interface{}) ([]{{$exportModelName}}, error) {
	var result []{{$exportModelName}}

	s := "select * from {{.TablePrefixName}} where 1=1 "
	// where
	var params []interface{}
	switch val := args[0].(type) {
	case Arr:
		for k, v := range val {
			switch vType := v.(type) {
			case string:
				if vType == "" {
					continue
				}
			}
			s += fmt.Sprintf("and {{.TablePrefixName}}.%s = ? ", k)
			params = append(params, v)
		}
	}
	// order
	if len(args) > 1 {
		switch val := args[1].(type) {
		case string:
			if val != "" {
				s += val + " "
			}
		}
	}
	// limit
	if len(args) > 2 {
		switch val := args[2].(type) {
		case []int:
			s += "limit " + comhelper.IntToString(val[0]) + "," + comhelper.IntToString(val[1])
		}
	}

	err := dbhelper.GormConnection().Raw(s, params...).Find(&result).Error
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " QueryStructureByMap")
		return []{{$exportModelName}}{}, err
	}
	return result, nil
}

// Count All records
func (m *{{$exportModelName}}) Count() (int64, error) {
	var count int64
	err := dbhelper.GormConnection().Model(m).Count(&count).Error
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " Count")
		return int64(0), err
	}
	return count, nil
}

// Delete Assemble the sql query with parameters
func (m *{{$exportModelName}}) Delete(args Arr) error {
	if len(args) == 0 {
		return errors.New("the where condition cannot be empty")
	}

	s := "delete from {{.TablePrefixName}} where 1=1 "
	var params []interface{}
	for k, v := range args {
		s += fmt.Sprintf("and {{.TablePrefixName}}.%s = ? ", k)
		params = append(params, v)
	}
	err := dbhelper.GormConnection().Raw(s, params...).Row().Err()
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " Delete")
		return err
	}
	return nil
}

// DeleteByStructure By primary key id
func (m *{{$exportModelName}}) DeleteByStructure(id interface{}) error {
	err := dbhelper.GormConnection().Delete(m, id).Error
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " DeleteByStructure")
		return err
	}
	return nil
}

// Update Assemble the sql query with parameters
func (m *{{$exportModelName}}) Update(set, args Arr) error {
	if len(args) == 0 {
		return errors.New("the where condition cannot be empty")
	}

	s := "update {{.TablePrefixName}} set "
	var params []interface{}
	for k, v := range set {
		s += fmt.Sprintf(" {{.TablePrefixName}}.%s = ?,", k)
		params = append(params, v)
	}

	s = s[:len(s)-1]
	s += " where 1=1 "
	for k, v := range args {
		s += fmt.Sprintf("and {{.TablePrefixName}}.%s = ? ", k)
		params = append(params, v)
	}

	err := dbhelper.GormConnection().Raw(s, params...).Row().Err()
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " Update")
		return err
	}
	return nil
}

// UpdateByStructure The primary key id is required
func (m *{{$exportModelName}}) UpdateByStructure(args *{{$exportModelName}}) error {
	err := dbhelper.GormConnection().Model(m).Updates(args).Error
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " UpdateByStructure")
		return err
	}
	return nil
}

// Insert Assemble the sql query with parameters
func (m *{{$exportModelName}}) Insert(args Arr) error {
	s := "insert into {{.TablePrefixName}} "
	sqlFiele := ""
	sqlValues := ""
	var params []interface{}
	for k, v := range args {
		sqlFiele += fmt.Sprintf("%s,", k)
		sqlValues += "?,"
		params = append(params, v)
	}

	sqlFiele = sqlFiele[:len(sqlFiele)-1]
	sqlValues = sqlValues[:len(sqlValues)-1]
	s += "(" + sqlFiele + ") values (" + sqlValues + ")"

	err := dbhelper.GormConnection().Raw(s, params...).Row().Err()
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " Insert")
		return err
	}
	return nil
}

// InsertByStructure /**
func (m *{{$exportModelName}}) InsertByStructure(args ...string) error {
	err := dbhelper.GormConnection().Omit(args...).Create(m).Error
	if err != nil {
		log.Print(err.Error(), "{{$exportModelName}}", " InsertByStructure")
		return err
	}
	return nil
}

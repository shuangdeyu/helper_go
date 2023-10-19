{{$exportModelName := .ModelName | ExportColumn}}

package {{.PackageName}}

import (
	"helper_go/dbhelper"
	"helper_go/comhelper"
	"fmt"
	"log"
)

type {{$exportModelName}} struct {
{{range .TableSchema}} {{.Field | ExportColumn}} {{.Type | TypeConvert}} {{.Field | Tags}} // {{.Comment}}
{{end}}}

var Default{{$exportModelName}} = &{{$exportModelName}}{}
{{if eq .CheckFirstTable .ModelName}}
type Arr map[string]interface{}
{{end}}

/**
 * 执行原生sql查询，返回string类型的map
 */
func (m *{{$exportModelName}}) Query(args ...interface{}) ([]map[string]string, error) {
	// 基础sql语句
	sql := ""
	switch val := args[0].(type) {
	case string:
		sql = val
	}
	// 映射参数
	params := []interface{}{}
	if len(args) > 1 {
		switch val := args[1].(type) {
		case []interface{}:
			params = val
		}
	}
	// order 语句拼接
	if len(args) > 2 {
		switch val := args[2].(type) {
		case string:
			sql += " " + val
		}
	}
	// limit 语句拼接
	if len(args) > 3 {
		switch val := args[3].(type) {
		case []int:
			sql += " limit " + comhelper.IntToString(val[0]) + "," + comhelper.IntToString(val[1])
		}
	}

	ret, err := dbhelper.NewEngineInit().SQL(sql, params...).QueryString()
	if err != nil {
		log.Print(err.Error())
		return []map[string]string{}, err
	}
	return ret, nil
}

/**
 * 执行原生sql，返回定义的结构体类型
 */
func (m *{{$exportModelName}}) QueryStructure(args ...interface{}) ([]{{$exportModelName}}, error) {
	// 基础sql语句
	sql := ""
	switch val := args[0].(type) {
	case string:
		sql = val
	}
	// 映射参数
	params := []interface{}{}
	if len(args) > 1 {
		switch val := args[1].(type) {
		case []interface{}:
			params = val
		}
	}
	// order 语句拼接
	if len(args) > 2 {
		switch val := args[2].(type) {
		case string:
			sql += " " + val
		}
	}
	// limit 语句拼接
	if len(args) > 3 {
		switch val := args[3].(type) {
		case []int:
			sql += " limit " + comhelper.IntToString(val[0]) + "," + comhelper.IntToString(val[1])
		}
	}

	result := []{{$exportModelName}}{}
	err := dbhelper.NewEngineInit().SQL(sql, params...).Find(&result)
	if err != nil {
		log.Print(err.Error())
		return []{{$exportModelName}}{}, err
	}
	return result, nil
}

/**
 * 通过参数构造sql查询，返回string类型的map
 */
func (m *{{$exportModelName}}) QueryByMap(args ...interface{}) ([]map[string]string, error) {
	sql := "select * from {{.TablePrefixName}} where 1=1 "
	// 拼接where语句
	var params []interface{}
	switch val := args[0].(type) {
	case Arr:
		for k, v := range val {
			switch v_type := v.(type) {
			case string:
				if v_type == "" {
					continue
				}
			}
			sql += fmt.Sprintf("and {{.TablePrefixName}}.%s = ? ", k)
			params = append(params, v)
		}
	}
	// 拼接order语句
	if len(args) > 1 {
		switch val := args[1].(type) {
		case string:
			if val != "" {
				sql += val + " "
			}
		}
	}
	// 拼接limit语句
	if len(args) > 2 {
		switch val := args[2].(type) {
		case []int:
			sql += "limit " + comhelper.IntToString(val[0]) + "," + comhelper.IntToString(val[1])
		}
	}

	ret, err := dbhelper.NewEngineInit().SQL(sql, params...).QueryString()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return ret, nil
}

/**
 * 通过参数构造sql查询，返回定义的结构体类型
 */
func (m *{{$exportModelName}}) QueryStructureByMap(args ...interface{}) ([]{{$exportModelName}}, error) {
	result := []{{$exportModelName}}{}

	sql := "select * from {{.TablePrefixName}} where 1=1 "
	// 拼接where语句
	var params []interface{}
	switch val := args[0].(type) {
	case Arr:
		for k, v := range val {
			switch v_type := v.(type) {
			case string:
				if v_type == "" {
					continue
				}
			}
			sql += fmt.Sprintf("and {{.TablePrefixName}}.%s = ? ", k)
			params = append(params, v)
		}
	}
	// 拼接order语句
	if len(args) > 1 {
		switch val := args[1].(type) {
		case string:
			if val != "" {
				sql += val + " "
			}
		}
	}
	// 拼接limit语句
	if len(args) > 2 {
		switch val := args[2].(type) {
		case []int:
			sql += "limit " + comhelper.IntToString(val[0]) + "," + comhelper.IntToString(val[1])
		}
	}

	err := dbhelper.NewEngineInit().SQL(sql, params...).Find(&result)
	if err != nil {
		log.Println(err.Error())
		return []{{$exportModelName}}{}, err
	}
	return result, nil
}

/**
 * 获取count
 */
func (m *{{$exportModelName}}) Count() (int64, error) {
	ret, err := dbhelper.NewEngineInit().Count(m)
	if err != nil {
		log.Println(err)
		return int64(0), err
	}
	return ret, nil
}

/**
 * 删除，通过参数构造sql
 */
func (m *{{$exportModelName}}) Delete(args Arr) error {
	sql := "delete from {{.TablePrefixName}} where 1=1 "
	var params []interface{}
	for k, v := range args {
		sql += fmt.Sprintf("and {{.TablePrefixName}}.%s = ? ", k)
		params = append(params, v)
	}
	_, err := dbhelper.NewEngineInit().SQL(sql, params...).QueryString()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/**
 * 删除，绑定结构体
 */
func (m *{{$exportModelName}}) DeleteByStructure(id int) error {
	_, err := dbhelper.NewEngineInit().Id(id).Delete(m)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/**
 * 更新，通过参数构造sql
 */
func (m *{{$exportModelName}}) Update(set, args Arr) error {
	sql := "update {{.TablePrefixName}} set "
	var params []interface{}
	for k, v := range set {
		sql += fmt.Sprintf("{{.TablePrefixName}}.%s = ?,", k)
		params = append(params, v)
	}

	sql = sql[:len(sql)-1]
	sql += " where 1=1 "
	for k, v := range args {
		sql += fmt.Sprintf("and {{.TablePrefixName}}.%s = ? ", k)
		params = append(params, v)
	}

	_, err := dbhelper.NewEngineInit().SQL(sql, params...).QueryString()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/**
 * 更新，绑定结构体
 */
func (m *{{$exportModelName}}) UpdateByStructure(args *{{$exportModelName}}) error {
	_, err := dbhelper.NewEngineInit().Update(m, args)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/**
 * 新增，绑定结构体
 */
func (m *{{$exportModelName}}) InsertByStructure(args ...string) error {
	_, err := dbhelper.NewEngineInit().Omit(args...).Insert(m)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/**
 * 新增，通过参数构造sql
 */
func (m *{{$exportModelName}}) Insert(args Arr) error {
	sql := "insert into {{.TablePrefixName}} "
	sql_fiele := ""
    sql_values := ""
	var params []interface{}
	for k, v := range args {
    	sql_fiele += fmt.Sprintf("%s,", k)
    	sql_values += "?,"
    	params = append(params, v)
    }

    sql_fiele = sql_fiele[:len(sql_fiele)-1]
    sql_values = sql_values[:len(sql_values)-1]
    sql += "(" + sql_fiele + ") values (" + sql_values + ")"

	_, err := dbhelper.NewEngineInit().SQL(sql, params...).QueryString()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
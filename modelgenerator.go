package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shuangdeyu/helper_go/dbhelper"
)

func genModelFile(render *template.Template, dbName, tableName string) {
	tableSchema := []dbhelper.TABLE_SCHEMA{}
	err := dbhelper.NewEngineInit().SQL(
		"show full columns from " + tableName + " from " + dbName).Find(&tableSchema)

	if err != nil {
		fmt.Println(err)
		return
	}
	if prefix := dbhelper.LoadIni("DataBase", "table_prefix"); prefix != "" {
		tableName = tableName[len(prefix):]
	}
	fileName := dbhelper.LoadIni("File", "model_file") + strings.ToLower(tableName) + ".go"
	os.Remove(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &dbhelper.ModelInfo{
		PackageName:     "model",
		BDName:          dbName,
		TablePrefixName: dbhelper.LoadIni("DataBase", "table_prefix") + tableName,
		TableName:       tableName,
		ModelName:       tableName,
		TableSchema:     &tableSchema,
	}
	if err := render.Execute(f, model); err != nil {
		log.Fatal(err)
	}
	fmt.Println(fileName)
	cmd := exec.Command("goimports", "-w", fileName)
	//cmd := exec.Command("gofmt", "-w", fileName)
	cmd.Run()
}

var confFile = flag.String("cf", "./config.ini", "the path of ini file")

func main() {

	flag.Parse()
	confParam := &dbhelper.ConfInit{
		FilePath: *confFile,
	}
	confParam.FileInit()

	logDir, _ := filepath.Abs(dbhelper.LoadIni("File", "model_file"))
	if _, err := os.Stat(logDir); err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}

	data, err := ioutil.ReadFile(dbhelper.LoadIni("File", "tpl_file"))
	if nil != err {
		fmt.Printf("%v\n", err)
		return
	}

	render := template.Must(template.New("model").
		Funcs(template.FuncMap{
			"FirstCharUpper":       dbhelper.FirstCharUpper,
			"TypeConvert":          dbhelper.TypeConvert,
			"Tags":                 dbhelper.Tags,
			"ExportColumn":         dbhelper.ExportColumn,
			"Join":                 dbhelper.Join,
			"MakeQuestionMarkList": dbhelper.MakeQuestionMarkList,
			"ColumnAndType":        dbhelper.ColumnAndType,
			"ColumnWithPostfix":    dbhelper.ColumnWithPostfix,
		}).Parse(string(data)))

	tableName := dbhelper.LoadIni("DataBase", "table_names")
	tablePrefix := dbhelper.LoadIni("DataBase", "table_prefix")
	dbName := dbhelper.LoadIni("DataBase", "db_name")
	if tableName != "" {
		tableNameSlice := strings.Split(tableName, ",")
		for _, v := range tableNameSlice {
			if tablePrefix != "" {
				v = tablePrefix + v
			}
			genModelFile(render, dbName, v)
		}
	} else {
		getTablesNameSql := "show tables from " + dbName
		tablaNames, err := dbhelper.NewEngineInit().QueryString(getTablesNameSql)
		if err != nil {
			fmt.Println(err)
		}
		for _, table := range tablaNames {
			tableCol := "Tables_in_" + dbName
			tablePrefixName := table[tableCol]
			genModelFile(render, dbName, tablePrefixName)
		}
	}
}

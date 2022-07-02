package mybatis

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

var (
	FieldTypes = map[string]string{
		`bigint`:    "java.lang.Long",
		`bit`:       "java.lang.Boolean",
		`char`:      "java.lang.String",
		`date`:      "java.util.Date",
		`datetime`:  "java.util.Date",
		`decimal`:   "java.math.BigDecimal",
		`double`:    "java.lang.Double",
		`float`:     "java.lang.Float",
		`int`:       "java.lang.Integer",
		`integer`:   "java.lang.Long",
		`text`:      "java.lang.String",
		`longtext`:  "java.lang.String",
		`time`:      "java.util.Date",
		`timestamp`: "java.util.Date",
		`tinyint`:   "java.lang.Integer",
		`varchar`:   "java.lang.String",
	}
	JdbcTypes = map[string]string{
		`char`:     "CHAR",
		`varchar`:  "VARCHAR",
		`tinyint`:  "TINYINT",
		`smallint`: "SMALLINT",
		`int`:      "INTEGER",
		`float`:    "FLOAT",
		`bigint`:   "BIGINT",
		`double`:   "DOUBLE",
		`date`:     "TIMESTAMP",
		`datetime`: "TIMESTAMP",
		`time`:     "TIMESTAMP",
		`text`:     "VARCHAR",
		`longtext`: "LONGVARCHAR",
		`decimal`:  "DECIMAL",
	}
)

type Field struct {
	Annotations []string
	Field       string
	Comment     string
}
type JavaModel struct {
	Annotations  []string
	Package      string
	TableName    string
	Name         string
	Fields       []Field
	Descriptions string
}

type Mapper struct {
	Package      string
	Imports      string
	Descriptions string
	Annotations  string
	Name         string
	Model        string
}

type XmlModel struct {
	Mapper  string
	Model   string
	Results []Result
}
type Result struct {
	Name     string
	Column   string
	JdbcType string
	Property string
}

func initDb() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%s?charset=utf8", UserName, Password, Host, Port, Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}

func Run() {

	if !IsAllTables && Table == "" {
		IsAllTables = true
	}

	initDb()
	tables := GetTables()

	if len(tables) == 0 {
		fmt.Printf("未查询到有效的表[%s] \n", Table)
		return
	}

	for i := 0; i < len(tables); i++ {
		columns := GetTableColumns(tables[i]["table_name"])
		// 生成model
		GeneratorModel(columns, tables[i])
		// 生成mapper
		GeneratorMapper(tables[i])
		// 生成xml
		GeneratorXml(columns, tables[i])
	}

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CheckPath(path string) {
	exist, err := PathExists(path)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	if !exist {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(fmt.Sprintf("mkdir failed![%v]\n", err))
		}
	}
}

// ToJavaName 转换Java名称
func ToJavaName(s string) string {
	arr := strings.Split(s, "_")
	var result string = ""
	for _, str := range arr {
		slen := len(str)
		result = result + strings.ToUpper(str[0:1]) + str[1:slen]
	}
	return result
}

// ToJavaBeanField 转换属性
func ToJavaBeanField(field string, fieldType string) string {
	_fieldType := FieldTypes[fieldType]
	_fieldType = GetTypeName(_fieldType)
	_field := ToHumpField(field)
	return "private " + _fieldType + " " + _field + ";"
}

func Distinct(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

// ToJavaBeanFieldCommennt 字段备注
func ToJavaBeanFieldCommennt(commennt string) string {
	return "/**\n" +
		"	 * " + commennt +
		"\n	 */"
}
func WriteDescriptions(table string, tableComment string) string {
	return "/**\n" + " * @author " + Author + "\n" +
		" * @time " + time.Now().Format("2006-01-02 15:04:05") + "\n" +
		" * @description " + tableComment + "\n" + " */"
}

// GetTypeName 获取类型
func GetTypeName(str string) string {
	arr := strings.Split(str, ".")
	lens := len(arr)
	result := arr[lens-1]
	return result
}

// ToHumpField 转驼峰
func ToHumpField(field string) string {
	arr := strings.Split(field, "_")
	var result string = ""
	for i, str := range arr {
		if i != 0 {
			slen := len(str)
			result = result + strings.ToUpper(str[0:1]) + str[1:slen]
		} else {
			result = result + str
		}
	}
	return result
}

func GeneratorModel(items []map[string]string, table map[string]string) {
	javaName := ToJavaName(table["table_name"])
	tpl, _ := template.New("").Parse(ModelTpl)

	fields := make([]Field, 0)
	classAnnotations := []string{
		"import lombok.AllArgsConstructor;",
		"import lombok.Builder;",
		"import lombok.Data;",
		"import lombok.NoArgsConstructor;",
		"",
	}
	for _, item := range items {
		annotations := make([]string, 0)
		hasPri := strings.Contains(item["column_key"], "PRI")

		if hasPri {
			annotations = append(annotations, "@Id")

		}
		if hasPri && strings.Contains(item["extra"], "auto_increment") {
			annotations = append(annotations, "@GeneratedValue(strategy = GenerationType.IDENTITY)")
			classAnnotations = append(classAnnotations, "import javax.persistence.GeneratedValue;")
			classAnnotations = append(classAnnotations, "import javax.persistence.GenerationType;")
		}

		if hasPri {
			classAnnotations = append(classAnnotations, "import javax.persistence.Id;")
			classAnnotations = append(classAnnotations, "import javax.persistence.Table;")
		}

		classAnnotations = append(classAnnotations, "import java.io.Serializable;")

		// date
		if strings.Contains(FieldTypes[item["data_type"]], "Date") {
			classAnnotations = append(classAnnotations, "import java.util.Date;")
		}

		fields = append(fields, Field{
			Annotations: annotations,
			Field:       ToJavaBeanField(item["column_name"], item["data_type"]),
			Comment:     ToJavaBeanFieldCommennt(item["column_comment"]),
		})
	}

	javaModel := &JavaModel{
		Annotations:  Distinct(classAnnotations),
		Descriptions: WriteDescriptions(table["table_name"], table["table_comment"]),
		Package:      ModelPackage,
		Name:         javaName,
		TableName:    table["table_name"],
		Fields:       fields,
	}

	filePath := OutFileDir + strings.Replace(ModelPackage, ".", "/", -1)
	CheckPath(filePath)
	file, err := os.OpenFile(filePath+"/"+javaName+".java", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}

	tpl.Execute(file, javaModel)

}

func GeneratorMapper(table map[string]string) {
	tpl, _ := template.New("").Parse(MapperTpl)

	var (
		_imports      strings.Builder
		_descriptions strings.Builder
		_annotations  strings.Builder
	)

	_imports.WriteString("\n")

	javaName := ToJavaName(table["table_name"])

	// 导入实体包
	if !strings.Contains(_imports.String(), javaName) {
		_imports.WriteString(fmt.Sprintf("import %s.%s;", ModelPackage, javaName))
	}

	// 类注释
	if !strings.Contains(_descriptions.String(), "@author") {
		_descriptions.WriteString(WriteDescriptions(table["table_name"], table["table_comment"]))
	}

	mapper := &Mapper{
		Package:      MapperPackage,
		Imports:      _imports.String(),
		Descriptions: _descriptions.String(),
		Annotations:  _annotations.String(),
		Name:         javaName + "Mapper",
		Model:        javaName,
	}

	filePath := OutFileDir + strings.Replace(MapperPackage, ".", "/", -1)
	CheckPath(filePath)
	file, err := os.OpenFile(filePath+"/"+javaName+"Mapper.java", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}

	tpl.Execute(file, mapper)
}

func GeneratorXml(items []map[string]string, table map[string]string) {
	tpl, _ := template.New("").Parse(XmlTpl)
	javaName := ToJavaName(table["table_name"])

	results := make([]Result, 0)
	for _, item := range items {
		result := Result{}
		hasPri := strings.Contains(item["column_key"], "PRI")
		if hasPri {
			result.Name = "id"
		} else {
			result.Name = "result"
		}
		result.Column = item["column_name"]
		result.JdbcType = JdbcTypes[item["data_type"]]
		result.Property = ToHumpField(item["column_name"])
		results = append(results, result)
	}
	xmlModel := &XmlModel{
		Mapper:  fmt.Sprintf("%s.%sMapper", MapperPackage, javaName),
		Model:   fmt.Sprintf("%s.%s", ModelPackage, javaName),
		Results: results,
	}

	filePath := OutFileDir + strings.Replace(XmlPackage, ".", "/", -1)
	CheckPath(filePath)
	file, err := os.OpenFile(filePath+"/"+javaName+"Mapper.xml", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}

	tpl.Execute(file, xmlModel)
}

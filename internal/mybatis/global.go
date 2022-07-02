package mybatis

import "database/sql"

var (
	BatisType     = "tk"
	OutFileDir    = "./out/"
	Author        string
	Table         = ""
	IsAllTables   = false
	Database      = ""
	Host          = "127.0.0.1"
	Port          = 3306
	UserName      = "root"
	Password      = "123456"
	ModelPackage  = ""
	MapperPackage = ""
	XmlPackage    = ""
	PrintHelp     = false
	DB            *sql.DB

	ModelTpl = `package {{.Package}};
{{range $item := .Annotations}}
{{$item}}
{{- end}}

{{.Descriptions}}
@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Table(name = "{{.TableName}}")
public class {{.Name}} implements Serializable {

	private static final long serialVersionUID = 1L;
{{range $item := .Fields}}
	{{$item.Comment}}
	{{- range $it := $item.Annotations}}
	{{$it}}
	{{- end}}
	{{$item.Field}}
{{end}}
}`

	MapperTpl = `package {{.Package}};
{{.Imports}}
import tk.mybatis.mapper.common.Mapper;
import tk.mybatis.mapper.common.MySqlMapper;

{{.Descriptions}}
public interface {{.Name}} extends Mapper<{{.Model}}>, MySqlMapper<{{.Model}}> {

}`

	XmlTpl = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="{{.Mapper}}">
	<resultMap id="BaseResultMap" type="{{.Model}}">
		{{- range $item := .Results}}
		<{{ $item.Name }} column="{{ $item.Column }}" jdbcType="{{ $item.JdbcType }}" property="{{ $item.Property }}" />
		{{- end}}
	</resultMap>
</mapper>`
)

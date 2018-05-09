package pandora
import (
	"github.com/kataras/iris"
	"database/sql"
	"strings"
)

func Pandora(db *sql.DB,ctx iris.Context) {
	strclientid := ctx.Params().Get("clientid")
	strservice := ctx.FormValue("service")
	if len(strclientid) < 5{
		ctx.WriteString("error clientid:"+strclientid)
		return;
	}
	substr := strings.Split(strclientid,":")
	if len(substr) < 6 {
		ctx.WriteString("error clientid:"+strclientid)
		return
	}
	strquery:= "select service.url from service,clientid where   clientid.game = "+substr[1] +" and clientid.ggi = "+substr[2]+ " and clientid.version = \""+substr[3]+"\" and service.name = \""+strservice+"\" and service.dc = clientid.dc"
	//uery:= "select name,url from service where name = \""+strservice+"\" and dc = \""
	//查询数据，指定字段名，返回sql.Rows结果集
	rows, _ := db.Query(strquery);
	defer rows.Close();
	for rows.Next(){
		url := ""
		rows.Scan(&url)
		if strservice == "pandora"{
			ctx.WriteString(url+"/"+strclientid)
		}else {
			ctx.WriteString(url)
		}

		return
	}
	ctx.WriteString("error service:"+strservice)

}
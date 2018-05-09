package eve

import (
"github.com/kataras/iris"
	"database/sql"
	"github.com/user/hello/clientId"
	"encoding/json"
	"fmt"
)

func Eve(db *sql.DB,ctx iris.Context) {
	strclientid := ctx.Params().Get("clientid")
	if clientId.CheckClientId(db,strclientid){
		strquery:= "select name,url from service"
		//查询数据，指定字段名，返回sql.Rows结果集
		rows, _ := db.Query(strquery);

		var serviceMap map[string]string
		/* 创建集合 */
		serviceMap = make(map[string]string)
		for rows.Next(){
			name := ""
			url := ""
			rows.Scan(&name,&url)
			fmt.Println(url, name);
			if name == "pandora" {
				serviceMap[name] = url + "/"+strclientid
			}else{
				serviceMap[name] = url
			}
		}
		jsonstr,_:=json.Marshal(serviceMap)
		ctx.WriteString(string(jsonstr))
		defer rows.Close();
	} else {
		ctx.WriteString("error clientid:"+strclientid)
	}
}
func EveWithDC(db *sql.DB,ctx iris.Context) {
	strclientid := ctx.Params().Get("clientid")
	strdc := ctx.Params().Get("dc")
	if clientId.CheckClientId(db,strclientid){
		strquery:= "select name,url from service where dc = \""+strdc+"\""
		//查询数据，指定字段名，返回sql.Rows结果集
		rows, _ := db.Query(strquery);

		var serviceMap map[string]string
		/* 创建集合 */
		serviceMap = make(map[string]string)
		for rows.Next(){
			name := ""
			url := ""
			rows.Scan(&name,&url)
			fmt.Println(url, name);
			if name == "pandora" {
				serviceMap[name] = url + "/"+strclientid
			}else{
				serviceMap[name] = url
			}
		}
		jsonstr,_:=json.Marshal(serviceMap)
		ctx.WriteString(string(jsonstr))
		defer rows.Close();
	} else {
		ctx.WriteString("error clientid:"+strclientid)
	}
}
func Datacenters(db *sql.DB,ctx iris.Context) {
	strclientid := ctx.Params().Get("clientid")
	ctx.WriteString(clientId.DataCenter(db,strclientid))
}
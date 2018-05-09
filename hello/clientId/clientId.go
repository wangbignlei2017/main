package clientId
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"fmt"
	"encoding/json"
)

func CheckClientId(db *sql.DB,strclientid string) bool{
	if len(strclientid) < 5{
		return false;
	}
	substr := strings.Split(strclientid,":")
	if len(substr) < 6 {
		return false;
	}
	strquery:= "select * from clientid where game = "+substr[1] +" and ggi = "+substr[2]+ " and version = \""+substr[3]+"\""
	//查询数据，指定字段名，返回sql.Rows结果集
	rows, _ := db.Query(strquery);
	defer rows.Close();
	for rows.Next() {
		return true
	}
	return false
}
func DataCenter(db *sql.DB,strclientid string) string{
	if len(strclientid) < 5{
		return "{}";
	}
	substr := strings.Split(strclientid,":")
	if len(substr) < 6 {
		return "{}";
	}
	strquery:= "select clientid.dc,datacenter.name from clientid,datacenter where clientid.dc = datacenter.name and  clientid.game = "+substr[1] +" and clientid.ggi = "+substr[2]+ " and clientid.version = \""+substr[3]+"\""
	//查询数据，指定字段名，返回sql.Rows结果集
	rows, _ := db.Query(strquery);
	fmt.Println(strquery);
	var serviceMap map[string]string
	/* 创建集合 */
	serviceMap = make(map[string]string)
	serviceMap["status"] = "acitve"
	serviceMap["preferred"] = "true"
	serviceMap["country_code"] = "CA"
	for rows.Next(){
		dc := ""
		name := ""
		rows.Scan(&dc,&name)
		fmt.Println(dc, name);
		serviceMap["name"] = name
		serviceMap["_datacenter_id"] = dc
	}
	defer rows.Close();
	jsonstr,_:=json.Marshal(serviceMap)
	return "["+string(jsonstr)+"]";
}
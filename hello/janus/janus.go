package janus

import (
	"github.com/kataras/iris"
	"database/sql"
	"github.com/user/hello/clientId"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

func Authorize(db *sql.DB,ctx iris.Context) {
	strclientid := ctx.FormValue("client_id")
	scope:=ctx.FormValue("scope")
	access_token_only:=ctx.FormValue("access_token_only")
	username:=ctx.FormValue("username")
	password:=ctx.FormValue("password")
	fmt.Println(access_token_only, scope);
	if clientId.CheckClientId(db,strclientid){
		strquery:= "SELECT credential, password, fedid FROM credential where credential = \""+username+"\""
		//查询数据，指定字段名，返回sql.Rows结果集
		rows, _:= db.Query(strquery);
		fedid := ""
		if rows.Next(){
			credential := ""
			dbpassword := ""

			rows.Scan(&credential,&dbpassword,&fedid)
			println(credential,dbpassword,fedid)
			if password != dbpassword{
				fmt.Println(credential, fedid);
			}
		}else {
			fedid = CreateCredential(db,strclientid,username,password)
			ctx.WriteString("not hit:"+strclientid)
		}
		println(fedid)
		defer rows.Close();
		if verifyscope(db,fedid,scope) {
			acctoken:= genaccesstoken(fedid,username,strclientid,scope)
			ctx.WriteString(acctoken)
		}
	} else {
		ctx.WriteString("error clientid:"+strclientid)
	}
}
func CreateCredential(db *sql.DB,strclientid string,username string,password string) string{
	u1,_:= uuid.NewV4()
	struuid:=u1.String()
	t := time.Now().Format("2006-01-02 15:04:05")
	strquery:= "insert into account values(\""+struuid+"\",\""+t+"\",0)"
	_,err:=db.Exec(strquery)
	if err != nil {
		println(err.Error())
	}
	println(strquery)
	strquery= "insert into credential values(\""+username+"\",\""+password+"\",\""+struuid+"\",\""+t+"\",\""+t+"\")"
	_,err=db.Exec(strquery)
	if err != nil {
		println(err.Error())
	}
	println(strquery)
	strquery= "insert into account_clientid(fed_id, clientid, create_time) values(\""+struuid+"\",\""+strclientid+"\",\""+t+"\",0)"
	_,err=db.Exec(strquery)
	if err != nil {
		println(err.Error())
	}
	println(strquery)

	strquery= "insert into fed_authority(fed_id) values(\""+struuid+"\")"
	_,err=db.Exec(strquery)
	if err != nil {
		println(err.Error())
	}
	println(strquery)

	return struuid
}
func verifyscope(db *sql.DB,fedid string,scopes string) bool{
	strscopes:= strings.Replace(scopes," ",",",-1)
	println(strscopes)
	strquery:= "SELECT  "+strscopes+" FROM fed_authority where fed_id = \""+fedid+"\""
	//查询数据，指定字段名，返回sql.Rows结果集
	rows, err:= db.Query(strquery);
	println(strquery)
	rows.Close()
	if err == nil {
		if rows.Next() {
			substr := strings.Split(strscopes,",")
			scopevaluse:= make([]bool, len(substr))
			rows.Scan(&scopevaluse)
			println(scopevaluse)
			for i:=0;i<len(scopevaluse);i++{
				if !scopevaluse[i]{
					println("scope miss %s",substr[i])
					return false
				}
			}

		}
	}else {
		println("scope miss %s","not hit")
		return false;
	}

	return true
}
func genaccesstoken(fedid string,username string,strclientid string,scope string)string {
	t := time.Now().Unix()

	strall:=fedid+","+scope+","+strclientid+","+strconv.FormatInt(t,10)+","+username+",,"+"mdc"
	sk:="1234567890abc"
	alldata:=strall+sk
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(alldata))
	cipherStr := md5Ctx.Sum(nil)
	fmt.Println(cipherStr)
	return strall+"|"+hex.EncodeToString(cipherStr)
}
func VerifyAccessToken(accesstoken string) bool{
	sk:="1234567890abc"
	index:= strings.LastIndex(accesstoken,"|")
	if index > 0{
		info:=string([]rune(accesstoken)[:index])
		token:=string([]rune(accesstoken)[index+1:])
		alldata:=info+sk
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(alldata))
		cipherStr := md5Ctx.Sum(nil)
		realtoken:=hex.EncodeToString(cipherStr)

		if strings.Compare(realtoken,token) != 0{
			return false
		}
	}else {
		return false
	}
	return true
}
func VerifyAccessTokenForScope(accesstoken string,scope string) bool {
	if VerifyAccessToken(accesstoken){
		substr:=strings.Split(accesstoken,",")
		if len(substr) > 5 {
			substr=strings.Split( substr[1]," ")
			for i:=0;i<len(substr);i++{
				if strings.Compare(scope,substr[i]) == 0 {

					return true
				}
			}
		}
	}
	return false
}
func GetFedId(db *sql.DB,credential string )string{
	fedid:=""
	strquery:= "SELECT fedid  FROM credential where credential = \""+credential+"\""
	//查询数据，指定字段名，返回sql.Rows结果集
	rows, err:= db.Query(strquery);
	rows.Close()
	if err != nil {
		if rows.Next() {

			rows.Scan(&fedid)
		}
	}
	return fedid
}
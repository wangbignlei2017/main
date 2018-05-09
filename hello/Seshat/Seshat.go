package Seshat
import (
	"database/sql"
	"github.com/kataras/iris"
	"fmt"
	"github.com/user/hello/janus"
	"strings"
	"github.com/user/hello/Status"
	"crypto/sha1"
	"encoding/json"
)

func Storage(db *sql.DB,ctx iris.Context)  {
	path:= ctx.Params().Get("p")
	fmt.Println(path)
	access_token:=ctx.FormValue("access_token")
	operation:=ctx.FormValue("operation")
	object:=ctx.FormValue("object")
	data:=ctx.FormValue("data")
	bother:=false
	if strings.Index(path,"data/me/") == 0 || strings.Index(path,"profiles/me/")==0  {
		if !janus.VerifyAccessTokenForScope(access_token,"storage") {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.WriteString(string(Status.Unauthorized)+" Unauthorized : HMAC does not match")

			return
		}
	}else if strings.Index(path,"data/") == 0 || strings.Index(path,"profiles/")==0 {
		if !janus.VerifyAccessTokenForScope(access_token,"storage_restricted") {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.WriteString(string(Status.Unauthorized)+" Unauthorized : HMAC does not match")
			return
		}
		bother = true
	}else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	substr:=strings.Split(access_token,"|")
	if len(substr) < 2 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	substr = strings.Split(substr[0],",")
	if len(substr) < 7 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	fedid:=substr[0]
	//scopes:=substr[1]
	clientid:=substr[2]
	//credential:=substr[4]
	substr = strings.Split(clientid,":")
	if len(substr) < 5{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	ggi:= substr[2]
	substr = strings.Split(path,"/")
	if bother{
		if len(substr) < 3{
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(string(Status.BadRequest)+" BadRequest")
			return
		}
		if strings.Index(substr[1],"fed_id:") > -1{
			tmpsubstr:=strings.Split(substr[1],":")
			fedid = tmpsubstr[1]
		}else{
			fedid = janus.GetFedId(db,substr[1])
		}

		if fedid == ""{
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(string(Status.BadRequest)+" BadRequest")
			return
		}
	}
	prekey:=substr[0]+"/"+fedid+"/"+ggi
	mykey:=substr[len(substr)-1]
	if substr[0] == "profiles" && substr[2] == "myprofile"{
		prekey =substr[0]+"/"+fedid+"/"+ggi+"/"+substr[2]
		if len(substr) == 3{
			mykey=""
		}
	}
	fmt.Println(object)
	if operation=="" || operation=="set"{
		if object !=""{
			SetDataByKey(db,mykey,prekey,ctx,object)
		}else{
			SetDataByKey(db,mykey,prekey,ctx,data)
		}
	}
	if operation=="merge"{
		if object !=""{
			MergeDataByKey(db,mykey,prekey,ctx,object)
		}else{
			MergeDataByKey(db,mykey,prekey,ctx,data)
		}
	}

	if operation=="batch_delete"{
		BatchDeleteDataByKey(db,mykey,prekey,data,ctx)
	}
	if operation=="delete"{
		DeleteDataByKey(db,mykey,prekey,ctx)
	}

}
func StorageGet(db *sql.DB,ctx iris.Context)  {
	path:= ctx.Params().Get("p")
	fmt.Println(path)
	access_token:=ctx.FormValue("access_token")
	//operation:=ctx.FormValue("operation")
	//object:=ctx.FormValue("object")
	bother:=false
	if strings.Index(path,"data/me/") == 0 || strings.Index(path,"profiles/me/")==0  {
		if !janus.VerifyAccessTokenForScope(access_token,"storage") {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.WriteString(string(Status.Unauthorized)+" Unauthorized : HMAC does not match")
			return
		}
	}else if strings.Index(path,"data/") == 0 || strings.Index(path,"profiles/")==0 {
		if !janus.VerifyAccessTokenForScope(access_token,"storage_ro") {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.WriteString(string(Status.Unauthorized)+" Unauthorized : HMAC does not match")
			return
		}
		bother = true
	}else {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	substr:=strings.Split(access_token,"|")
	if len(substr) < 2 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	substr = strings.Split(substr[0],",")
	if len(substr) < 7 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	fedid:=substr[0]
	//scopes:=substr[1]
	clientid:=substr[2]
	//credential:=substr[4]
	substr = strings.Split(clientid,":")
	if len(substr) < 5{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(string(Status.BadRequest)+" BadRequest")
		return
	}
	ggi:= substr[2]
	substr = strings.Split(path,"/")
	if bother{
		if len(substr) < 3{
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(string(Status.BadRequest)+" BadRequest")
			return
		}
		fmt.Println(strings.Index(substr[1],"fed_id:"))
		if strings.Index(substr[1],"fed_id:") > -1{
			tmpsubstr:=strings.Split(substr[1],":")
			fedid = tmpsubstr[1]
		}else{
			fedid = janus.GetFedId(db,substr[1])
		}
		if fedid == ""{
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(string(Status.BadRequest)+" BadRequest")
			return
		}
	}
	prekey:=substr[0]+"/"+fedid+"/"+ggi
	mykey:=substr[len(substr)-1]
	if substr[0] == "profiles" && substr[2] == "myprofile"{
		prekey =substr[0]+"/"+fedid+"/"+ggi+"/"+substr[2]
		if len(substr) == 3{
			mykey=""
		}
	}
	GetDataByKey(db,mykey,prekey,ctx)
}
func GetDataByKey(db *sql.DB,key string,prekey string,ctx iris.Context){
	substr:=strings.Split(key,".")
	mykey:=substr[0]
	hash:= GetHashcode(prekey+"/"+mykey)
	prehash:=GetHashcode(prekey)
	println(prehash)
	if mykey ==""{
		strquery:= "SELECT sha, data, pre_sha, key  FROM  seshat where pre_sha = '"+string( prehash)+"'"
		println(strquery)
		rows, err:= db.Query(strquery);
		var mapvaluse map[string]string
		//mapvaluse = make(map[strings]strings)
		if err != nil {
			for rows.Next(){
				sha:=""
				data:=""
				per_sha:=""
				tmpkey:=""
				rows.Scan(&sha,&data,&per_sha,&tmpkey)
				mapvaluse[tmpkey] = data
			}
			jsonobj,_:=json.Marshal(mapvaluse)
			ctx.WriteString(string(jsonobj))
		}else {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.WriteString("not fount")

		}
		rows.Close()
		return
	}else{
		strquery:= "SELECT sha, data, pre_sha,mykey from seshat where pre_sha = \""+string( prehash)+"\" and sha =\""+hash+"\""
		println(strquery)
		rows, err:= db.Query(strquery);
		if err == nil {
			if rows.Next(){
				sha:=""
				data:=""
				per_sha:=""
				tmpkey:=""
				rows.Scan(&sha,&data,&per_sha,&tmpkey)
				fmt.Println(sha,data,per_sha,tmpkey)
				if len(substr) == 1{
					ctx.WriteString(data)
				}else{
					var tmp map[string]interface{}
					newerr:=json.Unmarshal([]byte(data),&tmp)
					if newerr == nil{
						var tmpv map[string]interface{}
						nothit:=false
						tmpv = tmp
						for i:=1;i<len(substr);i++{
							v,ok:=tmpv[substr[i]]
							println(tmpv)
							println(v,ok)
							if ok {
								switch e := v.(type) {
								case map[string]interface{}:
									var s int
									tmpv = e
									fmt.Println(s)
								default:
									nothit = true
									break
								}
							}else{
								nothit = true
								break
							}
						}
						if nothit {
							ctx.StatusCode(iris.StatusNotFound)
							ctx.WriteString("not fount")
						}else{
							jsonobj,_:=json.Marshal(tmpv)
							ctx.WriteString(string(jsonobj))
						}
					}else{
						ctx.StatusCode(iris.StatusNotFound)
						ctx.WriteString("not fount")
						fmt.Println(newerr.Error(),data)
					}
				}
			}else{
				fmt.Println("not fount")
				ctx.StatusCode(iris.StatusNotFound)
				ctx.WriteString("not fount")
			}
			rows.Close()
		}else {
			println(err.Error())
			ctx.StatusCode(iris.StatusNotFound)
			ctx.WriteString("not fount")
		}

	}
}
func GetHashcode(key string) string{
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(key))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)
	println(key,bs)

	return fmt.Sprintf("%x", bs)
}

func DeleteDataByKey(db *sql.DB,key string,prekey string,ctx iris.Context){
	substr:=strings.Split(key,".")
	mykey:=substr[0]
	hash:= GetHashcode(prekey+"/"+mykey)
	prehash:=GetHashcode(prekey)
	if mykey ==""{
		strquery:= "delete  seshat where pre_sha = '"+string( prehash)+"'"
		println(strquery)
		db.Exec(strquery)
	}else {
		if len(substr) == 0{
			strquery:= "delete  seshat where sha = '"+string( hash)+"'"
			println(strquery)
			db.Exec(strquery)
		}else {
			strquery:= "SELECT sha, data, pre_sha, mykey  FROM  seshat where pre_sha = '"+string( prehash)+"' and sha ='"+hash+"'"
			println(strquery)
			rows, err:= db.Query(strquery);
			if err != nil {
				if rows.Next(){
					sha:=""
					data:=""
					per_sha:=""
					tmpkey:=""
					rows.Scan(&sha,&data,&per_sha,&tmpkey)
					var tmp map[string]interface{}
					newerr:=json.Unmarshal([]byte(data),&tmp)
					if newerr != nil{
						var tmpv * map[string]interface{}
						nothit:=false
						tmpv = &tmp
						for i:=1;i<len(substr)-1;i++{
							v,ok:=(*tmpv)[substr[i]]
							if ok {
								switch e := v.(type) {
								case map[string]interface{}:
									var s int
									tmpv = &e
									fmt.Println(s)
								default:
									nothit = true
									break
								}
							}else{
								nothit = true
								break
							}
						}
						if !nothit{
							delete((*tmpv),substr[len(substr)-1])
						}
						jsonobj,_:=json.Marshal(tmp)
						strquery ="update seshat set data ='"+string(jsonobj)+"' where sha='"+hash+"'"
						db.Exec(strquery)
					}
				}
			}
			rows.Close()
		}
	}
}
func BatchDeleteDataByKey(db *sql.DB,key string,prekey string,data string,ctx iris.Context){
	substr:=strings.Split(data,",")
	for i:=0;i<len(substr);i++ {
		if key == ""{
			DeleteDataByKey(db,substr[i],prekey,ctx)
		}else {
			DeleteDataByKey(db,key+"."+substr[i],prekey,ctx)
		}
	}
}
func SetDataByKey(db *sql.DB,key string,prekey string,ctx iris.Context,data string){
	substr:=strings.Split(key,".")
	mykey:=substr[0]
	path:=prekey+"/"+mykey
	hash:= GetHashcode(prekey+"/"+mykey)
	prehash:=GetHashcode(prekey)
	if mykey ==""{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("BadRequest")
	}else {
		if len(substr) == 0{
			strquery:= "INSERT INTO seshat VALUES ('"+hash+"','"+data+"','"+prehash+"','"+mykey+"','"+path+"') ON DUPLICATE KEY UPDATE b=VALUES(b),c=VALUES(c)"
			fmt.Println(strquery)
			db.Exec(strquery)
		}else {
			strquery:= "SELECT sha, data, pre_sha, mykey  FROM  seshat where pre_sha = \""+string( prehash)+"\" and sha =\""+hash+"\""
			fmt.Println(strquery)
			rows, err:= db.Query(strquery);
			if err == nil {
				if rows.Next(){
					sha:=""
					mydata:=""
					per_sha:=""
					tmpkey:=""
					rows.Scan(&sha,&mydata,&per_sha,&tmpkey)
					var tmp map[string]interface{}
					newerr:=json.Unmarshal([]byte(mydata),&tmp)
					if newerr == nil{
						var tmpv * map[string]interface{}
						nothit:=false
						tmpv = &tmp
						for i:=1;i<len(substr)-1;i++{
							v,ok:=(*tmpv)[substr[i]]
							if ok {
								switch e := v.(type) {
								case map[string]interface{}:
									var s int
									tmpv = &e
									fmt.Println(s)
								default:
									nothit = true
									break
								}
							}else{
								nothit = true
								break
							}
						}
						if !nothit {
							var dataobj interface{}
							dataerr:=json.Unmarshal([]byte(data),&dataobj)
							(*tmpv)[substr[len(substr)-1]] = dataerr
						}
						jsonobj,_:=json.Marshal(tmp)
						strquery ="update seshat set data ='"+string(jsonobj)+"' where sha=\""+hash+"\""
						fmt.Println(strquery)
						db.Exec(strquery)
					}else{
							str:=data
							for i:=len(substr)-1;i>0;i--{
								var m map[string]interface{}
								m= make(map[string]interface{})
								var dataobj interface{}
								dataerr:=json.Unmarshal([]byte(str),&dataobj)
								fmt.Println(str)
								if dataerr != nil{
									fmt.Println(dataerr)
								}else{
									fmt.Println(dataobj)
								}
								m[substr[i]]=dataobj
								jsonobj,_:=json.Marshal(m)
								str= string(jsonobj)
							}
							strquery:= "INSERT INTO seshat VALUES ('"+hash+"','"+str+"','"+prehash+"','"+mykey+"','"+path+"') ON DUPLICATE KEY UPDATE data=VALUES(data)"
							fmt.Println(strquery)
							_,exerr:=db.Exec(strquery)
							if exerr!=nil{
								fmt.Println(exerr.Error())
								ctx.StatusCode(iris.StatusBadRequest)
								ctx.WriteString(exerr.Error())
								fmt.Println(exerr.Error())
							}
					}

				}else{
					str:=data

					for i:=len(substr)-1;i>0;i--{
						var m map[string]interface{}
						m= make(map[string]interface{})
						var dataobj interface{}
						dataerr:=json.Unmarshal([]byte(str),&dataobj)
						fmt.Println(str)
						if dataerr != nil{
							fmt.Println(dataerr)
						}else{
							fmt.Println(dataobj)
						}
						m[substr[i]]=dataobj
						jsonobj,_:=json.Marshal(m)
						str= string(jsonobj)
					}
					strquery:= "INSERT INTO seshat VALUES ('"+hash+"','"+str+"','"+prehash+"','"+mykey+"','"+path+"') ON DUPLICATE KEY UPDATE data=VALUES(data)"
					fmt.Println(strquery)
					_,exerr:=db.Exec(strquery)
					if exerr!=nil{
						fmt.Println(exerr.Error())
						ctx.StatusCode(iris.StatusBadRequest)
						ctx.WriteString(exerr.Error())
						fmt.Println(exerr.Error())
					}

				}
			}else{
				fmt.Println(err.Error())
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.WriteString(err.Error())
			}
			rows.Close()
		}
	}
}
func MergeDataByKey(db *sql.DB,key string,prekey string,ctx iris.Context,data string){
	substr:=strings.Split(key,".")
	mykey:=substr[0]
	path:=prekey+"/"+mykey
	hash:= GetHashcode(path)

	prehash:=GetHashcode(prekey)
	if mykey ==""{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("BadRequest")
	}else {
		if len(substr) == 0{
			strquery:= "INSERT INTO seshat VALUES ('"+hash+"','"+data+"','"+prehash+"','"+mykey+"','"+path+"') ON DUPLICATE KEY UPDATE b=VALUES(b),c=VALUES(c)"
			fmt.Println(strquery)
			db.Exec(strquery)
		}else {
			strquery:= "SELECT sha, data, pre_sha, mykey  FROM  seshat where pre_sha = \""+string( prehash)+"\" and sha =\""+hash+"\""
			fmt.Println(strquery)
			rows, err:= db.Query(strquery);
			if err == nil {
				if rows.Next(){
					sha:=""
					mydata:=""
					per_sha:=""
					tmpkey:=""
					rows.Scan(&sha,&mydata,&per_sha,&tmpkey)
					var tmp map[string]interface{}
					newerr:=json.Unmarshal([]byte(mydata),&tmp)

					if newerr == nil{
						var tmpv * map[string]interface{}
						nothit:=false
						tmpv = &tmp
						for i:=1;i<len(substr);i++{
							v,ok:=(*tmpv)[substr[i]]
							if ok {
								switch e := v.(type) {
								case map[string]interface{}:
									tmpv = &e
									fmt.Println(*tmpv)
								default:
									nothit = true
									fmt.Println(*tmpv,e)
									break
								}
							}else{
								nothit = true
								fmt.Println(*tmpv)
								break
							}
						}
						if !nothit {
							var newtmp map[string]interface{}
							json.Unmarshal([]byte(data),&newtmp)
							fmt.Println(newtmp)
							for k,v :=range newtmp{
								(*tmpv)[k] = v
							}
							fmt.Println(*tmpv,nothit)
						}
						jsonobj,_:=json.Marshal(tmp)
						strquery ="update seshat set data = '"+string(jsonobj)+"' where sha=\""+hash+"\""
						fmt.Println(strquery)
						db.Exec(strquery)
					}else{
						ctx.StatusCode(iris.StatusBadRequest)
						ctx.WriteString(newerr.Error())
						fmt.Println(newerr.Error())
					}

				}else{
					str:=data
					for i:=len(substr)-1;i>0;i--{
						var m map[string]interface{}
						m= make(map[string]interface{})
						var newtmp interface{}
						json.Unmarshal([]byte(str),&newtmp)
						m[substr[i]]=newtmp
						jsonobj,_:=json.Marshal(m)
						str= string(jsonobj)
					}
					strquery:= "INSERT INTO seshat VALUES ('"+hash+"','"+str+"','"+prehash+"','"+mykey+"','"+path+"') ON DUPLICATE KEY UPDATE data=VALUES(data)"
					fmt.Println(strquery)
					_,exerr:=db.Exec(strquery)
					if exerr!=nil{
						fmt.Println(exerr.Error())
						ctx.StatusCode(iris.StatusBadRequest)
						ctx.WriteString(exerr.Error())

					}

				}
			}else{
				println(err.Error())
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.WriteString(err.Error())
			}
			rows.Close()
		}
	}
}
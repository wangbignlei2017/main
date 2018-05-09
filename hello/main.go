package main

import (
	"github.com/kataras/iris"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/user/hello/eve"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/user/hello/pandora"
	"encoding/json"
	"fmt"
	"crypto/sha1"
	"github.com/user/hello/janus"
	"github.com/user/hello/Seshat"
)
var db *sql.DB
func main() {
	app := iris.New()
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/fed_db?charset=utf8")
	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		name := "{\"version\": \"3.9.12-0\", \"eve\": \"on\"}"
		ctx.WriteString(name)
	})
	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/config
	app.Get("/config", func(ctx iris.Context) {
		name := "{\"version\": \"3.9.12-0\", \"eve\": \"on\"}"
		ctx.WriteString(name)
	})

	app.Get("/eve/config/{clientid}", func(ctx iris.Context) {
		eve.Eve(db,ctx);
	})
	app.Get("/eve/config/{clientid}/datacenters/{dc}/urls", func(ctx iris.Context) {
		eve.EveWithDC(db,ctx);
	})
	app.Get("/eve/config/{clientid}/datacenters", func(ctx iris.Context) {
		eve.Datacenters(db,ctx);
	})
	app.Get("/pandora/{clientid}/locate", func(ctx iris.Context) {
		pandora.Pandora(db,ctx)
	})
	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})
	app.Get("/test", func(ctx iris.Context) {
		data := `{"account":"d98eedfa-4c5f-11e8-b1e4-b8ca3a636aa8","client_ids":[],"installations":[{"carrier":null,"language":null,"country":null,"download_code":null,"model":null,"firmware":null,"resolution":null,"device_id":"3067551625"}],"alias":[],"last_login":null,"is_ghost":false,"credentials":["anonymous:d2luMzJfACQAAPRSAABSAAAAJCEAAA=="]}`

		var dbgInfos map[string]interface{}
		err:=json.Unmarshal([]byte(data), &dbgInfos)
		if len(dbgInfos) > 0 {
			fmt.Println(111)
		}else {
			fmt.Println(222)
		}
		fmt.Println(err)
		fmt.Println(len(dbgInfos))
		fmt.Println(dbgInfos)
		s := "sha1 this string"

		// The pattern for generating a hash is `sha1.New()`,
		// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
		// Here we start with a new hash.
		h := sha1.New()

		// `Write` expects bytes. If you have a string `s`,
		// use `[]byte(s)` to coerce it to bytes.
		h.Write([]byte(data))

		// This gets the finalized hash result as a byte
		// slice. The argument to `Sum` can be used to append
		// to an existing byte slice: it usually isn't needed.
		bs := h.Sum(nil)

		// SHA1 values are often printed in hex, for example
		// in git commits. Use the `%x` format verb to convert
		// a hash results to a hex string.
		fmt.Println(s)
		fmt.Printf("%x\n", bs)

	})
	app.Post("/janus/authorize",func(ctx iris.Context) {
		janus.Authorize(db,ctx)
	})
	app.Post("/janus/verify",func(ctx iris.Context) {
		token:= ctx.FormValue("access_token")
		janus.VerifyAccessToken(token)
	})
	app.Post("/storage/{p:path}",func(ctx iris.Context) {
		Seshat.Storage(db,ctx)
	})
	app.Get("/storage/{p:path}",func(ctx iris.Context) {
		Seshat.StorageGet(db,ctx)
	})
	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	t2t "github.com/takanoriyanagitani/go-pg2txt/table2txt"
	t2j "github.com/takanoriyanagitani/go-pg2txt/table2txt/std"
	p2t "github.com/takanoriyanagitani/go-pg2txt/table2txt/std/pgx"
)

func must[T any](t T, e error) T {
	if nil != e {
		panic(e)
	}
	return t
}

var database *sql.DB = must(p2t.NewDB(""))

var tab2jsons t2j.Tab2Jsons = t2j.Tab2JsonsNew(database)
var tab2str t2t.TableToString = tab2jsons.AsIf()

var tabchki t2j.TabChkInfo = t2j.TabChkInfoNew(database)
var tabchk t2t.TableChecker = tabchki.AsIf()
var tchkf t2t.TabChkFn = tabchk.Check

var t2c t2t.TableToString = tchkf.ToChecked(tab2str)

var tableName string = os.Getenv("ENV_TABLE_NAME")

func main() {
	defer database.Close()
	if 0 == len(tableName) {
		log.Printf("table name missing: ENV_TABLE_NAME")
		return
	}
	var jsons []string = must(t2c.Query(
		context.Background(),
		tableName,
	))
	for _, j := range jsons {
		fmt.Println(j)
	}
}

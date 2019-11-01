package main

import (
	"fmt"
	"github.com/masami10/rush/services/testORM/orm"

	_ "github.com/lib/pq"

	"io/ioutil"
	"os"
)

//`xorm:"varchar(25) notnull unique 'usr_name'"`
func main() {

	file, _ := os.Open("/home/kevin/Downloads/gopath/src/github.com/masami10/rush/services/testORM/workorder.json")

	var wor []byte
	wor, _ = ioutil.ReadAll(file)

	fmt.Println(string(wor))

	r := orm.NewRush()
	//插入新工单
	r.WorkorderIn(wor)

	fmt.Println("###############################################")
	//根据工单号,查询工单
	_, rr := r.WorkorderOut("MO0001")
	fmt.Println(string(rr))

	//results, err := engine.Query("select (jsons::json#>>'{product,url}')::text from testjson where (jsons::json#>>'{code}')::text = 'MO0002'")
	//
	//fmt.Println(string(results[0]["text"]))
	//	//则会在控制台打印出生成的SQL语句；
	//engine.ShowSQL(true)
	//	engine.Sync2(new(User))

}

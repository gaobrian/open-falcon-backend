package dashboard

import (
	"github.com/gaobrian/open-falcon-backend/modules/portal/g"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"fmt"
	"errors"
)

type Counter struct {
	Counter string
	Step    int
	Type    string
}

func QueryEndpintByNameRegx(queryStr string, limit int) (enp []Endpoint, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.Limit {
		limit = config.GraphDB.Limit
	}
	q := orm.NewOrm()
	q.Using("graph")
	_, err = q.Raw("select * from `endpoint` where endpoint regexp ? limit ?", queryStr, limit).QueryRows(&enp)
	return
}



func QueryCounterByNameRegx(endpoints string,queryStr string, limit int) (counters []Counter, err error) {
	config := g.Config()
	if limit == 0 || limit > config.GraphDB.Limit {
		limit = config.GraphDB.Limit
	}
	q := orm.NewOrm()
	q.Using("graph")

	jsonString := `{"hosts":` + endpoints +`}`
	var endps map[string][]string
	if err = json.Unmarshal([]byte(jsonString), &endps); err != nil {
		return nil,errors.New("")
	}

	if len(endps) == 0 {
		return nil,errors.New("")
	}

	sql_endp := `select id from endpoint where endpoint in (` + `"` + endps["hosts"][0] + `"`
	for i:=1;i<len(endps["hosts"]);i++{
		sql_endp = sql_endp + `,"` + endps["hosts"][i] + `"`
	}

	sql_endp = sql_endp + ")"


	var endpoint_ids []int
	_,err  = q.Raw(sql_endp).QueryRows(&endpoint_ids)

	if len(endpoint_ids) ==  0 || err!=nil {
		return nil,errors.New("")
	}

	if len(endpoint_ids) == 0{
		return nil,errors.New("")
	}

	sql_stmt := "select distinct(counter),type,step from endpoint_counter where endpoint_id in ("
	sql_stmt = sql_stmt + fmt.Sprintf("%d",endpoint_ids[0])
	for i:=1;i<len(endpoint_ids);i++{
		sql_stmt = sql_stmt +  fmt.Sprintf(",%d",endpoint_ids[i])
	}
	sql_stmt = sql_stmt + ") and counter regexp ? limit ?"

	_, err = q.Raw(sql_stmt ,queryStr, limit).QueryRows(&counters)
	return
}

func QueryEndpintByNameRegxForOps(queryStr string) (enp []Hosts, err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	_, err = q.Raw("select * from `host` where hostname regexp ?", queryStr).QueryRows(&enp)
	return
}

func CountNumOfHost() (c int, err error) {
	var h []Endpoint
	q := getOrmObj()
	_, err = q.Raw("select id from `endpoint`").QueryRows(&h)
	c = len(h)
	return
}

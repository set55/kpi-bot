package main

import (
	"database/sql"
	// dbQuery "kpi-bot/db"

	"kpi-bot/lib/bot"
	"fmt"
	"log"

	"kpi-bot/lib/rd"
	"kpi-bot/lib/test"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Define the data source name (DSN) Mh-mJ?sp.G"43*_HrCXRP9+^QS%3Et2yZE
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "developer", "Mh-mJ?sp.G\"43*_HrCXRP9+^QS%3Et2yZE", "192.168.2.8", "32606", "zentao")

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
		return
	}

	fmt.Println("Connected to MariaDB successfully!")

	// init bot, and date range
	beginDatetime := "2025-03-01 00:00:00"
	endDatetime := "2025-03-31 23:59:59"
	robot := bot.NewBot(db)


	// app 研发
	rds := []string{
		"alan.tin",
		"jihuaqing",
		"liuhongtao",
	}
	rdProjectPms := []string{
		"guoqiao.chen",
		"shawn.wang",
		"simon.chen",
		"qixiaofeng",
		"set.su",
		"justin.lee",
		"jiangjiahui",
		"caojianni",
	}
	err = robot.ProduceRdKpi("./excel/kpi-rd.xlsx", beginDatetime, endDatetime, rds, rdProjectPms)
	if err != nil {
		log.Fatalf("Error produceRdKpi: %v", err)
	}

	
	// other 研发
	rdsWithoutTest := []string{
		// embed system
		"shiwen.tin",
		"wangtuhe",
		"chenyuanchong",
		"qihongquan",
		"zhangzhilu",
		"zhuangjianyong",
		"wangxianming",
		// devops
		// "justin.lee",
		// 阿崔部门 研发
		"zengyi",
		"chenbo",
		"lixiaolong",
		"tangjilin",
		"jiaoxiangjie",
		"bieji",
		"suiguanyou",
		"lishuaipeng",
		"liuxiaoyun",
	}
	
	rdsWithoutTestProjectPms := []string{
		"guoqiao.chen",
		"shawn.wang",
		"simon.chen",
		"qixiaofeng",
		"set.su",
		"justin.lee",
		"caojianni",
	}
	err = robot.ProduceRdKpiWithoutTestReport("./excel/kpi-rd-without.xlsx", beginDatetime, endDatetime, rdsWithoutTest, rdsWithoutTestProjectPms)
	if err != nil {
		log.Fatalf("error produceRdKpiWithoutTestreport: %v", err)
	}

	

	// 项目经理
	pms := []string{
		"qixiaofeng",
		"jiangjiahui",
		"caojianni",
		"shawn.wang",
		"simon.chen",
	}
	err = robot.ProducePmKpi("./excel/pmnew.xlsx", beginDatetime, endDatetime, pms)
	if err != nil {
		log.Fatalf("error ProducePmKpi: %v", err)
	}

	// 项目经理(不包含测试)
	pmsWithoutTest := []string{"guoqiao.chen"}
	err = robot.ProducePmKpiWithoutTestReport("./excel/kpi-pm.xlsx", beginDatetime, endDatetime, pmsWithoutTest)
	if err != nil {
		log.Fatalf("error ProducePmKpiWithoutTestReport: %v", err)
	}

	// SSC RD
	rds2 := []string{
		"set.su", 
		"paul.gao",
		"justin.lee",
		"samy.gou", 
		"deakin.han", 
		"xiechen", 
		"zouyanling", 
		"ruanbanyong", 
		"zhouyao",
		"liuxiaoyan",
	}
	for _, v := range rds2 {
		tmp := rd.NewRdKpi2(db, v, beginDatetime, endDatetime)
		err := tmp.MakeRdReport("./excel/kpi-rd2-2.xlsx")
		if err != nil {
			log.Fatalf("error MakeRdReport: %v", err)
		}
	}

	// SSC TEST
	test2 := []string{
		"linyanhai",
		"wangshaoyu",
		"pengzijie",
		"xiezhiren",
	}
	for _, v := range test2 {
		tmp := test.NewTestKpi2(db, v, beginDatetime, endDatetime)
		err := tmp.MakeTestReport("./excel/kpi-test2.xlsx")
		if err != nil {
			log.Fatalf("error MakeTestReport: %v", err)
		}
	}

}
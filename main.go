package main

import (
	"database/sql"
	// dbQuery "kpi-bot/db"

	"kpi-bot/lib/bot"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Define the data source name (DSN)
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

	
	// all team members
	rds := []string{"set.su", "paul.gao", "samy.gou", "champion.fu", "alan.tin", "jihuaqing", "liuhongtao", "deakin.han"}
	rdsWithoutTest := []string{"shiwen.tin", "xiechen", "zouyanling", "ruanbanyong", "zhouyao", "liuxiaoyan", "wangtuhe"}
	tests := []string{"linyanhai", "wangshaoyu"}
	pms := []string{"qixiaofeng"}
	pmsWithoutTest := []string{"guoqiao.chen", "shawn.wang"}
	deveops := []string{"justin.lee"}

	beginDatetime := "2024-09-01 00:00:00"
	endDatetime := "2024-09-30 23:59:59"



	robot := bot.NewBot(db)
	err = robot.ProduceRdKpi("./excel/绩效考核模板-研发.xlsx", beginDatetime, endDatetime, rds)
	if err != nil {
		log.Fatalf("Error produceRdKpi: %v", err)
	}

	err = robot.ProduceRdKpiWithoutTestReport("./excel/绩效考核模板-研发(无测试报告).xlsx", beginDatetime, endDatetime, rdsWithoutTest)
	if err != nil {
		log.Fatalf("error produceRdKpiWithoutTestreport: %v", err)
	}

	err = robot.ProducePmKpi("./excel/绩效考核模板-项目.xlsx", beginDatetime, endDatetime, pms)
	if err != nil {
		log.Fatalf("error ProducePmKpi: %v", err)
	}

	err = robot.ProducePmKpiWithoutTestReport("./excel/绩效考核模板-项目.xlsx", beginDatetime, endDatetime, pmsWithoutTest)
	if err != nil {
		log.Fatalf("error ProducePmKpiWithoutTestReport: %v", err)
	}

	err = robot.ProduceTestKpi("./excel/绩效考核模板-测试.xlsx", beginDatetime, endDatetime, tests)
	if err != nil {
		log.Fatalf("error ProduceTestKpi: %v", err)
	}

	err = robot.ProduceStatisticKpi("./excel/kpi统计.xlsx", beginDatetime, endDatetime, rds, rdsWithoutTest, pms, pmsWithoutTest, tests)
	if err != nil {
		log.Fatalf("error ProduceStatisticKpi: %v", err)
	}

	err = robot.ProduceDeveopsKpi("./excel/绩效考核模板-运维.xlsx", beginDatetime, endDatetime, deveops)
	if err != nil {
		log.Fatalf("error ProduceDeveopsKpi: %v", err)
	}

	// dbQuery.QueryRdBugCarryOver(db, rds, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryRdBugCarryOverDetail(db, rds, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryRdBugCarryOverWithoutTestReport(db, rds, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryRdBugCarryOverDetailWithoutTestReport(db, rds, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryRdProjectProgress(db, []string{"zhouyao"}, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryRdProjectProgressDetail(db, []string{"zhouyao"}, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryProjectProgress(db, []string{"shawn.wang", "qixiaofeng", "guoqiao.chen"}, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryProjectProgressWithout(db, []string{"shawn.wang", "qixiaofeng", "guoqiao.chen"}, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryProjectCompleteRate(db, []string{"shawn.wang", "qixiaofeng", "guoqiao.chen"}, "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// dbQuery.QueryProjectCompleteRateDetail(db, []string{"shawn.wang", "qixiaofeng", "guoqiao.chen"}, "2024-07-01 00:00:00", "2024-07-31 23:59:59")

}

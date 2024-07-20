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
	rds := []string{"set.su", "paul.gao", "justin.lee", "samy.gou", "champion.fu", "alan.tin", "jihuaqing", "liuhongtao", "deakin.han"}
	rdsWithoutTest := []string{"shiwen.tin", "xiechen", "zouyanling", "ruanbanyong", "zhouyao", "liuxiaoyan", "wangtuhe"}
	tests := []string{"linyanhai", "wangshaoyu"}
	pms := []string{"shawn.wang", "qixiaofeng"}
	pmsWithoutTest := []string{"guoqiao.chen"}


	robot := bot.NewBot(db)
	err = robot.ProduceRdKpi("./excel/绩效考核模板-研发.xlsx", "2024-07-01 00:00:00", "2024-07-31 23:59:59", rds)
	if err != nil {
		log.Fatalf("Error produceRdKpi: %v", err)
	}

	err = robot.ProduceRdKpiWithoutTestReport("./excel/绩效考核模板-研发(无测试报告).xlsx", "2024-07-01 00:00:00", "2024-07-31 23:59:59", rdsWithoutTest)
	if err != nil {
		log.Fatalf("error produceRdKpiWithoutTestreport: %v", err)
	}

	err = robot.ProducePmKpi("./excel/绩效考核模板-项目.xlsx", "2024-07-01 00:00:00", "2024-07-31 23:59:59", pms)
	if err != nil {
		log.Fatalf("error ProducePmKpi: %v", err)
	}

	err = robot.ProducePmKpiWithoutTestReport("./excel/绩效考核模板-项目.xlsx", "2024-07-01 00:00:00", "2024-07-31 23:59:59", pmsWithoutTest)
	if err != nil {
		log.Fatalf("error ProducePmKpiWithoutTestReport: %v", err)
	}

	err = robot.ProduceTestKpi("./excel/绩效考核模板-测试.xlsx", "2024-07-01 00:00:00", "2024-07-31 23:59:59", tests)
	if err != nil {
		log.Fatalf("error ProduceTestKpi: %v", err)
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

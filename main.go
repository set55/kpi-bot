package main

import (
	"database/sql"
	// dbQuery "kpi-bot/db"

	// "kpi-bot/lib/bot"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"kpi-bot/lib/rd"
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

	
	// all team members

	// 研发
	rds := []string{"set.su", "paul.gao", "samy.gou", "alan.tin", "jihuaqing", "liuhongtao", "deakin.han"}
	// rdProjectPms := []string{"guoqiao.chen", "shawn.wang","simon.chen", "qixiaofeng", "set.su", "justin.lee", "jiangjiahui", "caojianni"}
	// rds := []string{"champion.fu"}
	
	
	// 软件服务中心 研发
	// rdsWithoutTest := []string{"shiwen.tin", "xiechen", "zouyanling", "ruanbanyong", "zhouyao", "liuxiaoyan", "wangtuhe", "chenyuanchong", "qihongquan", "zhangzhilun", "zhuangjianyong", "wangxianming"}
	// rdsWithoutTestProjectPms := []string{"guoqiao.chen", "shawn.wang","simon.chen", "qixiaofeng", "set.su", "justin.lee", "caojianni"}
	// rdsWithoutTest := []string{"liuxiaoyan"}

	// 阿崔部门 研发
	// rdsWithoutTest := []string{"zengyi", "chenbo", "lixiaolong", "tangjilin", "jiaoxiangjie", "bieji", "suiguanyou", "lishuaipeng"}

	// 测试
	// tests := []string{"linyanhai", "wangshaoyu", "pengzijie"}
	// tests := []string{"wangshaoyu"}

	// 项目经理
	// pms := []string{"qixiaofeng", "jiangjiahui", "caojianni"}
	// pmsWithoutTest := []string{"guoqiao.chen", "shawn.wang","simon.chen"}
	// pmsWithoutTest := []string{"guoqiao.chen"}

	// 运维
	// deveops := []string{"justin.lee"}

	beginDatetime := "2024-12-01 00:00:00"
	endDatetime := "2024-12-31 23:59:59"



	// robot := bot.NewBot(db)
	// err = robot.ProduceRdKpi("./excel/kpi-rd.xlsx", beginDatetime, endDatetime, rds, rdProjectPms)
	// if err != nil {
	// 	log.Fatalf("Error produceRdKpi: %v", err)
	// }

	// err = robot.ProduceRdKpiWithoutTestReport("./excel/kpi-rd-without.xlsx", beginDatetime, endDatetime, rdsWithoutTest, rdsWithoutTestProjectPms)
	// if err != nil {
	// 	log.Fatalf("error produceRdKpiWithoutTestreport: %v", err)
	// }

	// err = robot.ProducePmKpi("./excel/kpi-pm.xlsx", beginDatetime, endDatetime, pms)
	// if err != nil {
	// 	log.Fatalf("error ProducePmKpi: %v", err)
	// }

	// err = robot.ProducePmKpiWithoutTestReport("./excel/kpi-pm.xlsx", beginDatetime, endDatetime, pmsWithoutTest)
	// if err != nil {
	// 	log.Fatalf("error ProducePmKpiWithoutTestReport: %v", err)
	// }

	// storyPms := []string{"qixiaofeng", "guoqiao.chen", "shawn.wang", "simon.chen", "huangweiqi"}
	// err = robot.ProduceTestKpi("./excel/kpi-test.xlsx", beginDatetime, endDatetime, tests, storyPms)
	// if err != nil {
	// 	log.Fatalf("error ProduceTestKpi: %v", err)
	// }

	// err = robot.ProduceStatisticKpi("./excel/kpi-statistic.xlsx", beginDatetime, endDatetime, rds, rdsWithoutTest, pms, pmsWithoutTest, tests, storyPms, rdProjectPms, rdsWithoutTestProjectPms)
	// if err != nil {
	// 	log.Fatalf("error ProduceStatisticKpi: %v", err)
	// }

	// err = robot.ProduceDeveopsKpi("./excel/kpi-devops.xlsx", beginDatetime, endDatetime, deveops, rdProjectPms)
	// if err != nil {
	// 	log.Fatalf("error ProduceDeveopsKpi: %v", err)
	// }

	// err = robot.ProduceWhatEverKpi("./excel/管理岗人才胜任力盘点2025.xlsx")
	// if err != nil {
	// 	log.Fatalf("error ProduceWhatEverKpi: %v", err)
	// }

	// err = robot.ProduceWhatEverNormalKpi("./excel/一般员工人才胜任力盘点2025.xlsx")
	// if err != nil {
	// 	log.Fatalf("error ProduceWhatEverNormalKpi: %v", err)
	// }



	for _, v := range rds {
		tmp := rd.NewRdKpi2(db, v, beginDatetime, endDatetime)
		err := tmp.MakeRdReport("./excel/kpi-rd2.xlsx")
		if err != nil {
			log.Fatalf("error MakeRdReport: %v", err)
		}

	}
}
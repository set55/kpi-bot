package main

import (
	"database/sql"
	// dbQuery "kpi-bot/db"

	"kpi-bot/lib/bot"
	"fmt"
	"log"

	"kpi-bot/lib/rd"
	// "kpi-bot/lib/test"

	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	// Define the data source name (DSN) Mh-mJ?sp.G"43*_HrCXRP9+^QS%3Et2yZE
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "kpi", "DMdhKuWmUSHevcFBLpxC6R", "192.168.2.8", "30945", "zentao")

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

	// Calculate the first and last day of the previous month
	now := time.Now()
	fmt.Printf("Current time: %s\n", now.Format("2006-01-02 15:04:05"))
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	beginDatetime := firstOfMonth.AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
	endDatetime := firstOfMonth.Add(-time.Second).Format("2006-01-02 15:04:05")
	
	// 手动控制时间
	// beginDatetime := "2025-05-01 00:00:00"
	// endDatetime := "2025-05-31 23:59:59"
	
	fmt.Println("beginDatetime:", beginDatetime)
	fmt.Println("endDatetime:", endDatetime)

	// Initialize the bot with the database connection
	// bot only for old kpi calculation
	robot := bot.NewBot(db)

	
	// // other 研发
	// rdsWithoutTest := []string{
	// 	// embed system
	// 	"qihongquan",
	// 	"zhangzhilu",
	// 	"zhuangjianyong",
	// 	"wangxianming",
	// }
	
	// rdsWithoutTestProjectPms := []string{
	// 	"guoqiao.chen",
	// 	"shawn.wang",
	// 	"simon.chen",
	// 	"qixiaofeng",
	// 	"set.su",
	// 	"justin.lee",
	// 	"caojianni",
	// }
	// err = robot.ProduceRdKpiWithoutTestReport("./excel/kpi-rd-without.xlsx", beginDatetime, endDatetime, rdsWithoutTest, rdsWithoutTestProjectPms)
	// if err != nil {
	// 	log.Fatalf("error produceRdKpiWithoutTestreport: %v", err)
	// }

	// // 项目经理
	// pms := []string{
	// 	"qixiaofeng",
	// 	"jiangjiahui",
	// 	"caojianni",
	// 	"shawn.wang",
	// 	"simon.chen",
	// }
	// err = robot.ProducePmKpi("./excel/pmnew.xlsx", beginDatetime, endDatetime, pms)
	// if err != nil {
	// 	log.Fatalf("error ProducePmKpi: %v", err)
	// }

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
		tmp := rd.NewRdKpi2(db, v, beginDatetime, endDatetime, rd.RdCoefficient{
			PROJECT_PROGRESS_STANDARD: 40,
			DELAY_DAYS_SCORE: 4,
			STORY_STANDARD: 30,
			STORY_BASE_TIME: 0.1,
			STORY_BASE_SCORE: 0.025,
			BUG_CARRY_OVER_STANDARD: 30,
			BUG_ONE_SCORE: 2,
			TOP_COEFFICIENT: 1.2,
			SECOND_COEFFICIENT: 1.0,
			THIRD_COEFFICIENT: 0.8,
		})
		err := tmp.MakeRdReport("软件服务中心", "研发工程师", "ssc-rd", "Set", "./excel/kpi-rd2-2.xlsx")
		if err != nil {
			log.Fatalf("error MakeRdReport: %v", err)
		}
	}

	// SSC TEST
	// test2 := []string{
	// 	"linyanhai",
	// 	"wangshaoyu",
	// 	"pengzijie",
	// 	"xiezhiren",
	// }
	// for _, v := range test2 {
	// 	tmp := test.NewTestKpi2(db, v, beginDatetime, endDatetime, test.TestCoefficient{
	// 		TEST_PROGRESS_STANDARD: 40,
	// 		DELAY_DAYS_SCORE: 4,
	// 		VALIDATE_BUG_RATE_STANDARD: 40,
	// 		BUG_TO_STORY_NUM_STANDARD: 20,
	// 		BUG_ONE_GRADE: 2,
	// 		TOP_COEFFICIENT: 1.2,
	// 		SECOND_COEFFICIENT: 1.0,
	// 		THIRD_COEFFICIENT: 0.8,
	// 	})
	// 	err := tmp.MakeTestReport("软件服务中心", "测试工程师", "ssc-test", "Set", "./excel/kpi-test2.xlsx")
	// 	if err != nil {
	// 		log.Fatalf("error MakeTestReport: %v", err)
	// 	}
	// }

	// APP RD 阿崔部门 研发
	app := []string{
		"alan.tin",
		"jihuaqing",
		"liuhongtao",
		"yuanpengfei",
		"zengyi",
		"chenbo",
		"lixiaolong",
		"jiaoxiangjie",
		"bieji",
		"chenqi",
		"yuanhenghui",
		"liusang",
	}
	for _, v := range app {
		tmp := rd.NewRdKpi2(db, v, beginDatetime, endDatetime, rd.RdCoefficient{
			PROJECT_PROGRESS_STANDARD: 30,
			DELAY_DAYS_SCORE: 2,
			STORY_STANDARD: 30,
			STORY_BASE_TIME: 0.1,
			STORY_BASE_SCORE: 0.025,
			BUG_CARRY_OVER_STANDARD: 30,
			BUG_ONE_SCORE: 0.4,
			BUG_ONE_SCORE_SEVERITY: 0.6,
			TOP_COEFFICIENT: 1.2,
			SECOND_COEFFICIENT: 1.0,
			THIRD_COEFFICIENT: 0.8,
		})
		err := tmp.MakeAppRdReport("APP开发部", "研发工程师", "app-rd", "曾宪崔", "./excel/appkpi.xlsx")
		if err != nil {
			log.Fatalf("error MakeAppRdReport: %v", err)
		}
	}
	// for _, v := range app {
	// 	tmp := rd.NewAppRdKpi(db, v, beginDatetime, endDatetime, rd.AppRdCoefficient{
	// 		PROJECT_PROGRESS_STANDARD: 30,
	// 		DELAY_DAYS_SCORE: 2,
	// 		STORY_STANDARD: 35,
	// 		STORY_BASE_TIME: 0.1,
	// 		STORY_BASE_SCORE: 0.025,
	// 		BUG_CARRY_OVER_STANDARD: 25,
			
	// 		FIRST_BUG_RATE: 0,
	// 		SECOND_BUG_RATE: 0.1,
	// 		THIRD_BUG_RATE: 0.2,
	// 		FORTH_BUG_RATE: 0.3,
	// 		FIFTH_BUG_RATE: 0.4,
	// 		SIXTH_BUG_RATE: 0.5,
			
	// 		FIRST_BUG_BASE: 1.0,
	// 		SECOND_BUG_BASE: 0.9,
	// 		THIRD_BUG_BASE: 0.8,
	// 		FORTH_BUG_BASE: 0.7,
	// 		FIFTH_BUG_BASE: 0.6,
	// 		SIXTH_BUG_BASE: 0.5,
	// 		ZERO_BUG_BASE: 0,
	// 	})
	// 	err := tmp.MakeAppRdReport("APP开发部", "研发工程师", "app-rd", "曾宪崔", "./excel/appkpi.xlsx")
	// 	if err != nil {
	// 		log.Fatalf("error MakeRdReport: %v", err)
	// 	}
	// }

	// 软件6部 方世文
	embed6 := []string{
		"shiwen.tin",
		"wangtuhe",
		"chenyuanchong",
	}
	for _, v := range embed6 {
		tmp := rd.NewRdKpi2(db, v, beginDatetime, endDatetime, rd.RdCoefficient{
			PROJECT_PROGRESS_STANDARD: 40,
			DELAY_DAYS_SCORE: 4,
			STORY_STANDARD: 30,
			STORY_BASE_TIME: 0.1,
			STORY_BASE_SCORE: 0.025,
			BUG_CARRY_OVER_STANDARD: 30,
			BUG_ONE_SCORE: 2,
			TOP_COEFFICIENT: 1.2,
			SECOND_COEFFICIENT: 1.0,
			THIRD_COEFFICIENT: 0.8,
		})
		err := tmp.MakeRdReport("软件6部", "研发工程师", "embed6-rd", "方世文", "./excel/kpi-rd2-2.xlsx")
		if err != nil {
			log.Fatalf("error MakeRdReport: %v", err)
		}
	}

}
package bot

import (
	"database/sql"
	"fmt"
	"kpi-bot/lib/deveops"
	"kpi-bot/lib/excel"
	"kpi-bot/lib/pm"
	"kpi-bot/lib/rd"
	"kpi-bot/lib/test"
)

type (
	Bot struct {
		Db *sql.DB
	}
)


func NewBot(db *sql.DB) *Bot {
	return &Bot{
		Db: db,
	}
}



func (l *Bot) ProduceRdKpi(templatePath, startTime, endTime string, accounts []string) error {
	kpiManager := rd.NewRdKpi(l.Db, accounts, startTime, endTime)
	kpiGrades := kpiManager.GetRdKpiGrade()

	for _, kpiGrade := range kpiGrades {
		err := excel.MakeRdExcel(templatePath, kpiGrade)
		if err != nil {
			return fmt.Errorf("make rd excel fail: %v", err)
		}
	}
	return nil
}


func (l *Bot) ProduceRdKpiWithoutTestReport(templatePath, startTime, endTime string, accounts []string) error {
	kpiManager := rd.NewRdKpiWithoutTestReport(l.Db, accounts, startTime, endTime)
	kpiGrades := kpiManager.GetRdKpiWithoutTestReportGrade()

	for _, kpiGrade := range kpiGrades {
		err := excel.MakeRdWithoutTestreportExcel(templatePath, kpiGrade)
		if err != nil {
			return fmt.Errorf("make rd without testreport excel fail: %v", err)
		}
	}
	return nil
}

func (l *Bot) ProducePmKpi(templatePath, startTime, endTime string, accounts []string) error {
	kpiManager := pm.NewPmKpi(l.Db, accounts, startTime, endTime)
	kpiGrades := kpiManager.GetPmKpiGrade()

	for _, kpiGrade := range kpiGrades {
		err := excel.MakePmExcel(templatePath, kpiGrade)
		if err != nil {
			return fmt.Errorf("make pm excel fail: %v", err)
		}
	}
	return nil
}

func (l *Bot) ProducePmKpiWithoutTestReport(templatePath, startTime, endTime string, accounts []string) error {
	kpiManager := pm.NewPmKpiWithoutTestReport(l.Db, accounts, startTime, endTime)
	kpiGrades := kpiManager.GetPmKpiGradeWithoutTestReport()

	for _, kpiGrade := range kpiGrades {
		err := excel.MakePmExcelWithoutTestReport(templatePath, kpiGrade)
		if err != nil {
			return fmt.Errorf("make pm excel fail: %v", err)
		}
	}
	return nil
}

func (l *Bot) ProduceTestKpi(templatePath, startTime, endTime string, accounts []string) error {
	kpiManager := test.NewTestKpi(l.Db, accounts, startTime, endTime)
	kpiGrades := kpiManager.GetTestKpiGrade()

	for _, kpiGrade := range kpiGrades {
		err := excel.MakeTestExcel(templatePath, kpiGrade)
		if err != nil {
			return fmt.Errorf("make test excel fail: %v", err)
		}
	}
	return nil
}



func (l *Bot) ProduceStatisticKpi(tmplatePath, startTime, endTime string, rds, rdWithouts, pms, pmWithouts, tests []string) error {
	rdkpiManager := rd.NewRdKpi(l.Db, rds, startTime, endTime)
	rdWithoutKpiManager := rd.NewRdKpiWithoutTestReport(l.Db, rdWithouts, startTime, endTime)
	pmKpiManager := pm.NewPmKpi(l.Db, pms, startTime, endTime)
	pmWithoutKpiManager := pm.NewPmKpiWithoutTestReport(l.Db, pmWithouts, startTime, endTime)
	testKpiManager := test.NewTestKpi(l.Db, tests, startTime, endTime)


	rdKpiGrades := rdkpiManager.GetRdKpiGrade()
	rdWithoutKpiGrades := rdWithoutKpiManager.GetRdKpiWithoutTestReportGrade()
	pmKpiGrades := pmKpiManager.GetPmKpiGrade()
	pmWithoutKpiGrades := pmWithoutKpiManager.GetPmKpiGradeWithoutTestReport()
	testKpiGrades := testKpiManager.GetTestKpiGrade()


	err := excel.MakeKpiStatisticsExcel(tmplatePath, startTime, rdKpiGrades, rdWithoutKpiGrades, pmKpiGrades, pmWithoutKpiGrades, testKpiGrades)
	if err != nil {
		return fmt.Errorf("make kpi statistics excel fail: %v", err)
	}
	return nil
}

func (l *Bot) ProduceDeveopsKpi(templatePath, startTime, endTime string, accounts []string) error {
	kpiManager := deveops.NewDeveopsKpi(l.Db, accounts, startTime, endTime)
	kpiGrades := kpiManager.GetDeveopsKpiGrade()

	for _, kpiGrade := range kpiGrades {
		err := excel.MakeDeveopsExcel(templatePath, kpiGrade)
		if err != nil {
			return fmt.Errorf("make rd without testreport excel fail: %v", err)
		}
	}
	return nil
}
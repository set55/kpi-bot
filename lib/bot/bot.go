package bot

import (
	"database/sql"
	"fmt"
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
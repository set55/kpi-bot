package test

import (
	"database/sql"
	"fmt"
	"kpi-bot/common"
	dbQuery "kpi-bot/db"
	"os"
	"slices"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	// 测试软件项目进度达成率 分值
	TEST_PROGRESS_STANDARD2 = 40
	DELAY_DAYS_SCORE = 4

	// 测试软件项目有效bug率
	VALIDATE_BUG_RATE_STANDARD2 = 40

	// bug转需求数
	BUG_TO_STORY_NUM_STANDARD2 = 20
	BUG_ONE_GRADE2             = 2

	// 系数
	TOP_COEFFICIENT2    = 1.2
	SECOND_COEFFICIENT2 = 1.0
	THIRD_COEFFICIENT2  = 0.8
)

type (
	TestKpi2 struct {
		Account string
		Db      *sql.DB
		// 起始时间
		StartTime string
		// 结束时间
		EndTime string
	}

	TestKpiResult struct {
		ReportDetail string
		ReportGrade  float64
		BugDetail    string
		BugGrade     float64
		ToStoryDetail string
		ToStoryGrade  float64
		TotalGrade    float64
		StartTime     string
		EndTime       string
		AccountName string
		Coefficient float64
	}

)

// NewTestKpi 创建一个测试KPI对象
func NewTestKpi2(db *sql.DB, account, startTime, endTime string) *TestKpi2 {
	return &TestKpi2{
		Account:  account,
		Db:        db,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// GetTestKpiGrade 获取测试KPI信息
func (l *TestKpi2) GetTestKpiGrade() (result TestKpiResult) {
	result.AccountName = common.AccountToName(l.Account)
	result.StartTime = l.StartTime
	result.EndTime = l.EndTime
	
	// 获取测试报告
	testreports := dbQuery.QueryTestReport(l.Db, l.Account, l.StartTime, l.EndTime)
	delaydays := 0

	for _, testreport := range testreports {
		singledelaydays := common.CalculateDelayDays(testreport.ReportCreatedAt, testreport.TaskEnd)
		delaydays += singledelaydays
		result.ReportDetail += fmt.Sprintf("測試單id: %d, 測試報告id: %d, 測試單名称: %s, 測試單结束时间: %s, 测试报告创建时间: %s, 延遲天數: %d\n\n",
		 testreport.TaskId, testreport.ReportId, testreport.TaskName, testreport.TaskEnd, testreport.ReportCreatedAt, singledelaydays)
	}

	// 
	result.ReportGrade = float64(TEST_PROGRESS_STANDARD2 - delaydays * DELAY_DAYS_SCORE)
	if result.ReportGrade < 0 {
		result.ReportGrade = 0
	}

	// 获取有效bug率
	bugs := dbQuery.TestBugs(l.Db, l.Account, l.StartTime, l.EndTime)
	validateBugs := 0
	toStoryBugs := 0
	for _, bug := range bugs {
		if l.IsValidateBug(bug) {
			validateBugs++
		}
		if bug.BugResolution == "tostory" {
			toStoryBugs++
		}
		result.BugDetail += fmt.Sprintf("bug id: %d, bug标题: %s, bug创建人: %s, bug状态: %s, bug解决方案: %s\n\n", bug.BugId, bug.BugTitle, bug.BugCreator, bug.BugStatus, bug.BugResolution)
	}

	// 有效bug率
	bugRate := float64(validateBugs) / float64(len(bugs))
	result.BugGrade = l.ConverBugRateToBaseNumber(bugRate) * VALIDATE_BUG_RATE_STANDARD2
	result.BugDetail = fmt.Sprintf("有效bug数: %d, 总bug数: %d, 有效bug率: %.2f\n\n", validateBugs, len(bugs), bugRate) + result.BugDetail

	// bug转需求数
	result.ToStoryDetail = fmt.Sprintf("bug转需求数: %d\n\n", toStoryBugs)
	result.ToStoryGrade = float64(toStoryBugs) * BUG_ONE_GRADE2

	// 计算总分
	result.TotalGrade = result.ReportGrade + result.BugGrade + result.ToStoryGrade
	result.Coefficient = l.GetKpiGradeStandard(result.TotalGrade)
	return
}

// 计算得分系数
func (l *TestKpi2) GetKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 90 {
		return TOP_COEFFICIENT
	} else if totalGrade < 90 && totalGrade >= 70 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 60 && totalGrade >= 70 {
		return THIRD_COEFFICIENT
	}
	return 0
}

func (l *TestKpi2) IsValidateBug(bug dbQuery.TestBug) bool {
	if slices.Contains([]string{"", "fixed", "postponed", "tostory", "willnotfix"}, bug.BugResolution) {
		return true
	}
	return false
}

func (l *TestKpi2) ConverBugRateToBaseNumber(bugRate float64) float64 {
	if bugRate >= 0.9 {
		return 1.0
	} else if bugRate < 0.9 && bugRate >= 0.8 {
		return 0.8
	} else if bugRate < 0.8 && bugRate >= 0.7 {
		return 0.5
	} else {
		return 0.0
	}
}

func (l *TestKpi2) MakeTestReport(path string) error {
	data := l.GetTestKpiGrade()

	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("open file fail: %v", err)
	}
	defer f.Close()

	// Parse the date string into a time.Time object
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, data.StartTime)
	if err != nil {
		return fmt.Errorf("parse time err: %v", err)
	}

	// Extract the year and month
	year := t.Year()
	month := t.Month()

	// A1. 标题
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 測試工程师岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", data.AccountName))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度延时率 完成情况
	
	f.SetCellValue("Sheet1", "G4", data.ReportDetail)

	// H4. 项目进度延时率 最终得分
	f.SetCellValue("Sheet1", "H4", data.ReportGrade)

	// G5. 有效bug率  完成情况
	f.SetCellValue("Sheet1", "G5", data.BugDetail)

	// H5. 有效bug率  最终得分
	f.SetCellValue("Sheet1", "H5", data.BugGrade)

	// G6. bug转需求数 完成情况
	f.SetCellValue("Sheet1", "G6", data.ToStoryDetail)

	// H6. bug转需求数 最终得分
	f.SetCellValue("Sheet1", "H6", data.ToStoryGrade)

	// H9. 总分数
	f.SetCellValue("Sheet1", "H9", data.TotalGrade)

	// G11. 绩效基数
	f.SetCellValue("Sheet1", "G11", data.Coefficient)

	// G14. 最终得分系数

	// G17. 当月绩效奖金

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板-測試-%v.xlsx", year, int(month), year, int(month), data.AccountName)
	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Create the folder if it does not exist
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating folder: %v", err)
		}
		fmt.Println("Folder created successfully.")
	} else {
		fmt.Println("Folder already exists.")
	}

	// Save the modified file
	if err = f.SaveAs(filePath); err != nil {
		return fmt.Errorf("save as %v, err: %v", filePath, err)
	}
	return nil
}

package test

import (
	"database/sql"
	"fmt"
	"kpi-bot/common"
	dbQuery "kpi-bot/db"
	"kpi-bot/lib/rd"
)

const (
	// 测试软件项目进度达成率 分值
	TEST_PROGRESS_STANDARD = 30

	// 测试软件项目有效bug率
	VALIDATE_BUG_RATE_STANDARD = 30

	// bug转需求数
	BUG_TO_STORY_NUM_STANDARD = 20
	BUG_ONE_GRADE             = 3

	// 用例发现bug率
	CASE_BUG_RATE_STANDARD = 20

	// 系数
	TOP_COEFFICIENT    = 1.2
	SECOND_COEFFICIENT = 1.0
	THIRD_COEFFICIENT  = 0.8
)

type (
	TestKpi struct {
		Accounts  []string // test的账号
		Db        *sql.DB  // 数据库连接
		StartTime string   // 开始时间
		EndTime   string   // 结束时间
		StoryPms  []string // 需求PM
	}

	TestKpiGrade struct {
		Account string // 禅道账号

		StartTime string // 开始时间
		EndTime   string // 结束时间

		ProjectTotalSaturdays float64
		ProjectTotalSundays   float64

		RealTotalSaturdays float64
		RealTotalSundays   float64

		// 测试软件项目进度 完成情况
		TestProgressInfos []TestProgressDetail

		// 测试软件项目进度达成率
		SumRealTestDiffDays                  float64 // 總實際測試天數
		SumTestDiffDays                      float64 // 總預估測試天數
		DiffRate                             float64 // 平均项目进度延时率
		TestProgressAvgDiffDaysStandard      float64 // 项目测试进度达成基数
		TestProgressAvgDiffDaysStandardGrade float64 // 项目测试进度达成率 实际分数

		// 测试软件项目有效bug率
		// 1、测试报告结束时间是当月的
		// 2、bug未被删除，bug关联项目属于测试报告关联项目，bug关联版本是测试报告所属版本，bug是焰海打开的，bug解决状态是转需求，延期处理和已解决的，不予解决，这些叫有效bug。
		// 3、版本内所有bug，为项目与测试报告相等，并且不是指派给黄卫旗
		ValidateBugRate              float64 // 有效bug率
		ValidateBugRateStandard      float64 // 有效bug率基数
		ValidateBugRateStandardGrade float64 // 有效bug率实际分数

		// bug转需求数
		BugToStoryNum   int     // bug转需求数
		BugToStoryGrade float64 // bug转需求数实际分数

		// 用例发现bug率
		CaseBugRate              float64 // 用例发现bug率
		CaseBugRateStandard      float64 // 用例发现bug率基数
		CaseBugRateStandardGrade float64 // 用例发现bug率实际分数

		TotalGrade         float64 // 总分数
		TotalGradeStandard float64 // 总分数基数
	}

	TestProgressDetail struct {
		Account         string // 禅道账号
		TestTaskName    string // 测试任务名称
		TestReportTitle string // 测试报告标题
		TestTaskBegin   string // 测试任务开始时间
		TestTaskEnd     string // 测试任务预估结束时间
		TestReportEnd   string // 测试报告实际结束时间
	}
)

// NewTestKpi 创建一个测试KPI对象
func NewTestKpi(db *sql.DB, accounts, storyPms []string, startTime, endTime string) *TestKpi {
	return &TestKpi{
		Accounts:  accounts,
		Db:        db,
		StartTime: startTime,
		EndTime:   endTime,
		StoryPms:  storyPms,
	}
}

// GetTestKpiGrade 获取测试KPI信息
func (l *TestKpi) GetTestKpiGrade() map[string]TestKpiGrade {
	kpiGrades := make(map[string]TestKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = TestKpiGrade{
			Account:   account,
			StartTime: l.StartTime,
			EndTime:   l.EndTime,
		}
	}

	// 测试软件项目进度 完成情况
	fmt.Print("测试软件项目进度 完成情况\n")
	testProgressInfos := dbQuery.QueryTestProjectProgressResultDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)

	for account, result := range testProgressInfos {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				planSaturdays, planSundays := common.CountWeekends(r.TestTaskBegin, r.TestTaskEnd)
				realSaturdays, realSundays := common.CountWeekends(r.TestTaskBegin, r.TestReportEnd)
				tmp.TestProgressInfos = append(tmp.TestProgressInfos, TestProgressDetail{
					Account:         r.Account,
					TestTaskName:    r.TestTaskName,
					TestReportTitle: r.TestReportTitle,
					TestTaskBegin:   r.TestTaskBegin,
					TestTaskEnd:     r.TestTaskEnd,
					TestReportEnd:   r.TestReportEnd,
				})
				tmp.ProjectTotalSaturdays += float64(planSaturdays) / 2
				tmp.ProjectTotalSundays += float64(planSundays)

				tmp.RealTotalSaturdays += float64(realSaturdays) / 2
				tmp.RealTotalSundays += float64(realSundays)
				tmp.ProjectTotalSaturdays = 0
				tmp.ProjectTotalSundays = 0
				tmp.RealTotalSaturdays = 0
				tmp.RealTotalSundays = 0
			}
			kpiGrades[account] = tmp
		}
	}


	// 测试软件项目进度达成率
	fmt.Print("测试软件项目进度达成率\n")
	testProgressResult := dbQuery.QueryTestProjectProgress(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range testProgressResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.SumRealTestDiffDays = result.SumRealTestDiffDays - tmp.RealTotalSaturdays - tmp.RealTotalSundays
			tmp.SumTestDiffDays = result.SumTestDiffDays - tmp.ProjectTotalSaturdays - tmp.ProjectTotalSundays
			tmp.DiffRate = common.GetProjectProgressExpectRate(tmp.SumTestDiffDays, tmp.SumRealTestDiffDays)
			tmp.TestProgressAvgDiffDaysStandard = rd.GetRdProjectProgressStandard(tmp.DiffRate)
			tmp.TestProgressAvgDiffDaysStandardGrade = tmp.TestProgressAvgDiffDaysStandard * TEST_PROGRESS_STANDARD
			tmp.TotalGrade += tmp.TestProgressAvgDiffDaysStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 测试软件项目有效bug率
	fmt.Print("测试软件项目有效bug率\n")
	validateBugRateResult := dbQuery.QueryTestValidBugRate(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range validateBugRateResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.ValidateBugRate = result.ValidBugRate
			tmp.ValidateBugRateStandard = result.ValidBugRateStandard
			tmp.ValidateBugRateStandardGrade = result.ValidBugRateStandard * VALIDATE_BUG_RATE_STANDARD
			tmp.TotalGrade += tmp.ValidateBugRateStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// bug转需求数
	fmt.Print("bug转需求数\n")
	bugToStoryNumResult := dbQuery.QueryTestBugToStory(l.Db, l.Accounts, l.StoryPms, l.StartTime, l.EndTime)
	for account, result := range bugToStoryNumResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.BugToStoryNum = result.ToStory
			tmp.BugToStoryGrade = float64(result.ToStory) * BUG_ONE_GRADE
			if tmp.BugToStoryGrade > BUG_TO_STORY_NUM_STANDARD {
				tmp.BugToStoryGrade = BUG_TO_STORY_NUM_STANDARD
			}
			tmp.TotalGrade += tmp.BugToStoryGrade
			kpiGrades[account] = tmp
		}
	}

	// 用例发现bug率
	fmt.Print("用例发现bug率\n")
	caseBugRateResult := dbQuery.QueryTestBugCaseRate(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range caseBugRateResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.CaseBugRate = result.CaseBugRate
			tmp.CaseBugRateStandard = result.CaseBugStandard
			tmp.CaseBugRateStandardGrade = result.CaseBugStandard * CASE_BUG_RATE_STANDARD
			tmp.TotalGrade += tmp.CaseBugRateStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 结算系数
	for account, kpiGrade := range kpiGrades {
		tmp := kpiGrades[account]
		tmp.TotalGradeStandard = l.GetKpiGradeStandard(kpiGrade.TotalGrade)
		kpiGrades[account] = tmp
	}

	return kpiGrades
}

// 计算得分系数
func (l *TestKpi) GetKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 90 {
		return TOP_COEFFICIENT
	} else if totalGrade < 90 && totalGrade >= 70 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 60 && totalGrade >= 70 {
		return THIRD_COEFFICIENT
	}
	return 0
}

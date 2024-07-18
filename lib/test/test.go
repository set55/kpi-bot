package test

import (
	"database/sql"
	dbQuery "kpi-bot/db"
)

const (
	// 测试软件项目进度达成率 分值
	TEST_PROGRESS_STANDARD = 30

	// 测试软件项目有效bug率
	VALIDATE_BUG_RATE_STANDARD = 30

	// bug转需求数
	BUG_TO_STORY_NUM_STANDARD = 20
	BUG_ONE_GRADE = 3

	// 用例发现bug率
	CASE_BUG_RATE_STANDARD = 20

	// 系数
	TOP_COEFFICIENT = 1.2
	SECOND_COEFFICIENT = 1.0
	THIRD_COEFFICIENT = 0.7


)

type (
	TestKpi struct {
		Accounts []string // test的账号
		Db       *sql.DB  // 数据库连接
		StartTime string  // 开始时间
		EndTime string // 结束时间
	}

	TestKpiGrade struct {
		Account string // 禅道账号

		StartTime string // 开始时间
		EndTime   string // 结束时间

		// 测试软件项目进度达成率
		TestProgressAvgDiffDays              float64 // 平均项目测试进度预估天数差值
		TestProgressAvgDiffDaysStandard      float64 // 项目测试进度达成基数
		TestProgressAvgDiffDaysStandardGrade float64 // 项目测试进度达成率 实际分数

		// 测试软件项目有效bug率
		// 1、测试报告结束时间是当月的
		// 2、bug未被删除，bug关联项目属于测试报告关联项目，bug关联版本是测试报告所属版本，bug是焰海打开的，bug解决状态是转需求，延期处理和已解决的，不予解决，这些叫有效bug。
		// 3、版本内所有bug，为项目与测试报告相等，并且不是指派给黄卫旗
		ValidateBugRate float64 // 有效bug率
		ValidateBugRateStandard float64 // 有效bug率基数
		ValidateBugRateStandardGrade float64 // 有效bug率实际分数

		// bug转需求数
		BugToStoryNum int // bug转需求数
		BugToStoryGrade float64 // bug转需求数实际分数

		// 用例发现bug率
		CaseBugRate float64 // 用例发现bug率
		CaseBugRateStandard float64 // 用例发现bug率基数
		CaseBugRateStandardGrade float64 // 用例发现bug率实际分数

		TotalGrade float64 // 总分数
		TotalGradeStandard float64 // 总分数基数
	}
)



// NewTestKpi 创建一个测试KPI对象
func NewTestKpi(db *sql.DB, accounts []string, startTime, endTime string) *TestKpi {
	return &TestKpi{
		Accounts: accounts,
		Db:       db,
		StartTime: startTime,
		EndTime: endTime,
	}
}

// GetTestKpiGrade 获取测试KPI信息
func (l *TestKpi) GetTestKpiGrade() map[string]TestKpiGrade {
	kpiGrades := make(map[string]TestKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = TestKpiGrade{
			Account: account,
			StartTime: l.StartTime,
			EndTime: l.EndTime,
		}
	}


	// 测试软件项目进度达成率
	testProgressResult := dbQuery.QueryTestProjectProgress(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range testProgressResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.TestProgressAvgDiffDays = result.AvgDiffDays
			tmp.TestProgressAvgDiffDaysStandard = result.DiffDaysStandard
			tmp.TestProgressAvgDiffDaysStandardGrade = result.DiffDaysStandard * TEST_PROGRESS_STANDARD
			tmp.TotalGrade += tmp.TestProgressAvgDiffDaysStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 测试软件项目有效bug率
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
	bugToStoryNumResult := dbQuery.QueryTestBugToStory(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range bugToStoryNumResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.BugToStoryNum = result.ToStory
			tmp.BugToStoryGrade = float64(result.ToStory) * BUG_ONE_GRADE
			if tmp.BugToStoryGrade > BUG_TO_STORY_NUM_STANDARD {
				tmp.TotalGrade += BUG_TO_STORY_NUM_STANDARD
			} else {
				tmp.TotalGrade += tmp.BugToStoryGrade
			}
			kpiGrades[account] = tmp
		}
	}


	// 用例发现bug率
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
		tmp.TotalGradeStandard = l.GetRdKpiGradeStandard(kpiGrade.TotalGrade)
		kpiGrades[account] = tmp
	}



	return kpiGrades
}

// 计算得分系数
func (l *TestKpi) GetRdKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 100 {
		return TOP_COEFFICIENT
	} else if totalGrade < 100 && totalGrade >= 80 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 80 && totalGrade >= 60 {
		return THIRD_COEFFICIENT
	}
	return 0
}

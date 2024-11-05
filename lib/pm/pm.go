package pm

import (
	"database/sql"
	"fmt"
	"kpi-bot/common"
	dbQuery "kpi-bot/db"
)

const (

	// 项目进度延时率基数
	PM_PROJECT_PROGRESS_Level1 = 0.8
	PM_PROJECT_PROGRESS_Level2 = 1
	PM_PROJECT_PROGRESS_Level3 = 1.2
	PM_PROJECT_PROGRESS_Level4 = 0.8
	PM_PROJECT_PROGRESS_Level5 = 0.5

	// 项目软件项目进度达成率 分值
	PROJECT_PROGRESS_STANDARD = 30

	// 项目成果完成率 分值
	PROJECT_COMPLETEMENT_STANDARD = 20

	// 项目规划需求数 分值
	PROJECT_STORY_NUM_STANDARD = 20
	PROJECTED_STORY_STANDARD   = 0.5
	DEVELOPED_STORY_STANDARD   = 1
	CLOSED_STORY_STANDARD      = 1

	// 预估承诺完成率 分值
	PROJECT_ESTIMATE_STANDARD = 0

	// 项目预估工时准确率 分值
	TIME_ESTIMATE_STANDARD = 30

	// 系数
	TOP_COEFFICIENT    = 1.2
	SECOND_COEFFICIENT = 1.0
	THIRD_COEFFICIENT  = 0.7
)

type (
	PmKpi struct {
		Accounts  []string // pm的账号
		Db        *sql.DB  // 数据库连接
		StartTime string   // 开始时间
		EndTime   string   // 结束时间
	}

	PmKpiGrade struct {
		Account string // 禅道账号

		StartTime string // 开始时间
		EndTime   string // 结束时间

		ProjectTotalSaturdays float64
		ProjectTotalSundays   float64

		RealTotalSaturdays float64
		RealTotalSundays   float64

		// 项目软件项目进度达成率
		SumRealProjectDiff    float64 // 项目实际结束天數
		SumProjectDiff        float64 // 项目预估结束天数
		SumRealTestDiff       float64 // 测试实际结束天数
		SumTestDiff           float64 // 测试预估结束天数
		DiffRate              float64 // 平均项目进度延时率
		ProgressStandard      float64 // 项目进度达成基数
		ProgressStandardGrade float64 // 项目进度达成率 实际分数

		// 项目软件项目进度达成率 完成情况
		ProjectProgressList []ProjectProgressInfo

		// 项目成果完成率
		CompleteRate              float64 // 项目成果完成率
		CompleteRateStandard      float64 // 项目成果完成率基数
		CompleteRateStandardGrade float64 // 项目成果完成率实际分数

		// 项目成果完成率,完成情况
		ProjectCompleteList []ProjectCompleteInfo

		// 项目规划需求数
		ProjectedStoryNum int     // 评审完的需求数
		DevelopedStoryNum int     // 开发完的需求数
		ClosedStoryNum    int     // 关闭的需求数
		StoryNumGrade     float64 // 需求数实际分数

		// 预估承诺完成率
		PromiseDiffDays      float64 // 预估承诺完成率
		PromiseStandard      float64 // 预估承诺完成率基数
		PromiseStandardGrade float64 // 预估承诺完成率实际分数

		// 项目预估工时准确率
		TimeEstimateRate     float64 // 项目预估工时准确率
		TimeEstimateStandard float64 // 项目预估工时准确率基数
		TimeEstimateGrade    float64 // 项目预估工时准确率实际分数

		// 项目预估工时准确率 完成情况
		TimeEstimateList []TimeEstimateInfo

		TotalGrade         float64 // 总分数
		TotalGradeStandard float64 // 总分数基数

	}

	ProjectProgressInfo struct {
		Account        string // 禅道账号
		ProjectId      int64  // 项目id
		ProjectName    string // 项目名称
		ProjectType    string // 项目类型
		ProjectBegin   string // 项目开始时间
		ProjectEnd     string // 项目预估结束时间
		ProjectRealEnd string // 项目实际结束时间
		ProjectDiff    int    // 与预期相差天数
		TestStart      string // 测试开始时间
		TestEnd        string // 测试預估结束时间
		TestRealEnd    string // 测试实际结束时间
	}

	ProjectCompleteInfo struct {
		ProjectName  string  // 项目名称
		CompleteRate float64 // 项目成果完成率
	}

	TimeEstimateInfo struct {
		Account       string  // 禅道账号
		StoryId       int64   // 需求id
		Title         string  // 需求标题
		Estimate      float64 // 预估工时
		StoryConsumed float64 // 需求消耗工时
		EstimateRate  float64 // 预估工时准确率
	}
)

// NewPmKpi creates a new PmKpi object
func NewPmKpi(db *sql.DB, accounts []string, startTime, endTime string) *PmKpi {
	return &PmKpi{
		Accounts:  accounts,
		Db:        db,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// GetPmKpiGrade gets the PM KPI information
func (l *PmKpi) GetPmKpiGrade() map[string]PmKpiGrade {
	kpiGrades := make(map[string]PmKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = PmKpiGrade{
			Account:   account,
			StartTime: l.StartTime,
			EndTime:   l.EndTime,
		}
	}

	// 项目软件项目进度达成率 完成情况
	fmt.Print("项目软件项目进度达成率 完成情况\n")
	progressDetailResult := dbQuery.QueryProjectProgressDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range progressDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				planSaturdays, planSundays := common.CountWeekends(r.ProjectBegin, r.ProjectEnd)
				realSaturdays, realSundays := common.CountWeekends(r.ProjectBegin, r.ProjectRealEnd)
				fmt.Printf("projectBegin: %v, ProjectEnd: %v, ProjectRealEnd: %v, planSaturdays: %v, planSundays: %v, realSaturdays: %v, realSundays: %v\n", r.ProjectBegin, r.ProjectEnd, r.ProjectRealEnd, planSaturdays, planSundays, realSaturdays, realSundays)

				testPlanSaturdays, testPlanSundays := common.CountWeekends(r.TestStart, r.TestEnd)
				testRealSaturdays, testRealSundays := common.CountWeekends(r.TestStart, r.TestRealEnd)
				tmp.ProjectProgressList = append(tmp.ProjectProgressList, ProjectProgressInfo{
					Account:        r.Account,
					ProjectId:      r.ProjectId,
					ProjectName:    r.ProjectName,
					ProjectType:    r.ProjectType,
					ProjectBegin:   r.ProjectBegin,
					ProjectEnd:     r.ProjectEnd,
					ProjectRealEnd: r.ProjectRealEnd,
					ProjectDiff:    r.ProjectDiff,
					TestStart:      r.TestStart,
					TestEnd:        r.TestEnd,
					TestRealEnd:    r.TestRealEnd,
				})

				tmp.ProjectTotalSaturdays += float64(planSaturdays + testPlanSaturdays) / 2
				tmp.ProjectTotalSundays += float64(planSundays + testPlanSundays)

				tmp.RealTotalSaturdays += float64(realSaturdays + testRealSaturdays) / 2
				tmp.RealTotalSundays += float64(realSundays + testRealSundays)

				tmp.ProjectTotalSaturdays = 0
				tmp.ProjectTotalSundays = 0
				tmp.RealTotalSaturdays = 0
				tmp.RealTotalSundays = 0
			}
			kpiGrades[account] = tmp
		}
	}


	// 项目软件项目进度达成率
	fmt.Print("项目软件项目进度达成率\n")
	progressResult := dbQuery.QueryProjectProgress(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range progressResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.SumRealProjectDiff = result.SumRealProjectDiff
			tmp.SumProjectDiff = result.SumProjectDiff
			tmp.SumRealTestDiff = result.SumRealTestDiff
			tmp.SumTestDiff = result.SumTestDiff
			fmt.Printf("QueryProjectProgress,account: %v, SumRealProjectDiff: %v, SumProjectDiff: %v, SumRealTestDiff: %v, SumTestDiff: %v, ProjectTotalSaturdays: %v, ProjectTotalSundays: %v, RealTotalSaturdays: %v, RealTotalSundays: %v\n",
			account, tmp.SumRealProjectDiff, tmp.SumProjectDiff, tmp.SumRealTestDiff, tmp.SumTestDiff, tmp.ProjectTotalSaturdays, tmp.ProjectTotalSundays, tmp.RealTotalSaturdays, tmp.RealTotalSundays)
			tmp.DiffRate = common.GetProjectProgressExpectRate((tmp.SumProjectDiff + tmp.SumTestDiff - tmp.ProjectTotalSaturdays - tmp.ProjectTotalSundays), (tmp.SumRealProjectDiff + tmp.SumRealTestDiff - tmp.RealTotalSaturdays - tmp.RealTotalSundays))
			fmt.Printf("QueryProjectProgress,account: %v, DiffRate: %v\n", account, tmp.DiffRate)
			tmp.ProgressStandard = GetPmProjectProgressStandard(tmp.DiffRate)
			tmp.ProgressStandardGrade = tmp.ProgressStandard * PROJECT_PROGRESS_STANDARD
			tmp.TotalGrade += tmp.ProgressStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 项目成果完成率
	fmt.Print("项目成果完成率\n")
	completeRateResult := dbQuery.QueryProjectCompleteRate(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range completeRateResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.CompleteRate = result.CompleteRate
			tmp.CompleteRateStandard = result.CompleteRateStandard
			tmp.CompleteRateStandardGrade = result.CompleteRateStandard * PROJECT_COMPLETEMENT_STANDARD
			tmp.TotalGrade += tmp.CompleteRateStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 项目成果完成率,完成情况
	fmt.Print("项目成果完成率,完成情况\n")
	completeRateDetailResult := dbQuery.QueryProjectCompleteRateDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range completeRateDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmp.ProjectCompleteList = append(tmp.ProjectCompleteList, ProjectCompleteInfo{
					ProjectName:  r.ProjectName,
					CompleteRate: r.CompleteRate,
				})
			}
			kpiGrades[account] = tmp
		}
	}

	// 项目规划需求数
	fmt.Print("项目规划需求数\n")
	storyNumResult := dbQuery.QueryProjectStoryNum(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range storyNumResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				switch r.Stage {
				case "projected":
					tmp.ProjectedStoryNum = r.StoryNum
					tmp.StoryNumGrade += float64(tmp.ProjectedStoryNum) * PROJECTED_STORY_STANDARD
				case "developed":
					tmp.DevelopedStoryNum = r.StoryNum
					tmp.StoryNumGrade += float64(tmp.DevelopedStoryNum) * DEVELOPED_STORY_STANDARD
				case "closed":
					tmp.ClosedStoryNum = r.StoryNum
					tmp.StoryNumGrade += float64(tmp.ClosedStoryNum) * CLOSED_STORY_STANDARD
				}
			}

			if tmp.StoryNumGrade > PROJECT_STORY_NUM_STANDARD {
				tmp.StoryNumGrade = PROJECT_STORY_NUM_STANDARD
			}

			tmp.TotalGrade += tmp.StoryNumGrade
			kpiGrades[account] = tmp
		}
	}

	// 预估承诺完成率
	fmt.Print("预估承诺完成率\n")
	projectPromiseResult := dbQuery.QueryProjectEstimateRate(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range projectPromiseResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.PromiseDiffDays = result.DiffDays
			tmp.PromiseStandard = result.ProgressStandard
			tmp.PromiseStandardGrade = result.ProgressStandard * PROJECT_ESTIMATE_STANDARD
			tmp.TotalGrade += tmp.PromiseStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 项目预估工时准确率
	fmt.Print("项目预估工时准确率\n")
	timeEstimateResult := dbQuery.QueryProjectTimeEstimateRate(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range timeEstimateResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.TimeEstimateRate = result.TimeEstimateRate
			tmp.TimeEstimateStandard = result.TimeEstimateStandard
			tmp.TimeEstimateGrade = tmp.TimeEstimateStandard * TIME_ESTIMATE_STANDARD
			tmp.TotalGrade += tmp.TimeEstimateGrade
			kpiGrades[account] = tmp
		}
	}

	// 项目预估工时准确率 完成情况
	fmt.Print("项目预估工时准确率 完成情况\n")
	timeEstimateDetailResult := dbQuery.QueryProjectTimeEstimateRateDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range timeEstimateDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmpEstimateRate := float64(0)
				if r.EstimateRate != nil {
					tmpEstimateRate = *r.EstimateRate
				}

				tmp.TimeEstimateList = append(tmp.TimeEstimateList, TimeEstimateInfo{
					Account:       r.Account,
					StoryId:       r.StoryId,
					Title:         r.Title,
					Estimate:      r.Estimate,
					StoryConsumed: r.StoryConsumed,
					EstimateRate:  tmpEstimateRate,
				})
			}
			kpiGrades[account] = tmp
		}
	}

	for account, kpiGrade := range kpiGrades {
		tmp := kpiGrades[account]
		tmp.TotalGradeStandard = l.GetRdKpiGradeStandard(kpiGrade.TotalGrade)
		kpiGrades[account] = tmp
	}

	return kpiGrades
}

// 计算得分系数
func (l *PmKpi) GetRdKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 100 {
		return TOP_COEFFICIENT
	} else if totalGrade < 100 && totalGrade >= 80 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 80 && totalGrade >= 60 {
		return THIRD_COEFFICIENT
	}
	return 0
}

func GetPmProjectProgressStandard(avgDiffRate float64) float64 {
	if avgDiffRate <= -0.5 {
		return PM_PROJECT_PROGRESS_Level1
	} else if avgDiffRate > -0.5 && avgDiffRate <= -0.2 {
		return PM_PROJECT_PROGRESS_Level2
	}else if avgDiffRate > -0.2 && avgDiffRate <= 0 {
		return PM_PROJECT_PROGRESS_Level3
	} else if avgDiffRate > 0 && avgDiffRate <= 0.2 {
		return PM_PROJECT_PROGRESS_Level4
	} else if avgDiffRate > 0.2 && avgDiffRate <= 0.5 {
		return PM_PROJECT_PROGRESS_Level5
	} else {
		return 0
	}
}

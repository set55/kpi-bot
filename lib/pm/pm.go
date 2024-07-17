package pm

import (
	"database/sql"
	dbQuery "kpi-bot/db"
)


const (
	// 项目软件项目进度达成率 分值
	PROJECT_PROGRESS_STANDARD = 40

	// 项目成果完成率 分值
	PROJECT_COMPLETEMENT_STANDARD = 35

	// 项目规划需求数 分值
	PROJECT_STORY_NUM_STANDARD = 15
	PROJECTED_STORY_STANDARD = 0.5
	DEVELOPED_STORY_STANDARD = 1
	CLOSED_STORY_STANDARD = 1

	// 预估承诺完成率 分值
	PROJECT_ESTIMATE_STANDARD = 10
)

type (
	PmKpi struct {
		Accounts []string // pm的账号
		Db       *sql.DB  // 数据库连接
	}

	PmKpiGrade struct {
		Account string // 禅道账号

		StartTime string // 开始时间
		EndTime   string // 结束时间

		// 项目软件项目进度达成率
		ProgressAvgDiffDays float64 // 平均项目进度预估天数差值
		ProgressStandard float64 // 项目进度达成基数
		ProgressStandardGrade float64 // 项目进度达成率 实际分数

		// 项目成果完成率
		CompleteRate float64 // 项目成果完成率
		CompleteRateStandard float64 // 项目成果完成率基数
		CompleteRateStandardGrade float64 // 项目成果完成率实际分数

		// 项目成果完成率,完成情况
		ProjectCompleteList []ProjectCompleteInfo

		// 项目规划需求数
		ProjectedStoryNum int // 评审完的需求数
		DevelopedStoryNum int // 开发完的需求数
		ClosedStoryNum int // 关闭的需求数
		StoryNumGrade float64 // 需求数实际分数

		// 预估承诺完成率
		PromiseDiffDays float64 // 预估承诺完成率
		PromiseStandard float64 // 预估承诺完成率基数
		PromiseStandardGrade float64 // 预估承诺完成率实际分数

		TotalGrade float64 // 总分数

	}

	ProjectCompleteInfo struct {
		ProjectName string // 项目名称
		CompleteRate float64 // 项目成果完成率
	}
)
// NewPmKpi creates a new PmKpi object
func NewPmKpi(db *sql.DB, accounts []string) *PmKpi {
	return &PmKpi{
		Accounts: accounts,
		Db:       db,
	}
}


// GetPmKpiGrade gets the PM KPI information
func (l *PmKpi) GetPmKpiGrade() map[string]PmKpiGrade {
	kpiGrades := make(map[string]PmKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = PmKpiGrade{
			Account: account,
		}
	}

	// 项目软件项目进度达成率
	progressResult := dbQuery.QueryProjectProgress(l.Db, l.Accounts)
	for account, result := range progressResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.ProgressAvgDiffDays = result.AvgDiffDays
			tmp.ProgressStandard = result.ProgressStandard
			tmp.ProgressStandardGrade = result.ProgressStandard * PROJECT_PROGRESS_STANDARD
			tmp.TotalGrade += tmp.ProgressStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 项目成果完成率
	completeRateResult := dbQuery.QueryProjectCompleteRate(l.Db, l.Accounts)
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
	completeRateDetailResult := dbQuery.QueryProjectCompleteRateDetail(l.Db, l.Accounts)
	for account, result := range completeRateDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmp.ProjectCompleteList = append(tmp.ProjectCompleteList, ProjectCompleteInfo{
					ProjectName: r.ProjectName,
					CompleteRate: r.CompleteRate,
				})
			}
			kpiGrades[account] = tmp
		}
	}

	// 项目规划需求数
	storyNumResult := dbQuery.QueryProjectStoryNum(l.Db, l.Accounts)
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
	projectPromiseResult := dbQuery.QueryProjectEstimateRate(l.Db, l.Accounts)
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

	return kpiGrades
}
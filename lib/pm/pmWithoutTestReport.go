package pm

import (
	"database/sql"
	dbQuery "kpi-bot/db"
)


const (
	// 项目软件项目进度达成率 分值
	PROJECT_PROGRESS_STANDARD_WITHOUT_TESTREPORT = 40

	// 项目成果完成率 分值
	PROJECT_COMPLETEMENT_STANDARD_WITHOUT_TESTREPORT = 20

	// 项目规划需求数 分值
	PROJECT_STORY_NUM_STANDARD_WITHOUT_TESTREPORT = 10
	PROJECTED_STORY_STANDARD_WITHOUT_TESTREPORT = 0.5
	DEVELOPED_STORY_STANDARD_WITHOUT_TESTREPORT = 1
	CLOSED_STORY_STANDARD_WITHOUT_TESTREPORT = 1

	// 预估承诺完成率 分值
	PROJECT_ESTIMATE_STANDARD_WITHOUT_TESTREPORT = 0

	// 项目预估工时准确率 分值
	TIME_ESTIMATE_STANDARD_WITHOUT_TESTREPORT = 30

	// 系数
	TOP_COEFFICIENT_WITHOUT_TESTREPORT = 1.2
	SECOND_COEFFICIENT_WITHOUT_TESTREPORT = 1.0
	THIRD_COEFFICIENT_WITHOUT_TESTREPORT = 0.7
)

type (
	PmKpiWithoutTestReport struct {
		Accounts []string // pm的账号
		Db       *sql.DB  // 数据库连接
		StartTime string // 开始时间
		EndTime string // 结束时间
	}

	PmKpiGradeWithoutTestReport struct {
		Account string // 禅道账号

		StartTime string // 开始时间
		EndTime   string // 结束时间

		// 项目软件项目进度达成率
		ProgressAvgDiffDays float64 // 平均项目进度预估天数差值
		ProgressStandard float64 // 项目进度达成基数
		ProgressStandardGrade float64 // 项目进度达成率 实际分数

		// 项目软件项目进度达成率 完成情况
		ProjectProgressList []ProjectProgressInfoWithoutTestReport

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

		// 项目预估工时准确率
		TimeEstimateRate float64 // 项目预估工时准确率
		TimeEstimateStandard float64 // 项目预估工时准确率基数
		TimeEstimateGrade float64 // 项目预估工时准确率实际分数

		// 项目预估工时准确率 完成情况
		TimeEstimateList []TimeEstimateInfo

		TotalGrade float64 // 总分数
		TotalGradeStandard float64 // 总分数基数

	}

	ProjectProgressInfoWithoutTestReport struct {
		Account        string // 禅道账号
		ProjectId      int64  // 项目id
		ProjectName    string // 项目名称
		ProjectType    string // 项目类型
		ProjectEnd     string // 项目预估结束时间
		ProjectRealEnd string // 项目实际结束时间
		ProjectDiff    int    // 与预期相差天数
	}

	ProjectCompleteInfoWithoutTestReport struct {
		ProjectName string // 项目名称
		CompleteRate float64 // 项目成果完成率
	}
)
// NewPmKpi creates a new PmKpi object
func NewPmKpiWithoutTestReport(db *sql.DB, accounts []string, startTime, endTime string) *PmKpiWithoutTestReport {
	return &PmKpiWithoutTestReport{
		Accounts: accounts,
		Db:       db,
		StartTime: startTime,
		EndTime: endTime,
	}
}


// GetPmKpiGrade gets the PM KPI information
func (l *PmKpiWithoutTestReport) GetPmKpiGradeWithoutTestReport() map[string]PmKpiGradeWithoutTestReport {
	kpiGrades := make(map[string]PmKpiGradeWithoutTestReport)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = PmKpiGradeWithoutTestReport{
			Account: account,
			StartTime: l.StartTime,
			EndTime: l.EndTime,
		}
	}

	// 项目软件项目进度达成率
	progressResult := dbQuery.QueryProjectProgressWithout(l.Db, l.Accounts, l.StartTime, l.EndTime)
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

	// 项目软件项目进度达成率 完成情况
	progressDetailResult := dbQuery.QueryProjectProgressDetailWithoutTestReport(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range progressDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmp.ProjectProgressList = append(tmp.ProjectProgressList, ProjectProgressInfoWithoutTestReport{
					Account:        r.Account,
					ProjectId:      r.ProjectId,
					ProjectName:    r.ProjectName,
					ProjectType:    r.ProjectType,
					ProjectEnd:     r.ProjectEnd,
					ProjectRealEnd: r.ProjectRealEnd,
					ProjectDiff:    r.ProjectDiff,
				})
			}
			kpiGrades[account] = tmp
		}
	}

	// 项目成果完成率
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
	completeRateDetailResult := dbQuery.QueryProjectCompleteRateDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
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
	timeEstimateDetailResult := dbQuery.QueryProjectTimeEstimateRateDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range timeEstimateDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmp.TimeEstimateList = append(tmp.TimeEstimateList, TimeEstimateInfo{
					Account:       r.Account,
					StoryId:       r.StoryId,
					Title:         r.Title,
					Estimate:      r.Estimate,
					StoryConsumed: r.StoryConsumed,
					EstimateRate:  r.EstimateRate,
				})
			}
			kpiGrades[account] = tmp
		}
	}
	
	for account, kpiGrade := range kpiGrades {
		tmp := kpiGrades[account]
		tmp.TotalGradeStandard = l.GetRdKpiGradeStandardWithoutTestReport(kpiGrade.TotalGrade)
		kpiGrades[account] = tmp
	}


	return kpiGrades
}


// 计算得分系数
func (l *PmKpiWithoutTestReport) GetRdKpiGradeStandardWithoutTestReport(totalGrade float64) float64 {
	if totalGrade >= 100 {
		return TOP_COEFFICIENT
	} else if totalGrade < 100 && totalGrade >= 80 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 80 && totalGrade >= 60 {
		return THIRD_COEFFICIENT
	}
	return 0
}
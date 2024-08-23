package rd

import (
	"database/sql"
	"fmt"
	"kpi-bot/common"
	dbQuery "kpi-bot/db"
)

const (
	// 项目进度延时率 分值
	PROJECT_PROGRESS_STANDARD_WITHOUT_TESTREPORT = 30

	// 需求完成率 分值
	STORY_STANDARD_WITHOUT_TESTREPORT = 45

	// bug遗留率 分值
	BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT = 20
	
	// 工时预估达成比 分值
	TIME_ESTIMATE_STANDARD_WITHOUT_TESTREPORT = 15
)

type (
	RdWithoutTestReportKpi struct {
		Accounts []string // rd的账号
		Db       *sql.DB  // 数据库连接
		StartTime string
		EndTime string
	}

	RdWithoutTestReportKpiGrade struct {
		Account string // 禅道账号

		StartTime string // 开始时间
		EndTime   string // 结束时间

		ProjectTotalSaturdays float64
		ProjectTotalSundays   float64

		RealTotalSaturdays float64
		RealTotalSundays   float64

		// 项目进度延时率
		// AvgDiffExpect            float64               // 平均项目进度预估天数差值
		SumPlanDiffDays          float64               // 平均项目进度预估天数差值
		SumRealDiffDays          float64               // 平均项目进度实际天数差值
		AvgDiffRate              float64               // 平均项目进度延时率
		AvgProgressStandard      float64               // 项目进度延时率基数
		AvgProgressStandardGrade float64               // 项目进度达成率 实际分数
		ProjectProgressList      []ProjectProgressInfo // 项目进度达成率 详情

		// 需求完成率
		TotalStoryScore float64     // 需求总分数
		StoryList       []StoryInfo // 需求详情

		// bug遗留率
		BugCarryOverRate float64        // bug遗留率
		BugCarryStandard float64        // bug遗留率基数
		BugCarryStandardGrade float64   // bug遗留率 实际分数
		BugInfoList      []BugCarryInfoWithoutTestReport // bug遗留详情

		// 工时预估达成比
		TimeEstimateRate     float64 // 工时预估达成比
		TimeEstimateStandard float64 // 工时预估基数
		TimeEstimateStandardGrade float64 // 工时预估实际分数

		TotalGrade float64 // 总分数
		TotalGradeStandard float64 // 总分数基数
	}

	BugCarryInfoWithoutTestReport struct {
		Account       string // 禅道账号
		ProjectName    string // 项目名称
		BugId         int64  // bug id
		BugTitle      string // bug标题
		BugResolution string // bug解决方案
		BugStatus     string // bug状态
	}
)

// NewRdKpi 创建一个研发KPI对象
func NewRdKpiWithoutTestReport(db *sql.DB, accounts []string, startTime, endTime string) *RdWithoutTestReportKpi {
	return &RdWithoutTestReportKpi{
		Accounts: accounts,
		Db:       db,
		StartTime: startTime,
		EndTime: endTime,
	}
}

// GetRdKpiGrade 获取研发KPI信息
func (l *RdWithoutTestReportKpi) GetRdKpiWithoutTestReportGrade() map[string]RdWithoutTestReportKpiGrade {
	kpiGrades := make(map[string]RdWithoutTestReportKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = RdWithoutTestReportKpiGrade{
			Account: account,
			StartTime: l.StartTime,
			EndTime: l.EndTime,
		}
	}

	// 项目进度完成情况
	projectProgressDetailResult := dbQuery.QueryRdProjectProgressDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range projectProgressDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				planSaturdays, planSundays := common.CountWeekends(r.Begin, r.End)
				realSaturdays, realSundays := common.CountWeekends(r.Begin, r.RealEnd)
				tmp.ProjectProgressList = append(tmp.ProjectProgressList, ProjectProgressInfo{
					ProjectId:  r.ProjectId,
					ProjectName: r.ProjectName,
					Begin:      r.Begin,
					End:        r.End,
					RealEnd:    r.RealEnd,
					PlanDiff:   r.PlanDiff - float64(planSaturdays + planSundays),
					RealDiff:   r.RealDiff - float64(realSaturdays + realSundays),
				})
				tmp.ProjectTotalSaturdays += float64(planSaturdays) / 2
				tmp.ProjectTotalSundays += float64(planSundays)

				tmp.RealTotalSaturdays += float64(realSaturdays) / 2
				tmp.RealTotalSundays += float64(realSundays)
			}
			fmt.Printf("account: %v, ProjectTotalSaturdays: %v, ProjectTotalSundays: %v, RealTotalSaturdays: %v, RealTotalSundays: %v\n", account, tmp.ProjectTotalSaturdays, tmp.ProjectTotalSundays, tmp.RealTotalSaturdays, tmp.RealTotalSundays)
			kpiGrades[account] = tmp
		}
	}


	// 项目进度达成率
	projectProgressResult := dbQuery.QueryRdProjectProgress(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range projectProgressResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.SumPlanDiffDays = result.SumPlanDiffDays - tmp.ProjectTotalSaturdays - tmp.ProjectTotalSundays
			tmp.SumRealDiffDays = result.SumRealDiffDays - tmp.RealTotalSaturdays - tmp.RealTotalSundays
			tmp.AvgDiffRate = common.GetProjectProgressExpectRate(tmp.SumPlanDiffDays, tmp.SumRealDiffDays)
			tmp.AvgProgressStandard = GetRdProjectProgressStandard(tmp.AvgDiffRate)
			tmp.AvgProgressStandardGrade = tmp.AvgProgressStandard * PROJECT_PROGRESS_STANDARD
			tmp.TotalGrade += tmp.AvgProgressStandardGrade
			fmt.Printf("account: %v, SumPlanDiffDays: %v, SumRealDiffDays: %v, AvgDiffRate: %v, AvgProgressStandard: %v, AvgProgressStandardGrade: %v, TotalGrade: %v\n", account, tmp.SumPlanDiffDays, tmp.SumRealDiffDays, tmp.AvgDiffRate, tmp.AvgProgressStandard, tmp.AvgProgressStandardGrade, tmp.TotalGrade)
			kpiGrades[account] = tmp
		}
	}

	// 需求达成率
	// storyScoreResult := dbQuery.QueryRdStoryScore(l.Db, l.Accounts, l.StartTime, l.EndTime)
	// for account, result := range storyScoreResult {
	// 	if _, ok := kpiGrades[account]; ok {
	// 		tmp := kpiGrades[account]
	// 		if result.Score > STORY_STANDARD_WITHOUT_TESTREPORT {
	// 			tmp.TotalStoryScore = STORY_STANDARD_WITHOUT_TESTREPORT
	// 		} else {
	// 			tmp.TotalStoryScore = result.Score
	// 		}
	// 		tmp.TotalGrade += tmp.TotalStoryScore
	// 		kpiGrades[account] = tmp
	// 	}
	// }

	// 需求完成情况
	storyDetailResult := dbQuery.QueryRdStoryDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range storyDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			totalScore := float64(0)
			for _, r := range result {
				score := r.Estimate / STORY_BASE_TIME * STORY_BASE_SCORE
				totalScore += score
				tmp.StoryList = append(tmp.StoryList, StoryInfo{
					Id: r.StoryId,
					Account: r.Account,
					Title: r.Title,
					Estimate: r.Estimate,
					Score: score,
				})
			}
			if totalScore > STORY_STANDARD_WITHOUT_TESTREPORT {
				tmp.TotalStoryScore = STORY_STANDARD_WITHOUT_TESTREPORT
			} else {
				tmp.TotalStoryScore = totalScore
			}
			tmp.TotalGrade += tmp.TotalStoryScore
			kpiGrades[account] = tmp
		}
	}

	// 项目版本bug遗留率 无测试报告
	bugCarryOverResult := dbQuery.QueryRdBugCarryOverWithoutTestReport(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range bugCarryOverResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.BugCarryOverRate = result.BugCarryOverRate
			tmp.BugCarryStandard = result.BugCarryOverRateStandard
			tmp.BugCarryStandardGrade = result.BugCarryOverRateStandard * BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT
			tmp.TotalGrade += tmp.BugCarryStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 項目版本bug遗留實際情況 无测试报告
	bugCarryDetailResult := dbQuery.QueryRdBugCarryOverDetailWithoutTestReport(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range bugCarryDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmp.BugInfoList = append(tmp.BugInfoList, BugCarryInfoWithoutTestReport{
					Account: r.Account,
					ProjectName: r.ProjectName,
					BugId: r.BugId,
					BugTitle: r.BugTitle,
					BugResolution: r.BugResolution,
					BugStatus: r.BugStatus,
				})
			}
			kpiGrades[account] = tmp
		}
	}



	// 工时预估达成比
	timeEstimateRateResult := dbQuery.QueryRdTimeEstimateRate(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range timeEstimateRateResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.TimeEstimateRate = result.TimeEstimateRate
			tmp.TimeEstimateStandard = result.TimeEstimateRateStandard
			tmp.TimeEstimateStandardGrade = result.TimeEstimateRateStandard * TIME_ESTIMATE_STANDARD_WITHOUT_TESTREPORT
			tmp.TotalGrade += tmp.TimeEstimateStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 结算系数
	for account := range kpiGrades {
		tmp := kpiGrades[account]
		// if len(tmp.BugInfoList) == 0 {
		// 	tmp.TotalGrade -= tmp.BugCarryStandardGrade
		// 	tmp.TotalGrade += BUG_CARRY_OVER_STANDARD
		// }
		tmp.TotalGradeStandard = l.GetRdKpiGradeStandard(tmp.TotalGrade)
		kpiGrades[account] = tmp
	}

	return kpiGrades
}

// 计算得分系数
func (l *RdWithoutTestReportKpi) GetRdKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 100 {
		return TOP_COEFFICIENT
	} else if totalGrade < 100 && totalGrade >= 80 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 80 && totalGrade >= 60 {
		return THIRD_COEFFICIENT
	}
	return 0
}
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

	STORY_BASE_SCORE_WITHOUT_TESTREPORT = 0.021 // 分值

	// 需求完成率 分值
	STORY_STANDARD_WITHOUT_TESTREPORT = 30

	// bug遗留率 分值
	BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT = 25

	// 工时预估达成比 分值
	TIME_ESTIMATE_STANDARD_WITHOUT_TESTREPORT = 15
)

type (
	RdWithoutTestReportKpi struct {
		Accounts  []string // rd的账号
		Db        *sql.DB  // 数据库连接
		StartTime string
		EndTime   string
		Pms       []string // 关联的项目经理
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
		BugCarryOverRate      float64                         // bug遗留率
		BugCarryStandard      float64                         // bug遗留率基数
		BugCarryStandardGrade float64                         // bug遗留率 实际分数
		BugInfoList           []BugCarryInfoWithoutTestReport // bug遗留详情

		// 工时预估达成比
		TimeEstimateRate          float64 // 工时预估达成比
		TimeEstimateStandard      float64 // 工时预估基数
		TimeEstimateStandardGrade float64 // 工时预估实际分数

		TotalGrade         float64 // 总分数
		TotalGradeStandard float64 // 总分数基数
	}

	BugCarryInfoWithoutTestReport struct {
		Account       string // 禅道账号
		ProjectName   string // 项目名称
		BugId         int64  // bug id
		BugTitle      string // bug标题
		BugResolution string // bug解决方案
		BugStatus     string // bug状态
	}
)

// NewRdKpi 创建一个研发KPI对象
func NewRdKpiWithoutTestReport(db *sql.DB, accounts, rdpms []string, startTime, endTime string) *RdWithoutTestReportKpi {
	return &RdWithoutTestReportKpi{
		Accounts:  accounts,
		Db:        db,
		StartTime: startTime,
		EndTime:   endTime,
		Pms:       rdpms,
	}
}

// GetRdKpiGrade 获取研发KPI信息
func (l *RdWithoutTestReportKpi) GetRdKpiWithoutTestReportGrade() map[string]RdWithoutTestReportKpiGrade {
	kpiGrades := make(map[string]RdWithoutTestReportKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = RdWithoutTestReportKpiGrade{
			Account:               account,
			StartTime:             l.StartTime,
			EndTime:               l.EndTime,
			BugCarryStandardGrade: BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT,
			TotalGrade: BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT, // 起始總分先加bug分 25分
		}
	}

	// 项目进度完成情况
	fmt.Print("项目进度完成情况\n")
	projectProgressDetailResult := dbQuery.QueryRdProjectProgressDetail(l.Db, l.Accounts, l.Pms, l.StartTime, l.EndTime)
	for account, result := range projectProgressDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				planSaturdays, planSundays := common.CountWeekends(r.Begin, r.End)
				realSaturdays, realSundays := common.CountWeekends(r.Begin, r.RealEnd)
				tmp.ProjectProgressList = append(tmp.ProjectProgressList, ProjectProgressInfo{
					ProjectId:   r.ProjectId,
					ProjectName: r.ProjectName,
					Begin:       r.Begin,
					End:         r.End,
					RealEnd:     r.RealEnd,
					PlanDiff:    r.PlanDiff - float64(planSaturdays+planSundays),
					RealDiff:    r.RealDiff - float64(realSaturdays+realSundays),
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
			fmt.Printf("account: %v, ProjectTotalSaturdays: %v, ProjectTotalSundays: %v, RealTotalSaturdays: %v, RealTotalSundays: %v\n", account, tmp.ProjectTotalSaturdays, tmp.ProjectTotalSundays, tmp.RealTotalSaturdays, tmp.RealTotalSundays)
			kpiGrades[account] = tmp
		}
	}

	// 项目进度达成率
	fmt.Print("项目进度达成率\n")
	projectProgressResult := dbQuery.QueryRdProjectProgress(l.Db, l.Accounts, l.Pms, l.StartTime, l.EndTime)
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
	fmt.Print("需求完成情况\n")
	storyDetailResult := dbQuery.QueryRdStoryDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range storyDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			totalScore := float64(0)
			for _, r := range result {
				score := r.Estimate / STORY_BASE_TIME * STORY_BASE_SCORE_WITHOUT_TESTREPORT
				totalScore += score
				tmp.StoryList = append(tmp.StoryList, StoryInfo{
					Id:       r.StoryId,
					Account:  r.Account,
					Title:    r.Title,
					Estimate: r.Estimate,
					Score:    score,
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
	// fmt.Print("项目版本bug遗留率 无测试报告\n")
	// bugCarryOverResult := dbQuery.QueryRdBugCarryOverWithoutTestReport(l.Db, l.Accounts, l.StartTime, l.EndTime)
	// for account, result := range bugCarryOverResult {
	// 	if _, ok := kpiGrades[account]; ok {
	// 		tmp := kpiGrades[account]
	// 		tmp.BugCarryOverRate = result.BugCarryOverRate
	// 		tmp.BugCarryStandard = result.BugCarryOverRateStandard
	// 		fmt.Printf("account: %v, BugCarryOverRate: %v, BugCarryStandard: %v\n", account, tmp.BugCarryOverRate, tmp.BugCarryStandard)
	// 		tmp.BugCarryStandardGrade = result.BugCarryOverRateStandard * BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT
	// 		tmp.TotalGrade += tmp.BugCarryStandardGrade
	// 		kpiGrades[account] = tmp
	// 	}
	// }

	// 項目版本bug遗留實際情況 无测试报告
	fmt.Print("項目版本bug遗留实际情况 无测试报告\n")
	bugCarryDetailResult := dbQuery.QueryRdBugCarryOverDetailWithoutTestReport(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range bugCarryDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.BugCarryOverRate = 0
			activeBugCount := float64(0)
			fixBugCount := float64(0)
			for _, r := range result {
				projectName := ""
				if r.ProjectName == nil {
					projectName = ""
				} else {
					projectName = *r.ProjectName
				}

				if r.BugStatus == "active" {
					activeBugCount++
				} else {
					fixBugCount++
				}
				tmp.BugInfoList = append(tmp.BugInfoList, BugCarryInfoWithoutTestReport{
					Account:       r.Account,
					ProjectName:   projectName,
					BugId:         r.BugId,
					BugTitle:      r.BugTitle,
					BugResolution: r.BugResolution,
					BugStatus:     r.BugStatus,
				})
			}
			// 如果bug总数大于0 算出bug遗留率
			if activeBugCount+fixBugCount > 0 {
				tmp.BugCarryOverRate = activeBugCount / (activeBugCount + fixBugCount)
			}
			tmp.BugCarryStandard = common.GetBugStandard(tmp.BugCarryOverRate)
			tmp.BugCarryStandardGrade = tmp.BugCarryStandard * BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT
			tmp.TotalGrade += tmp.BugCarryStandardGrade - BUG_CARRY_OVER_STANDARD_WITHOUT_TESTREPORT // 有bug的記得先扣掉起始加的25分
			kpiGrades[account] = tmp
		}
	}

	// 工时预估达成比
	fmt.Print("工时预估达成比\n")
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
		tmp.TotalGradeStandard = l.GetKpiGradeStandard(tmp.TotalGrade)
		kpiGrades[account] = tmp
	}

	return kpiGrades
}

// 计算得分系数
func (l *RdWithoutTestReportKpi) GetKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 90 {
		return TOP_COEFFICIENT
	} else if totalGrade < 90 && totalGrade >= 70 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 70 && totalGrade >= 60 {
		return THIRD_COEFFICIENT
	}
	return 0
}

package rd

import (
	"database/sql"
	"kpi-bot/common"
	dbQuery "kpi-bot/db"
)

const (
	// 项目进度延时率 分值
	PROJECT_PROGRESS_STANDARD = 30

	// 需求完成率 分值
	STORY_STANDARD = 30

	// 需求基础 分值
	STORY_BASE_TIME  = 0.1   // 小时
	STORY_BASE_SCORE = 0.03 // 分值

	// bug遗留率 分值
	BUG_CARRY_OVER_STANDARD = 20

	// 工时预估达成比 分值
	TIME_ESTIMATE_STANDARD = 10

	// 版本提测次数 分值
	PUB_TIMES_STANDARD = 10

	// 系数
	TOP_COEFFICIENT    = 1.2
	SECOND_COEFFICIENT = 1.0
	THIRD_COEFFICIENT  = 0.8

	// 项目进度延时率基数
	PROJECT_PROGRESS_Level1 = 1.5
	PROJECT_PROGRESS_Level2 = 1.2
	PROJECT_PROGRESS_Level3 = 1
	PROJECT_PROGRESS_Level4 = 0.8
	PROJECT_PROGRESS_Level5 = 0.5
)

type (
	RdKpi struct {
		Accounts  []string // rd的账号
		Db        *sql.DB  // 数据库连接
		StartTime string   // 开始时间
		EndTime   string   // 结束时间
		Pms       []string // 关联的项目经理
	}

	RdKpiGrade struct {
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
		BugCarryOverRate      float64        // bug遗留率
		BugCarryStandard      float64        // bug遗留率基数
		BugCarryStandardGrade float64        // bug遗留率 实际分数
		BugInfoList           []BugCarryInfo // bug遗留详情

		// 測試單數目
		TestTaskCount int

		// 工时预估达成比
		TimeEstimateRate          float64 // 工时预估达成比
		TimeEstimateStandard      float64 // 工时预估基数
		TimeEstimateStandardGrade float64 // 工时预估实际分数

		// 版本提测次数
		AvgPubTimes              float64       // 平均提测次数
		AvgPubTimesStandard      float64       // 平均提测次数基数
		AvgPubTimesStandardGrade float64       // 平均提测次数实际分数
		PubTimeList              []PubTimeInfo // 提测详情

		TotalGrade         float64 // 总分数
		TotalGradeStandard float64 // 总分数基数
	}

	ProjectProgressInfo struct {
		ProjectId   int     // 项目id
		ProjectName string  // 项目名称
		Begin       string  // 计划开始时间
		End         string  // 计划结束时间
		RealEnd     string  // 实际结束时间
		PlanDiff    float64 // 计划天数差值
		RealDiff    float64 // 实际天数差值
	}

	StoryInfo struct {
		Id       int64   // 需求id
		Account  string  // 禅道账号
		Title    string  // 需求标题
		Estimate float64 // 预估工时
		Score    float64 // 需求分数
	}

	BugCarryInfo struct {
		Account       string // 禅道账号
		ProjectName   string // 项目名称
		BugId         int64  // bug id
		BugTitle      string // bug标题
		BugResolution string // bug解决方案
		BugStatus     string // bug状态
	}

	PubTimeInfo struct {
		Account     string // 禅道账号
		ProjectType string // 项目类型
		ProjectName string // 项目名称
		PubTimes    int    // 发版次数
		LastPubTime string // 最后一次提测时间
	}
)

// NewRdKpi 创建一个研发KPI对象
func NewRdKpi(db *sql.DB, accounts, pms []string, startTime, endTime string) *RdKpi {
	return &RdKpi{
		Accounts:  accounts,
		Db:        db,
		StartTime: startTime,
		EndTime:   endTime,
		Pms:       pms,
	}
}

// GetRdKpiGrade 获取研发KPI信息
func (l *RdKpi) GetRdKpiGrade() map[string]RdKpiGrade {
	kpiGrades := make(map[string]RdKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = RdKpiGrade{
			Account:               account,
			StartTime:             l.StartTime,
			EndTime:               l.EndTime,
			BugCarryStandardGrade: BUG_CARRY_OVER_STANDARD,
			TotalGrade: BUG_CARRY_OVER_STANDARD, // 起始總分先加bug分 25分
		}
	}

	// 项目进度完成情况
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
			kpiGrades[account] = tmp
		}
	}

	// 项目进度达成率
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
			kpiGrades[account] = tmp
		}
	}

	// 需求达成率
	// storyScoreResult := dbQuery.QueryRdStoryScore(l.Db, l.Accounts, l.StartTime, l.EndTime)
	// for account, result := range storyScoreResult {
	// 	if _, ok := kpiGrades[account]; ok {
	// 		tmp := kpiGrades[account]
	// 		if result.Score > STORY_STANDARD {
	// 			tmp.TotalStoryScore = STORY_STANDARD
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
					Id:       r.StoryId,
					Account:  r.Account,
					Title:    r.Title,
					Estimate: r.Estimate,
					Score:    score,
				})
			}

			if totalScore > STORY_STANDARD {
				tmp.TotalStoryScore = STORY_STANDARD
			} else {
				tmp.TotalStoryScore = totalScore
			}
			tmp.TotalGrade += tmp.TotalStoryScore
			kpiGrades[account] = tmp
		}
	}

	// 測試單數量
	taskCountResult := dbQuery.QueryCountTestTask(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range taskCountResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.TestTaskCount = result
			kpiGrades[account] = tmp
		}
	}

	// 项目版本bug遗留率 無測試報告
	// bugCarryOverResult := dbQuery.QueryRdBugCarryOverWithoutTestReport(l.Db, l.Accounts, l.StartTime, l.EndTime)
	// for account, result := range bugCarryOverResult {
	// 	if _, ok := kpiGrades[account]; ok {
	// 		tmp := kpiGrades[account]
	// 		if tmp.TestTaskCount > 0 {
	// 			tmp.BugCarryOverRate = result.BugCarryOverRate
	// 			tmp.BugCarryStandard = result.BugCarryOverRateStandard
	// 			tmp.BugCarryStandardGrade = result.BugCarryOverRateStandard * BUG_CARRY_OVER_STANDARD
	// 			tmp.TotalGrade += tmp.BugCarryStandardGrade
	// 			kpiGrades[account] = tmp
	// 		}

	// 	}
	// }

	// 項目版本bug遗留實際情況 無測試報告
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
				tmp.BugInfoList = append(tmp.BugInfoList, BugCarryInfo{
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
			tmp.BugCarryStandardGrade = tmp.BugCarryStandard * BUG_CARRY_OVER_STANDARD
			tmp.TotalGrade += tmp.BugCarryStandardGrade - BUG_CARRY_OVER_STANDARD // 有bug的記得先扣掉起始加的25分
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
			tmp.TimeEstimateStandardGrade = result.TimeEstimateRateStandard * TIME_ESTIMATE_STANDARD
			tmp.TotalGrade += tmp.TimeEstimateStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 版本发版次数平均发版次数
	pubTimesResult := dbQuery.QueryRdPubTimes(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range pubTimesResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			tmp.AvgPubTimes = result.AvgPubTimes
			tmp.AvgPubTimesStandard = result.AvgPubTimesStandard
			tmp.AvgPubTimesStandardGrade = result.AvgPubTimesStandard * PUB_TIMES_STANDARD
			tmp.TotalGrade += tmp.AvgPubTimesStandardGrade
			kpiGrades[account] = tmp
		}
	}

	// 版本发版次数详情
	pubTimesDetailResult := dbQuery.QueryRdPubTimesDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range pubTimesDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmp.PubTimeList = append(tmp.PubTimeList, PubTimeInfo{
					Account:     r.Account,
					ProjectType: r.ProjectType,
					ProjectName: r.ProjectName,
					PubTimes:    r.PubTimes,
					LastPubTime: r.LastPubTime,
				})
			}
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
func (l *RdKpi) GetKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 90 {
		return TOP_COEFFICIENT
	} else if totalGrade < 90 && totalGrade >= 70 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 70 && totalGrade >= 60 {
		return THIRD_COEFFICIENT
	}
	return 0
}

func GetRdProjectProgressStandard(avgDiffRate float64) float64 {
	if avgDiffRate <= -0.5 {
		return PROJECT_PROGRESS_Level1
	} else if avgDiffRate > -0.5 && avgDiffRate <= -0.2 {
		return PROJECT_PROGRESS_Level2
	} else if avgDiffRate > -0.2 && avgDiffRate <= 0 {
		return PROJECT_PROGRESS_Level3
	} else if avgDiffRate > 0 && avgDiffRate <= 0.2 {
		return PROJECT_PROGRESS_Level4
	} else if avgDiffRate > 0.2 && avgDiffRate <= 0.5 {
		return PROJECT_PROGRESS_Level5
	} else {
		return 0
	}
}

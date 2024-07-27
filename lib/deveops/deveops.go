package deveops

import (
	"database/sql"
	"fmt"
	"kpi-bot/common"
	dbQuery "kpi-bot/db"
	"kpi-bot/lib/rd"
)



const (
	// 项目进度延时率 分值
	PROJECT_PROGRESS_STANDARD = 30

	// 需求完成率 分值
	STORY_STANDARD = 30

	// 线上故障率 分值
	SERVICE_FAILURE_RATE = 30

	// 工时预估达成比 分值
	TIME_ESTIMATE_STANDARD = 10

	// 发布故障率 分值
	RELEASE_FAILURE_RATE = 10

	// 系数
	TOP_COEFFICIENT    = 1.2
	SECOND_COEFFICIENT = 1.0
	THIRD_COEFFICIENT  = 0.7
)


type (
	DeveopsKpi struct {
		Accounts []string // rd的账号
		Db       *sql.DB  // 数据库连接
		StartTime string
		EndTime string
	}

	DeveopsKpiGrade struct {
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

	ProjectProgressInfo struct {
		ProjectId   int  // 项目id
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
)


// NewDeveopsKpi 创建一个运维KPI对象
func NewDeveopsKpi(db *sql.DB, accounts []string, startTime, endTime string) *DeveopsKpi {
	return &DeveopsKpi{
		Accounts: accounts,
		Db:       db,
		StartTime: startTime,
		EndTime: endTime,
	}
}


// GetDeveopsKpi 获取运维KPI
func (l *DeveopsKpi) GetDeveopsKpiGrade() map[string]DeveopsKpiGrade {
	kpiGrades := make(map[string]DeveopsKpiGrade)

	// 建立所有账户啊kpi信息
	for _, account := range l.Accounts {
		kpiGrades[account] = DeveopsKpiGrade{
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
			tmp.AvgProgressStandard = rd.GetRdProjectProgressStandard(tmp.AvgDiffRate)
			tmp.AvgProgressStandardGrade = tmp.AvgProgressStandard * PROJECT_PROGRESS_STANDARD
			tmp.TotalGrade += tmp.AvgProgressStandardGrade
			fmt.Printf("account: %v, SumPlanDiffDays: %v, SumRealDiffDays: %v, AvgDiffRate: %v, AvgProgressStandard: %v, AvgProgressStandardGrade: %v, TotalGrade: %v\n", account, tmp.SumPlanDiffDays, tmp.SumRealDiffDays, tmp.AvgDiffRate, tmp.AvgProgressStandard, tmp.AvgProgressStandardGrade, tmp.TotalGrade)
			kpiGrades[account] = tmp
		}
	}

	// 需求达成率
	storyScoreResult := dbQuery.QueryRdStoryScore(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range storyScoreResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			if result.Score > STORY_STANDARD {
				tmp.TotalStoryScore = STORY_STANDARD
			} else {
				tmp.TotalStoryScore = result.Score
			}
			tmp.TotalGrade += tmp.TotalStoryScore
			kpiGrades[account] = tmp
		}
	}

	// 需求完成情况
	storyDetailResult := dbQuery.QueryRdStoryDetail(l.Db, l.Accounts, l.StartTime, l.EndTime)
	for account, result := range storyDetailResult {
		if _, ok := kpiGrades[account]; ok {
			tmp := kpiGrades[account]
			for _, r := range result {
				tmp.StoryList = append(tmp.StoryList, StoryInfo{
					Id: r.StoryId,
					Account: r.Account,
					Title: r.Title,
					Estimate: r.Estimate,
					Score: r.Score,
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
			tmp.TimeEstimateStandardGrade = result.TimeEstimateRateStandard * TIME_ESTIMATE_STANDARD
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
func (l *DeveopsKpi) GetRdKpiGradeStandard(totalGrade float64) float64 {
	if totalGrade >= 100 {
		return TOP_COEFFICIENT
	} else if totalGrade < 100 && totalGrade >= 80 {
		return SECOND_COEFFICIENT
	} else if totalGrade < 80 && totalGrade >= 60 {
		return THIRD_COEFFICIENT
	}
	return 0
}

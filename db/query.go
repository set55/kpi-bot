package dbQuery

import (
	"database/sql"
	"fmt"
	"kpi-bot/common"
	"log"
)

type (
	// 软件研发项目进度延时率
	QueryRdProjectProgressResult struct {
		Account string // 禅道账号
		// AvgDiffExpect   float64 // 平均项目进度预估天数差值
		SumPlanDiffDays float64 // 總项目进度预估天数差值
		SumRealDiffDays float64 // 總项目进度实际天数差值
	}

	// 软件研发项目进度延时率-完成情况
	QueryRdProjectProgressDetailResult struct {
		Account     string  // 禅道账号
		ProjectId   int     // 项目id
		ProjectName string  // 项目名称
		Begin       string  // 计划开始时间
		End         string  // 计划结束时间
		RealEnd     string  // 实际结束时间
		PlanDiff    float64 // 计划天数差值
		RealDiff    float64 // 实际天数差值
	}

	// 需求达成率
	QueryRdStoryResult struct {
		Account string  // 禅道账号
		Score   float64 // 需求总分数
	}

	// 需求达成率-完成情况
	QueryRdStoryDetailResult struct {
		StoryId  int64   // 需求id
		Account  string  // 禅道账号
		Title    string  // 需求标题
		Estimate float64 // 预估工时
		Stage    string  // 需求状态
		Score    float64 // 需求分数
	}

	// 项目版本bug遗留率情况 无/有测试报告
	QueryRdBugCarryOverResult struct {
		Account                  string  // 禅道账号
		BugCarryOverRate         float64 // bug遗留率
		BugCarryOverRateStandard float64 // bug遗留率基数
	}

	// 項目版本bug遗留實際情況 有测试报告
	QueryRdBugCarryOverDetailResult struct {
		Account       string // 禅道账号
		TestReport    string // 测试报告
		BugId         int64  // bug id
		BugTitle      string // bug 标题
		BugResolution string // bug 解决方案
		BugStatus     string // bug 状态
	}

	// 項目版本bug遗留實際情況 有测试报告
	QueryRdBugCarryOverDetailResultWithoutTestReport struct {
		Account       string // 禅道账号
		ProjectName   *string // 项目名称
		BugId         int64  // bug id
		BugTitle      string // bug 标题
		BugResolution string // bug 解决方案
		BugStatus     string // bug 状态
	}

	// 工时预估达成比
	QueryRdTimeEstimateRateResult struct {
		Account                  string  // 禅道账号
		TimeEstimateRate         float64 // 工时预估达成比
		TimeEstimateRateStandard float64 // 工时预估基数
	}

	// 版本发版次数平均发版次数
	QueryRdPubTimesResult struct {
		Account             string  // 禅道账号
		AvgPubTimes         float64 // 平均提测次数
		AvgPubTimesStandard float64 // 平均提测次数基数
	}

	// 版本发版次数详情
	QueryRdPubTimesDetailResult struct {
		Account     string // 禅道账号
		ProjectType string // 项目类型
		ProjectName string // 项目名称
		PubTimes    int    // 提测次数
		LastPubTime string // 最后提测时间
	}

	// 测试软件项目进度 完成情况
	QueryTestProjectProgressDetailResult struct {
		Account         string // 禅道账号
		TestTaskName    string // 测试任务名称
		TestReportTitle string // 测试报告标题
		TestTaskBegin   string // 测试任务开始时间
		TestTaskEnd     string // 测试任务预估结束时间
		TestReportEnd   string // 测试报告实际结束时间
	}

	// 测试软件项目进度达成率
	QueryTestProjectProgressResult struct {
		Account             string  // 禅道账号
		SumRealTestDiffDays float64 // 總測試进度实际天数差值
		SumTestDiffDays     float64 // 總测试进度預估天数差值
	}

	//测试软件项目有效bug率
	//1、测试报告结束时间是当月的
	//2、bug未被删除，bug关联项目属于测试报告关联项目，bug关联版本是测试报告所属版本，bug是焰海打开的，bug解决状态是转需求，延期处理和已解决的，不予解决，这些叫有效bug。
	//3、版本内所有bug，为项目与测试报告相等，并且不是指派给黄卫旗
	QueryTestValidBugRateResult struct {
		Account              string  // 禅道账号
		ValidBugRate         float64 // 有效bug率
		ValidBugRateStandard float64 // 有效bug率基数
	}

	// bug转需求数
	QueryTestBugToStoryResult struct {
		Account string // 禅道账号
		ToStory int    // bug转需求数
	}

	// 用例发现bug率--同build(版本)下，有多少是关联case的
	// "duplicate"重复bug,"bydesign"设计如此,"notrepro"未复现
	QueryTestBugCaseRateResult struct {
		Account         string  // 禅道账号
		CaseBugRate     float64 // 用例发现bug率
		CaseBugStandard float64 // 用例发现bug率基数
	}

	// 项目软件项目进度达成率
	QueryProjectProgressResult struct {
		Account string // 禅道账号
		// AvgDiffDays      float64 // 平均项目进度预估天数差值
		SumRealProjectDiff float64 // 總项目进度实际天数差值
		SumProjectDiff     float64 // 總项目进度預估天数差值
		SumRealTestDiff    float64 // 總测试进度实际天数差值
		SumTestDiff        float64 // 總测试进度預估天数差值
	}

	// 项目软件进度完成情况
	QueryProjectProgressDetailResult struct {
		Account        string // 禅道账号
		ProjectId      int64  // 项目id
		ProjectName    string // 项目名称
		ProjectType    string // 项目类型
		ProjectBegin   string // 项目开始时间
		ProjectEnd     string // 项目预估结束时间
		ProjectRealEnd string // 项目实际结束时间
		ProjectDiff    int    // 与预期相差天数
		TestStart      string // 测试开始时间
		TestEnd        string // 测试结束时间
		TestRealEnd    string // 测试实际结束时间
	}

	// 项目软件进度完成情况
	QueryProjectProgressDetailResultWithoutTestReport struct {
		Account        string // 禅道账号
		ProjectId      int64  // 项目id
		ProjectName    string // 项目名称
		ProjectType    string // 项目类型
		ProjectBegin   string // 项目开始时间
		ProjectEnd     string // 项目预估结束时间
		ProjectRealEnd string // 项目实际结束时间
	}

	// 项目成果完成率,不需要关注执行，只需要看项目需求完成度，因为有执行一定有项目
	QueryProjectCompleteRateResult struct {
		Account              string  // 禅道账号
		CompleteRate         float64 // 项目成果完成率
		CompleteRateStandard float64 // 项目成果完成率基数
	}

	// 项目成果完成率,完成情况
	QueryProjectCompleteRateDetailResult struct {
		Account      string  // 禅道账号
		ProjectName  string  // 项目名称
		CompleteRate float64 // 项目成果完成率
	}

	// 项目规划需求数
	QueryProjectStoryNumResult struct {
		Account  string // 禅道账号
		Stage    string // 阶段
		StoryNum int    // 需求数
	}

	// 预估承诺完成率，只看项目，和最后一个测试报告
	QueryProjectEstimateRateResult struct {
		Account          string  // 禅道账号
		DiffDays         float64 // 平均项目进度预估天数差值
		ProgressStandard float64 // 项目进度达成基数
	}

	// 项目预估工时准确率
	QueryProjectTimeEstimateRateResult struct {
		Account              string  // 禅道账号
		TimeEstimateRate     float64 // 工时预估达成比
		TimeEstimateStandard float64 // 工时预估基数
	}

	// 项目预估工时准确率 完成情况
	QueryProjectTimeEstimateRateResultDetail struct {
		Account       string  // 禅道账号
		StoryId       int64   // 需求id
		Title         string  // 需求标题
		Estimate      float64 // 预估工时
		StoryConsumed float64 // 需求消耗工时
		EstimateRate  *float64 // 预估工时准确率
	}
)

// 软件研发项目进度延时率
func QueryRdProjectProgress(db *sql.DB, accounts, rdpms []string, startTime, endTime string) map[string]QueryRdProjectProgressResult {
	results := map[string]QueryRdProjectProgressResult{}
	sqlCmd := fmt.Sprintf(`
		select tmp.account, SUM(tmp.plan_diff) as sum_plan_diff, SUM(tmp.real_diff) as sum_real_diff
		from (
		select a.account,c.name, c.begin ,c.end,c.realEnd,TIMESTAMPDIFF(DAY,c.begin,c.end) as plan_diff, TIMESTAMPDIFF(DAY,c.begin,c.realEnd) as real_diff
		from zt_user a 
		inner join zt_team b on b.account = a.account 
		inner join zt_project c on c.type in("sprint") and c.id = b.root and c.status = "closed" and c.acl in ("open", "private") and openedBy in (%s)
		where a.account in (%s) and c.realEnd between "%s" and "%s"
		order by a.account,c.realEnd desc
		) tmp
		group by account
	`, common.AccountArrayToString(rdpms), common.AccountArrayToString(accounts), startTime, endTime)
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var result QueryRdProjectProgressResult
		err = rows.Scan(
			&result.Account,
			&result.SumPlanDiffDays,
			&result.SumRealDiffDays,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, SumPlanDiffDays: %f, SumRealDiffDays: %f\n", result.Account, result.SumPlanDiffDays, result.SumRealDiffDays)

		results[result.Account] = result
	}
	return results
}

// 软件研发项目进度延时率-完成情况
func QueryRdProjectProgressDetail(db *sql.DB, accounts, rdpms []string, startTime, endTime string) map[string][]QueryRdProjectProgressDetailResult {
	results := map[string][]QueryRdProjectProgressDetailResult{}
	sqlCmd := fmt.Sprintf(`
		select a.account, c.id as project_id ,c.name as project_name, c.begin , c.end, c.realEnd, TIMESTAMPDIFF(DAY,c.begin,c.end) as plan_diff, TIMESTAMPDIFF(DAY,c.begin,c.realEnd) as real_diff
		from zt_user a 
		inner join zt_team b on b.account = a.account 
		inner join zt_project c on c.type in("sprint") and c.id = b.root and c.status = "closed" and c.deleted="0" and c.acl in ("open", "private") and openedBy in (%s)
		where a.account in (%s) and c.realEnd between "%s" and "%s"
		order by a.account,c.realEnd desc
	`, common.AccountArrayToString(rdpms), common.AccountArrayToString(accounts), startTime, endTime)
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdProjectProgressDetailResult
		err = rows.Scan(
			&result.Account,
			&result.ProjectId,
			&result.ProjectName,
			&result.Begin,
			&result.End,
			&result.RealEnd,
			&result.PlanDiff,
			&result.RealDiff,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, ProjectId: %v, ProjectName: %s, Begin: %v, End: %v, RealEnd: %v, PlanDiff: %v, RealDiff: %v\n", result.Account, result.ProjectId, result.ProjectName, result.Begin, result.End, result.RealEnd, result.PlanDiff, result.RealDiff)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 需求达成率
func QueryRdStoryScore(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryRdStoryResult {
	results := map[string]QueryRdStoryResult{}
	sqlCmd := fmt.Sprintf(`
		select tmp.account,sum(tmp.get_score) 
		from ( 
			select DISTINCT c.id,a.account,c.title,c.estimate, 
				case when c.estimate < 4 then 1
					when c.estimate < 8 and c.estimate >= 4 then 1.5
					when c.estimate < 16 and c.estimate >= 8 then 2
					when c.estimate >= 16 then 2.5
					else 0 end as get_score 
			from zt_user a 
			inner join zt_task b on b.finishedDate between "%s" and "%s" and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
			inner join zt_story c on c.id = b.story and c.stage not in ("waiting","planned","projected","developing") 
			where a.account in (%s) order by a.account desc 
		) tmp 
		group by account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdStoryResult
		err = rows.Scan(
			&result.Account,
			&result.Score,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, Score: %f\n", result.Account, result.Score)

		results[result.Account] = result
	}
	return results
}

// 需求达成率-完成情况
func QueryRdStoryDetail(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryRdStoryDetailResult {
	results := map[string][]QueryRdStoryDetailResult{}
	sqlCmd := fmt.Sprintf(`
		select DISTINCT c.id,a.account,c.title,c.estimate, c.stage,
			case when c.estimate < 4 then 1
					when c.estimate < 8 and c.estimate >= 4 then 1.5
					when c.estimate < 16 and c.estimate >= 8 then 2
					when c.estimate >= 16 then 2.5
					else 0 end as get_score  
		from zt_user a 
		inner join zt_task b on b.finishedDate between "%s" and "%s" and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
		inner join zt_story c on c.id = b.story and c.stage not in ("waiting","planned","projected","developing") 
		where a.account in (%s) order by a.account desc  
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdStoryDetailResult
		err = rows.Scan(
			&result.StoryId,
			&result.Account,
			&result.Title,
			&result.Estimate,
			&result.Stage,
			&result.Score,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		fmt.Printf("StoryId: %d, Account: %s, Title: %s, Estimate: %f, Stage: %s, Score: %f\n", result.StoryId, result.Account, result.Title, result.Estimate, result.Stage, result.Score)
		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 軟件项目版本bug遗留率情况 有测试报告
func QueryRdBugCarryOver(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryRdBugCarryOverResult {
	results := map[string]QueryRdBugCarryOverResult{}
	sqlCmd := fmt.Sprintf(`
		select account ,AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) as bug_carry_over_rate,
		case when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) = 0 then "1.0"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.1 then "0.9"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.2 then "0.8"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.3 then "0.7"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.4 then "0.6"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.5 then "0.5"
			else "0" end as bug_carry_over_rate_standard
		from ( 
			select a.account,c.title,c.project,c.execution,c.builds,count(1) as build_fix_bug,
			(select count(1) from zt_bug where openedBuild = c.builds and deleted = "0" and status = "active" and assignedTo = a.account) as build_postponed_bugs 
			from zt_user a 
			inner join zt_testreport c on c.end between "%s" and "%s" and c.deleted ="0" 
			inner join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.resolvedBy = a.account and b.resolution in ( "fixed") 
			where a.account in (%s) 
			group by a.account ,c.title,c.project,c.execution,c.builds 
			) tmp 
		GROUP BY account
		`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)

	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()
	for rows.Next() {
		var result QueryRdBugCarryOverResult
		err = rows.Scan(
			&result.Account,
			&result.BugCarryOverRate,
			&result.BugCarryOverRateStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, BugCarryOverRate: %f, BugCarryOverRateStandard: %f\n", result.Account, result.BugCarryOverRate, result.BugCarryOverRateStandard)

		results[result.Account] = result
	}
	return results
}

// 軟件項目版本bug遗留實際情況 有测试报告
func QueryRdBugCarryOverDetail(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryRdBugCarryOverDetailResult {
	results := map[string][]QueryRdBugCarryOverDetailResult{}
	sqlCmd := fmt.Sprintf(`
		select a.account as account, c.title as test_report, b.id as bug_id, b.title as bug_title, b.resolution as bug_resolution, b.status as bug_status
		from zt_testreport c
		inner join  zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds
		inner join zt_user a on a.account = b.resolvedBy or (a.account = b.assignedTo and b.status="active")
		where a.account in (%s) and c.end between "%s" and "%s" and c.deleted ="0" 
	`, common.AccountArrayToString(accounts), startTime, endTime)
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdBugCarryOverDetailResult
		err = rows.Scan(
			&result.Account,
			&result.TestReport,
			&result.BugId,
			&result.BugTitle,
			&result.BugResolution,
			&result.BugStatus,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, TestReport: %s, BugId: %d, BugTitle: %s, BugResolution: %s, BugStatus: %s\n", result.Account, result.TestReport, result.BugId, result.BugTitle, result.BugResolution, result.BugStatus)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 軟件研發測試單數量
func QueryCountTestTask(db *sql.DB, accounts []string, startTime, endTime string) map[string]int {
	results := map[string]int{}
	sqlCmd := fmt.Sprintf(`
		select 
			a.account , count(b.name) as testtask_num
			from zt_user a 
			inner join zt_testtask b on b.createdBy = a.account and  b.createdDate BETWEEN "%s" and "%s"
			where a.account in (%s) 
			group by a.account
		`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account string
		var count int
		err = rows.Scan(
			&account,
			&count,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, Testtask Count: %d\n", account, count)

		results[account] = count
	}
	return results

}

// select a.account,b.project,c.name,count(1) as build_fix_bug
// from zt_user a
// left join zt_bug b on b.deleted="0" and b.resolvedBy = a.account and b.resolution in ( "fixed") and assignedDate between "2024-10-01 00:00:00" and "2024-10-31 23:59:59"
// left join zt_project c on c.id = b.project
// where a.account in ("liuxiaoyan") and c.name is not null
// group by a.account ,b.project,c.name

// select a.account as account, c.name as project_name, b.id as bug_id, b.title as bug_title, b.resolution as bug_resolution, b.status as bug_status
// from zt_bug b
// left join zt_project c on c.id = b.project
// left join zt_user a on a.account = b.resolvedBy or (a.account = b.assignedTo and b.status="active")
// where a.account in ("liuxiaoyan") and b.assignedDate between "2024-10-01 00:00:00" and "2024-10-31 23:59:59" and b.deleted ="0"

// 軟件项目版本bug遗留率情况 无测试报告
func QueryRdBugCarryOverWithoutTestReport(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryRdBugCarryOverResult {
	results := map[string]QueryRdBugCarryOverResult{}
	sqlCmd := fmt.Sprintf(`
		select account ,AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) as bug_carry_over_rate,
		case when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) = 0 then "1.0"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.1 then "0.9"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.2 then "0.8"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.3 then "0.7"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.4 then "0.6"
			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.5 then "0.5"
			else "0" end as bug_carry_over_rate_standard
		from ( 
			select a.account,b.project,c.name,count(1) as build_fix_bug,
			(select count(1) from zt_bug where project = b.project and deleted = "0" and status = "active" and assignedTo = a.account and assignedDate between "%s" and "%s") as build_postponed_bugs 
			from zt_user a 
			left join zt_bug b on b.deleted="0" and b.resolvedBy = a.account and b.resolution in ( "fixed") and assignedDate between "%s" and "%s"
			left join zt_project c on c.id = b.project
			where a.account in (%s) and c.name is not null
			group by a.account ,b.project,c.name
		) tmp 
		GROUP BY account
		`, startTime, endTime, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)

	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()
	for rows.Next() {
		var result QueryRdBugCarryOverResult
		err = rows.Scan(
			&result.Account,
			&result.BugCarryOverRate,
			&result.BugCarryOverRateStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, BugCarryOverRate: %f, BugCarryOverRateStandard: %f\n", result.Account, result.BugCarryOverRate, result.BugCarryOverRateStandard)

		results[result.Account] = result
	}
	return results
}

// select account ,AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) as bug_carry_over_rate,
// 		case when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) = 0 then "1.0"
// 			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.1 then "0.9"
// 			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.2 then "0.8"
// 			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.3 then "0.7"
// 			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.4 then "0.6"
// 			when AVG(build_postponed_bugs/(build_postponed_bugs+build_fix_bug)) <= 0.5 then "0.5"
// 			else "0" end as bug_carry_over_rate_standard
// 		from ( 
// 			select a.account,b.project,c.name,count(1) as build_fix_bug,
// 			(select count(1) from zt_bug where project = b.project and deleted = "0" and status = "active" and assignedTo = a.account and assignedDate between "%s" and "%s") as build_postponed_bugs 
// 			from zt_user a 
// 			left join zt_bug b on b.deleted="0" and b.resolvedBy = a.account and b.resolution in ( "fixed") and assignedDate between "%s" and "%s"
// 			left join zt_project c on c.id = b.project
// 			where a.account in ("samy.gou","deakin.han")
// 			group by a.account ,b.project,c.name
// 		) tmp 
// 		GROUP BY account

// 項目版本bug遗留實際情況 无测试报告
func QueryRdBugCarryOverDetailWithoutTestReport(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryRdBugCarryOverDetailResultWithoutTestReport {
	results := map[string][]QueryRdBugCarryOverDetailResultWithoutTestReport{}
	sqlCmd := fmt.Sprintf(`
		select a.account as account, c.name as project_name, b.id as bug_id, b.title as bug_title, b.resolution as bug_resolution, b.status as bug_status
		from zt_bug b
		left join zt_project c on c.id = b.project
		left join zt_user a on a.account = b.resolvedBy or (a.account = b.assignedTo and b.status="active")
		where a.account in (%s) and b.assignedDate between "%s" and "%s" and b.deleted ="0"
		`, common.AccountArrayToString(accounts), startTime, endTime)
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdBugCarryOverDetailResultWithoutTestReport
		err = rows.Scan(
			&result.Account,
			&result.ProjectName,
			&result.BugId,
			&result.BugTitle,
			&result.BugResolution,
			&result.BugStatus,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		// fmt.Printf("Account: %s, ProjectName: %s, BugId: %d, BugTitle: %s, BugResolution: %s, BugStatus: %s\n", result.Account, *result.ProjectName, result.BugId, result.BugTitle, result.BugResolution, result.BugStatus)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 工时预估达成比
func QueryRdTimeEstimateRate(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryRdTimeEstimateRateResult {
	results := map[string]QueryRdTimeEstimateRateResult{}
	sqlCmd := fmt.Sprintf(`
		select a.account,AVG(b.consumed/b.estimate) as time_estimate_rate,
		case when AVG(b.consumed/b.estimate) <= 1 then "1.0"
			when AVG(b.consumed/b.estimate) > 1 and AVG(b.consumed/b.estimate) <= 1.2 then "0.8"
			when AVG(b.consumed/b.estimate) > 1.2 and AVG(b.consumed/b.estimate) <= 1.4 then "0.6"
			when AVG(b.consumed/b.estimate) > 1.4 and AVG(b.consumed/b.estimate) <= 1.6 then "0.4"
			else "0.5" end as time_estimate_level
		from zt_user a 
		inner join zt_task b on b.finishedDate between "%s" and "%s" and b.finishedBy=a.account and b.deleted="0" and b.parent=0 
		where a.account in (%s) 
		group by a.account 
		order by a.account desc
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdTimeEstimateRateResult
		err = rows.Scan(
			&result.Account,
			&result.TimeEstimateRate,
			&result.TimeEstimateRateStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, TimeEstimateRate: %f, TimeEstimateRateStandard: %f\n", result.Account, result.TimeEstimateRate, result.TimeEstimateRateStandard)

		results[result.Account] = result
	}
	return results
}

// 版本发版次数平均发版次数
func QueryRdPubTimes(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryRdPubTimesResult {
	results := map[string]QueryRdPubTimesResult{}
	sqlCmd := fmt.Sprintf(`
		select tmp2.account, AVG(tmp2.pub_times) as pub_times, 
		case when AVG(tmp2.pub_times) > 0 and AVG(tmp2.pub_times) <= 2 then 1.0
			when AVG(tmp2.pub_times) > 2 and AVG(tmp2.pub_times) <= 4 then 0.8
			when AVG(tmp2.pub_times) > 4 and AVG(tmp2.pub_times) <= 6 then 0.4
			else 0 end as pub_times_level
		from
		(
			select tmp.account, tmp.project_type, tmp.project_name, tmp.pub_times, tmp.last_pub_time
			from
			(
				select a.account as account ,c.type as project_type, c.name as project_name,count(1) as pub_times, max(b.createdDate) as last_pub_time
				from zt_user a 
				inner join zt_build b on b.builder = a.account 
				left join zt_project c on b.execution = c.id
				where a.account in (%s) 
				GROUP BY a.account,b.project,b.execution
			) tmp where tmp.last_pub_time between "%s" and "%s"
		) tmp2 group by tmp2.account
	`, common.AccountArrayToString(accounts), startTime, endTime)
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdPubTimesResult
		err = rows.Scan(
			&result.Account,
			&result.AvgPubTimes,
			&result.AvgPubTimesStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, AvgPubTimes: %f, AvgPubTimesStandard: %f\n", result.Account, result.AvgPubTimes, result.AvgPubTimesStandard)

		results[result.Account] = result
	}
	return results
}

// 版本发版次数详情
func QueryRdPubTimesDetail(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryRdPubTimesDetailResult {
	results := map[string][]QueryRdPubTimesDetailResult{}
	sqlCmd := fmt.Sprintf(`
		select tmp.account, tmp.project_type, tmp.project_name, tmp.pub_times, tmp.last_pub_time
		from
		(
			select a.account as account ,c.type as project_type, c.name as project_name,count(1) as pub_times, max(b.createdDate) as last_pub_time
			from zt_user a 
			inner join zt_build b on b.builder = a.account 
			left join zt_project c on b.execution = c.id
			where a.account in (%s) 
			GROUP BY a.account,b.project,b.execution
		) tmp where tmp.last_pub_time between "%s" and "%s"
	`, common.AccountArrayToString(accounts), startTime, endTime)
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryRdPubTimesDetailResult
		err = rows.Scan(
			&result.Account,
			&result.ProjectType,
			&result.ProjectName,
			&result.PubTimes,
			&result.LastPubTime,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, ProjectType: %s, ProjectName: %s, PubTimes: %d, LastPubTime: %s\n", result.Account, result.ProjectType, result.ProjectName, result.PubTimes, result.LastPubTime)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 测试软件项目进度达成率 完成情况
func QueryTestProjectProgressResultDetail(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryTestProjectProgressDetailResult {
	results := map[string][]QueryTestProjectProgressDetailResult{}
	sqlCmd := fmt.Sprintf(`
		select 
		a.account, c.name, b.title, c.begin, c.end, b.end as real_end
		from zt_user a
		inner join zt_testreport b on b.createdBy = a.account and b.end between "%s" and "%s" and b.deleted ="0"
		inner join zt_testtask c on c.id = b.tasks
		where a.account in (%s)
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryTestProjectProgressDetailResult
		err = rows.Scan(
			&result.Account,
			&result.TestTaskName,
			&result.TestReportTitle,
			&result.TestTaskBegin,
			&result.TestTaskEnd,
			&result.TestReportEnd,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, TestTaskName: %v, TestReportTitle: %v, TestTaskBegin: %v, TestTaskEnd: %v, TestReportEnd: %v\n",
			result.Account, result.TestTaskName, result.TestReportTitle, result.TestTaskBegin, result.TestTaskEnd, result.TestReportEnd)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 测试软件项目进度达成率
func QueryTestProjectProgress(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryTestProjectProgressResult {
	results := map[string]QueryTestProjectProgressResult{}
	sqlCmd := fmt.Sprintf(`
		select tmp.account,SUM(tmp.real_diff) as avg_real_diff, SUM(tmp.test_diff) as avg_test_diff
		from (
		select 
		a.account,c.name,b.title, c.begin,c.end,b.end as real_end, TIMESTAMPDIFF(DAY,c.begin,b.end) as real_diff, TIMESTAMPDIFF(DAY,c.begin,c.end) as test_diff
		from zt_user a
		inner join zt_testreport b on b.createdBy = a.account and b.end between "%s" and "%s" and b.deleted ="0"
		inner join zt_testtask c on c.id = b.tasks
		where a.account in (%s)
		) tmp group by tmp.account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryTestProjectProgressResult
		err = rows.Scan(
			&result.Account,
			&result.SumRealTestDiffDays,
			&result.SumTestDiffDays,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, SumRealTestDiffDays: %v, SumTestDiffDays: %v\n", result.Account, result.SumRealTestDiffDays, result.SumTestDiffDays)

		results[result.Account] = result
	}
	return results
}

// 测试软件项目有效bug率
// 1、测试报告结束时间是当月的
// 2、bug未被删除，bug关联项目属于测试报告关联项目，bug关联版本是测试报告所属版本，bug是焰海打开的，bug解决状态是转需求，延期处理和已解决的，不予解决，这些叫有效bug。
// 3、版本内所有bug，为项目与测试报告相等，并且不是指派给黄卫旗
func QueryTestValidBugRate(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryTestValidBugRateResult {
	results := map[string]QueryTestValidBugRateResult{}
	sqlCmd := fmt.Sprintf(`
		select account ,AVG(build_fix_bugs/project_build_bugs) as validate_bug_rate,
		case 
			when AVG(build_fix_bugs/project_build_bugs) >= 0.9 then 1.0
			when AVG(build_fix_bugs/project_build_bugs) >= 0.8 and AVG(build_fix_bugs/project_build_bugs) < 0.9 then 0.8
			when AVG(build_fix_bugs/project_build_bugs) >= 0.7 and AVG(build_fix_bugs/project_build_bugs) < 0.8 then 0.5
			else 0 end as validate_bug_rate_standard
		from 
		(
		select tmp.account,tmp.project,tmp.title, count(tmp.id) as build_fix_bugs,sum(DISTINCT tmp.build_bugs) as project_build_bugs
		from (
		select a.account,c.title,c.project,c.builds,b.id,
			(select count(1) from zt_bug where project = c.project and openedBuild = c.builds and deleted = "0" and assignedTo not in ("huangweiqi") and  openedBy = a.account ) as build_bugs 
		from zt_user a 
		inner join zt_testreport c on c.end between "%s" and "%s" and c.deleted ="0" 
		left join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.openedBy = a.account and b.resolution in ("tostory","postponed","willnotfix","fixed") 
		where a.account in (%s) and b.assignedTo not in ("huangweiqi")
		) tmp 
		group by tmp.account,tmp.project,tmp.title
		) tmp2 group by tmp2.account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryTestValidBugRateResult
		err = rows.Scan(
			&result.Account,
			&result.ValidBugRate,
			&result.ValidBugRateStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, ValidBugRate: %f, ValidBugRateStandard: %f\n", result.Account, result.ValidBugRate, result.ValidBugRateStandard)

		results[result.Account] = result
	}
	return results
}

// bug转需求数
func QueryTestBugToStory(db *sql.DB, accounts, storyPms []string, startTime, endTime string) map[string]QueryTestBugToStoryResult {
	results := map[string]QueryTestBugToStoryResult{}
	sqlCmd := fmt.Sprintf(`
		select tmp.account, count(1) as tostory_num
		from (
		select a.account,b.id,b.title
		from zt_user a 
		inner join zt_testreport c on c.end between "%s" and "%s" and c.deleted ="0" 
		inner join zt_bug b on b.deleted="0" and b.project = c.project and b.openedBuild = c.builds and b.openedBy = a.account and ( b.resolution in ("tostory")  or b.assignedTo in (%s))
		where a.account in (%s)
		) tmp group by tmp.account
	`, startTime, endTime, common.AccountArrayToString(storyPms), common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryTestBugToStoryResult
		err = rows.Scan(
			&result.Account,
			&result.ToStory,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, ToStory: %d\n", result.Account, result.ToStory)

		results[result.Account] = result
	}
	return results
}

// 用例发现bug率--同build(版本)下，有多少是关联case的
// "duplicate"重复bug,"bydesign"设计如此,"notrepro"未复现
func QueryTestBugCaseRate(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryTestBugCaseRateResult {
	results := map[string]QueryTestBugCaseRateResult{}
	sqlCmd := fmt.Sprintf(`
		select account ,AVG(build_case_bugs/project_build_bugs)  as case_bug_rate,
		case 
			when AVG(build_case_bugs/project_build_bugs) >= 0.9 then 1.0
			when AVG(build_case_bugs/project_build_bugs) >= 0.8 and AVG(build_case_bugs/project_build_bugs) < 0.9 then 0.8
			when AVG(build_case_bugs/project_build_bugs) >= 0.7 and AVG(build_case_bugs/project_build_bugs) < 0.8 then 0.6
			else 0.5 end as case_bug_rate_standard
		from 
		( 
			select tmp.account,tmp.project , count(tmp.id) as build_case_bugs,sum(DISTINCT tmp.build_bugs) as project_build_bugs
			from (
			select a.account,c.title,c.project,c.builds,b.id,
				(select count(1) from zt_bug where project = c.project and openedBuild = c.builds and deleted = "0" and assignedTo not in ("huangweiqi") and  openedBy = a.account ) as build_bugs 
			from zt_user a 
			inner join zt_testreport c on c.end between "%s" and "%s" and c.deleted ="0" 
			left join zt_bug b on b.deleted="0" and  b.project = c.project and  b.case <> 0 and b.openedBuild = c.builds and b.openedBy = a.account and b.resolution in ("tostory","postponed","willnotfix","fixed") 
			where a.account in (%s)  and b.assignedTo not in ("huangweiqi")
			) tmp 
			group by tmp.account,tmp.project
		) tmp2 group by tmp2.account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryTestBugCaseRateResult
		err = rows.Scan(
			&result.Account,
			&result.CaseBugRate,
			&result.CaseBugStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, CaseBugRate: %f, CaseBugStandard: %f\n", result.Account, result.CaseBugRate, result.CaseBugStandard)

		results[result.Account] = result
	}
	return results
}

// 项目软件项目进度达成率
func QueryProjectProgress(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryProjectProgressResult {
	results := map[string]QueryProjectProgressResult{}
	sqlCmd := fmt.Sprintf(`
		select account, SUM(project_diff) as sum_project_diff, SUM(real_project_diff) as sum_real_project_diff, SUM(test_diff) as sum_test_diff, SUM(real_test_diff) as sum_real_test_diff
		from ( 
			select a.account, b.id, b.name, b.type, b.begin as project_begin, b.end as project_end, b.realEnd as project_real_end,
			TIMESTAMPDIFF(DAY,b.begin,b.end) as project_diff, TIMESTAMPDIFF(DAY,b.begin,b.realEnd) as real_project_diff,
			d.begin as test_begin, d.end as test_end, c.end as test_real_end ,TIMESTAMPDIFF(DAY,d.begin,d.end) as test_diff, TIMESTAMPDIFF(DAY,d.begin,c.end) as real_test_diff
			from zt_user a 
			inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s"
			inner join zt_testreport c on (b.type="sprint" and c.execution = b.id) or (c.execution=0 and b.type="project" and c.project=b.id) 
			inner join zt_testtask d on d.id = c.tasks 
			where a.account in (%s) 
		) tmp group by tmp.account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectProgressResult
		err = rows.Scan(
			&result.Account,
			&result.SumProjectDiff,
			&result.SumRealProjectDiff,
			&result.SumTestDiff,
			&result.SumRealTestDiff,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, SumProjectDiff: %v, SumRealProjectDiff: %v, SumTestDiff: %v, SumRealTestDiff: %v\n", result.Account, result.SumProjectDiff, result.SumRealProjectDiff, result.SumTestDiff, result.SumRealTestDiff)

		results[result.Account] = result
	}
	return results
}

// 项目软件项目进度达成率 完成情况
func QueryProjectProgressDetail(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryProjectProgressDetailResult {
	results := map[string][]QueryProjectProgressDetailResult{}
	sqlCmd := fmt.Sprintf(`
	select a.account,b.id,b.name,b.type, b.begin as project_begin, b.end as project_end,b.realEnd as project_real_end,TIMESTAMPDIFF(DAY,b.end,b.realEnd) as project_diff,
			d.begin as test_begin, d.end as test_end, c.end as test_real_end
			from zt_user a
			inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s"
			inner join zt_testreport c on (b.type="sprint" and c.execution = b.id) or (c.execution=0 and b.type="project" and c.project=b.id)
			inner join zt_testtask d on d.id = c.tasks
			where a.account in (%s)
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result QueryProjectProgressDetailResult
		err = rows.Scan(
			&result.Account,
			&result.ProjectId,
			&result.ProjectName,
			&result.ProjectType,
			&result.ProjectBegin,
			&result.ProjectEnd,
			&result.ProjectRealEnd,
			&result.ProjectDiff,
			&result.TestStart,
			&result.TestEnd,
			&result.TestRealEnd,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, ProjectId: %d, ProjectName: %s, ProjectType: %s, ProjectBegin: %s, ProjectEnd: %s, ProjectRealEnd: %s, ProjectDiff: %d, TestStart: %s, TestEnd: %s\n",
		result.Account, result.ProjectId, result.ProjectName, result.ProjectType, result.ProjectBegin, result.ProjectEnd, result.ProjectRealEnd, result.ProjectDiff, result.TestStart, result.TestEnd)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 项目软件项目进度达成率，无测试报告
func QueryProjectProgressWithout(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryProjectProgressResult {
	results := map[string]QueryProjectProgressResult{}
	sqlCmd := fmt.Sprintf(`
		select account, SUM(real_project_diff) as sum_real_diff, SUM(project_diff) as sum_project_diff
		from ( 
			select a.account,b.name,b.type, b.begin as project_begin,b.end as project_end,b.realEnd as project_real_end,TIMESTAMPDIFF(DAY,b.begin,b.end) as project_diff, TIMESTAMPDIFF(DAY,b.begin,b.realEnd) as real_project_diff
			from zt_user a 
			inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s"
			where a.account in (%s) 
		) tmp group by tmp.account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectProgressResult
		err = rows.Scan(
			&result.Account,
			&result.SumRealProjectDiff,
			&result.SumProjectDiff,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, SumRealProjectDiff: %f, SumProjectDiff: %f\n", result.Account, result.SumRealProjectDiff, result.SumProjectDiff)

		results[result.Account] = result
	}
	return results
}

// 项目软件项目进度达成率 完成情况 无测试报告
func QueryProjectProgressDetailWithoutTestReport(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryProjectProgressDetailResultWithoutTestReport {
	results := map[string][]QueryProjectProgressDetailResultWithoutTestReport{}
	sqlCmd := fmt.Sprintf(`
	select a.account,b.id,b.name,b.type, b.begin, b.end as project_end,b.realEnd as project_real_end
			from zt_user a
			inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s"
			where a.account in (%s)
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result QueryProjectProgressDetailResultWithoutTestReport
		err = rows.Scan(
			&result.Account,
			&result.ProjectId,
			&result.ProjectName,
			&result.ProjectType,
			&result.ProjectBegin,
			&result.ProjectEnd,
			&result.ProjectRealEnd,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, ProjectId: %d, ProjectName: %s, ProjectType: %s, ProjectBegin: %v, ProjectEnd: %s, ProjectRealEnd: %s\n", result.Account, result.ProjectId, result.ProjectName, result.ProjectType, result.ProjectBegin, result.ProjectEnd, result.ProjectRealEnd)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 项目成果完成率,不需要关注执行，只需要看项目需求完成度，因为有执行一定有项目
func QueryProjectCompleteRate(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryProjectCompleteRateResult {
	results := map[string]QueryProjectCompleteRateResult{}
	sqlCmd := fmt.Sprintf(`
		select account,AVG(devleoped_num/project_storys) as complete_rate,
		case when AVG(devleoped_num/project_storys) >= 1 then 1.0
			when AVG(devleoped_num/project_storys) >= 0.85 and AVG(devleoped_num/project_storys) < 1 then 0.8
			when AVG(devleoped_num/project_storys) >= 0.75 and AVG(devleoped_num/project_storys) < 0.85 then 0.6
			when AVG(devleoped_num/project_storys) >= 0.65 and AVG(devleoped_num/project_storys) < 0.75 then 0.4
			else 0 end as complete_rate_standard
		from ( 
			select a.account,b.id,b.name,count(1) as devleoped_num,
			(select count(1) from zt_projectstory ztp 
				inner join zt_story zts on zts.id=ztp.story and zts.deleted="0" 
				where project = b.id) as project_storys 
			from zt_user a 
			inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s" and b.project = 0 
			inner join zt_projectstory d on d.project = b.id
			inner join zt_story c on c.id= d.story and c.deleted = "0" and c.stage not in ("waiting","planned","projected","developing") 
			where a.account in (%s) 
			group by a.account,b.id 
		) tmp group by account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectCompleteRateResult
		err = rows.Scan(
			&result.Account,
			&result.CompleteRate,
			&result.CompleteRateStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, CompleteRate: %f, CompleteRateStandard: %f\n", result.Account, result.CompleteRate, result.CompleteRateStandard)

		results[result.Account] = result
	}
	return results
}

// 项目成果完成率,不需要关注执行，只需要看项目需求完成度，因为有执行一定有项目
func QueryProjectCompleteRateWithout(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryProjectCompleteRateResult {
	results := map[string]QueryProjectCompleteRateResult{}
	sqlCmd := fmt.Sprintf(`
		select account,AVG(devleoped_num/project_storys) as complete_rate,
		case when AVG(devleoped_num/project_storys) >= 0.95 then 1.0
			when AVG(devleoped_num/project_storys) >= 0.9 and AVG(devleoped_num/project_storys) < 0.95 then 1.0
			when AVG(devleoped_num/project_storys) >= 0.8 and AVG(devleoped_num/project_storys) < 0.9 then 0.8
			when AVG(devleoped_num/project_storys) >= 0.7 and AVG(devleoped_num/project_storys) < 0.8 then 0.5
			else 0 end as complete_rate_standard
		from ( 
			select a.account,b.id,b.name,count(1) as devleoped_num,
			(select count(1) from zt_projectstory ztp 
				inner join zt_story zts on zts.id=ztp.story and zts.deleted="0" 
				where project = b.id) as project_storys 
			from zt_user a 
			inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s" and b.project = 0 
			inner join zt_projectstory d on d.project = b.id
			inner join zt_story c on c.id= d.story and c.deleted = "0" and c.stage not in ("waiting","planned","projected","developing") 
			where a.account in (%s) 
			group by a.account,b.id 
		) tmp group by account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectCompleteRateResult
		err = rows.Scan(
			&result.Account,
			&result.CompleteRate,
			&result.CompleteRateStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, CompleteRate: %f, CompleteRateStandard: %f\n", result.Account, result.CompleteRate, result.CompleteRateStandard)

		results[result.Account] = result
	}
	return results
}

// 项目成果完成率,完成情况
func QueryProjectCompleteRateDetail(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryProjectCompleteRateDetailResult {
	results := map[string][]QueryProjectCompleteRateDetailResult{}
	sqlCmd := fmt.Sprintf(`
		select account,name as project_name,devleoped_num/project_storys as complete_rate
		from (
			select a.account,b.id,b.name,count(1) as devleoped_num,
			(select count(1) from zt_projectstory ztp 
				inner join zt_story zts on zts.id=ztp.story and zts.deleted="0" 
				where project = b.id) as project_storys 
			from zt_user a 
			inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s" and b.project = 0 
			inner join zt_projectstory d on d.project = b.id
			inner join zt_story c on c.id= d.story and c.deleted = "0" and c.stage not in ("waiting","planned","projected","developing") 
			where a.account in (%s) 
			group by a.account,b.id 
		) tmp
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectCompleteRateDetailResult
		err = rows.Scan(
			&result.Account,
			&result.ProjectName,
			&result.CompleteRate,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, ProjectName: %s, CompleteRate: %f\n", result.Account, result.ProjectName, result.CompleteRate)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 项目规划需求数
func QueryProjectStoryNum(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryProjectStoryNumResult {
	results := map[string][]QueryProjectStoryNumResult{}
	sqlCmd := fmt.Sprintf(`
		with recursive cte as ( 
		select b.id,d.story 
		from zt_project b 
		inner join zt_projectstory d on d.project = b.id 
		where b.deleted="0" and b.realEnd between "%s" and "%s" and b.project = 0 
		) 

		select b.openedBy,b.stage ,count(1) 
		from zt_story b 
		inner join cte a on a.story = b.id 
		where b.openedBy in (%s) and b.deleted = "0"
		group by b.openedBy,b.stage
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectStoryNumResult
		err = rows.Scan(
			&result.Account,
			&result.Stage,
			&result.StoryNum,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, Stage: %s, StoryNum: %d\n", result.Account, result.Stage, result.StoryNum)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

// 预估承诺完成率，只看项目，和最后一个测试报告
func QueryProjectEstimateRate(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryProjectEstimateRateResult {
	results := map[string]QueryProjectEstimateRateResult{}
	sqlCmd := fmt.Sprintf(`
		select tmp2.account,AVG(tmp2.diff_days) as diff_days,
		case when AVG(tmp2.diff_days) <= -3 then 1.2
			when AVG(tmp2.diff_days) > -3 and AVG(tmp2.diff_days) <= 0 then 1.0
			when AVG(tmp2.diff_days) > 0 and AVG(tmp2.diff_days) <= 2 then 0.8
			when AVG(tmp2.diff_days) > 2 and AVG(tmp2.diff_days) <= 4 then 0.6
			when AVG(tmp2.diff_days) > 4 and AVG(tmp2.diff_days) <= 5 then 0.5
			else 0 end as progress_standard
		from ( 
			select account,name,plan,plan_end,testreport_end,TIMESTAMPDIFF(DAY,plan_end,testreport_end) as diff_days 
			from ( 
				select a.account,b.name,REPLACE(c.plan,",","") as plan,d.end as plan_end,
				(select end from zt_testreport a where id = (select max(id) from zt_testreport where project = b.id and product = d.product ) ) as testreport_end 
				from zt_user a 
				inner join zt_project b on b.PM = a.account and b.deleted="0" and b.realEnd between "%s" and "%s"
				inner join zt_projectproduct c on c.project = b.id 
				inner join zt_productplan d on d.id=REPLACE(c.plan,",","") 
				where a.account in (%s) 
			) tmp where testreport_end is not null
		)tmp2 
		group by account
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectEstimateRateResult
		err = rows.Scan(
			&result.Account,
			&result.DiffDays,
			&result.ProgressStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, AvgDiffDays: %f, ProgressStandard: %f\n", result.Account, result.DiffDays, result.ProgressStandard)

		results[result.Account] = result
	}
	return results
}

// 项目工时预估准确率
func QueryProjectTimeEstimateRate(db *sql.DB, accounts []string, startTime, endTime string) map[string]QueryProjectTimeEstimateRateResult {
	results := map[string]QueryProjectTimeEstimateRateResult{}
	sqlCmd := fmt.Sprintf(`
		select 
		tmp.account as account,(100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) as estimate_rate,
		case when (100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) > 90 then 1
		when (100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) <= 90 and (100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) > 80 then 0.8
		when (100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) <= 80 and (100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) > 70 then 0.6
		when (100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) <= 70 and (100 - abs(avg(tmp.story_consumed / tmp.estimate) - 1)*100) > 60 then 0.4
		else 0 end as estimate_standard
		from (
		select
		a.account,b.id,b.title,b.estimate,sum(c.consumed) as story_consumed
		from zt_user a
		inner  join zt_story b on b.openedBy = a.account and b.deleted = "0" 
		inner join zt_task c on c.story = b.id  and c.finishedDate BETWEEN "%s" and "%s" and c.deleted = "0" and c.parent= 0
		where a.account in (%s)
		GROUP BY a.account,b.id,b.estimate
		) tmp
		GROUP BY tmp.account;
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectTimeEstimateRateResult
		err = rows.Scan(
			&result.Account,
			&result.TimeEstimateRate,
			&result.TimeEstimateStandard,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, TimeEstimateRate: %f, TimeEstimateStandard: %f\n", result.Account, result.TimeEstimateRate, result.TimeEstimateStandard)

		results[result.Account] = result
	}
	return results
}

// 项目工时预估准确率 完成情况
func QueryProjectTimeEstimateRateDetail(db *sql.DB, accounts []string, startTime, endTime string) map[string][]QueryProjectTimeEstimateRateResultDetail {
	results := map[string][]QueryProjectTimeEstimateRateResultDetail{}
	sqlCmd := fmt.Sprintf(`
		select
		a.account,b.id,b.title,b.estimate,sum(c.consumed) as story_consumed, 100 - (1 - abs(sum(c.consumed)/b.estimate))*100 as estimate_rate
		from zt_user a
		inner  join zt_story b on b.openedBy = a.account and b.deleted = "0" 
		inner join zt_task c on c.story = b.id  and c.finishedDate BETWEEN "%s" and "%s" and c.deleted = "0" and c.parent= 0
		where a.account in (%s)
		GROUP BY a.account,b.id,b.estimate;
	`, startTime, endTime, common.AccountArrayToString(accounts))
	fmt.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result QueryProjectTimeEstimateRateResultDetail
		err = rows.Scan(
			&result.Account,
			&result.StoryId,
			&result.Title,
			&result.Estimate,
			&result.StoryConsumed,
			&result.EstimateRate,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		fmt.Printf("Account: %s, StoryId: %v, Title: %v, Estimate: %v, StoryConsumed: %v, EstimateRate: %v\n", result.Account, result.StoryId, result.Title, result.Estimate, result.StoryConsumed, result.EstimateRate)

		results[result.Account] = append(results[result.Account], result)
	}
	return results
}

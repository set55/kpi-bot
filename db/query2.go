package dbQuery

import (
	"database/sql"
	"fmt"
	// "kpi-bot/common"
	"log"
)

const (
	BugSeverity1 = 1 // 致命
	BugSeverity2 = 2 // 严重
	BugSeverity3 = 3 // 一般
	BugSeverity4 = 4 // 轻微
)

type (
	RdTask struct {
		StoryId       int
		StoryTitle    string
		TaskId        int
		TaskName      string
		FinishedBy    string
		StoryEstimate float64
		TaskConsumed  float64
		FinishedDate  string
	}

	RdBug struct {
		BugId         int
		BugTitle      string
		BugStatus     string
		BugResolution string
		BugSeverity  int
	}

	Project struct {
		Account   string
		Root      int   // zt_team.root is zt_project.id 项目id(可能是冲刺也可能是项目)
		Name      string
		Begin     *string
		End       *string
		RealBegan *string
		RealEnd   *string
	}

	Account struct {
		Account string
		RealName string
	}
)

func QueryRdTasks(db *sql.DB, account, startTime, endTime string) []RdTask {
	results := []RdTask{}
	sqlCmd := fmt.Sprintf(`
		select zt_story.id , zt_story.title, zt_task.id, zt_task.name, zt_task.finishedBy, zt_story.estimate, zt_task.consumed, zt_task.finishedDate
		from zt_task
		inner join zt_story on zt_task.story=zt_story.id
		where finishedBy='%s' and finishedDate > '%s' and finishedDate <= '%s';	
	`, account, startTime, endTime)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var result RdTask
		err := rows.Scan(&result.StoryId, &result.StoryTitle, &result.TaskId, &result.TaskName, &result.FinishedBy, &result.StoryEstimate, &result.TaskConsumed, &result.FinishedDate)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		results = append(results, result)
	}
	return results
}

func RdBugs(db *sql.DB, account string, endDate string) []RdBug {
	result := []RdBug{}
	sqlCmd := fmt.Sprintf(`
		select id, title, status, resolution, severity from zt_bug 
		where assignedTo='%s' and status="active" and deleted='0' and openedDate <= '%s';
	`, account, endDate)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var res RdBug
		err := rows.Scan(&res.BugId, &res.BugTitle, &res.BugStatus, &res.BugResolution, &res.BugSeverity)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result = append(result, res)
	}
	return result
}

func RdBugsApp(db *sql.DB, account string, startDate, endDate string) []RdBug {
	result := []RdBug{}
	sqlCmd := fmt.Sprintf(`
		select id, title, status, resolution, severity from zt_bug 
		where assignedTo='%s' and status="active" and deleted='0' and ((activatedDate >= '%s' and activatedDate <= '%s') or (openedDate >= '%s' and openedDate <= '%s') or (assignedDate >= '%s' and assignedDate <= '%s'));
	`, account, startDate, endDate, startDate, endDate, startDate, endDate)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var res RdBug
		err := rows.Scan(&res.BugId, &res.BugTitle, &res.BugStatus, &res.BugResolution, &res.BugSeverity)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result = append(result, res)
	}
	return result
}

func AppRdActiveBugs(db *sql.DB, account string) []RdBug {
	result := []RdBug{}
	sqlCmd := fmt.Sprintf(`
		select id, title, status, resolution from zt_bug
		where assignedTo='%s' and status='active' and deleted='0';
	`, account)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var res RdBug
		err := rows.Scan(&res.BugId, &res.BugTitle, &res.BugStatus, &res.BugResolution)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result = append(result, res)
	}
	return result
}

func AppRdResolvedBugs(db *sql.DB, account string) []RdBug {
	result := []RdBug{}
	sqlCmd := fmt.Sprintf(`
		select id, title, status, resolution from zt_bug
		where resolvedBy='%s' and (status='resolved' or status='closed') and deleted='0';
	`, account)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var res RdBug
		err := rows.Scan(&res.BugId, &res.BugTitle, &res.BugStatus, &res.BugResolution)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result = append(result, res)
	}
	return result
}

func QueryRdProjects(db *sql.DB, account, startTime, endTime string) []Project {
	results := []Project{}
	// todo: zt_team.type = 'execution' is sprint, type = 'project' is project, which one should we use?
	sqlCmd := fmt.Sprintf(`
		select zt_team.account, zt_team.root, zt_project.name, zt_project.begin, zt_project.end, zt_project.realBegan, zt_project.realEnd 
		from zt_team
		inner join zt_project on zt_team.root=zt_project.id
		where zt_team.account='%s' and zt_team.type='execution' and zt_project.realEnd > '%s' and zt_project.realEnd <= '%s';
	`, account, startTime, endTime)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var result Project
		err := rows.Scan(&result.Account, &result.Root, &result.Name, &result.Begin, &result.End, &result.RealBegan, &result.RealEnd)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		results = append(results, result)
	}
	return results
}

func QueryAccount(db *sql.DB, account string) Account {
	sqlCmd := fmt.Sprintf(`
		select account, realname from zt_user 
		where account='%s';
	`, account)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	var result Account
	for rows.Next() {
		err := rows.Scan(&result.Account, &result.RealName)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
	}
	return result
}

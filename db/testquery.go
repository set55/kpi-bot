package dbQuery

import (
	"database/sql"
	"fmt"
	// "kpi-bot/common"
	"log"
)

type (
	TestReport struct {
		ReportId        int
		TaskId          int
		ReportTitle     string
		TaskName        string
		ReportCreatedAt string
		ReportCreator   string
		TaskBegin       string
		TaskEnd         string
	}

	TestBug struct {
		BugId         int
		BugTitle      string
		BugCreator    string
		BugStatus     string
		BugResolution string
	}
)

func QueryTestReport(db *sql.DB, account, startTime, endTime string) []TestReport {
	results := []TestReport{}
	sqlCmd := fmt.Sprintf(`
		select 
		zt_testreport.id as reportId,
		zt_testtask.id as taskId,
		zt_testreport.title as reportTitle,
		zt_testtask.name as taskName,
		zt_testreport.createdDate as reportCreatedAt,
		zt_testreport.createdBy as reportCreator,
		zt_testtask.begin as taskBegin,
		zt_testtask.end as taskEnd
		from zt_testreport
		inner join zt_testtask on zt_testreport.objectID = zt_testtask.id
		where zt_testreport.createdBy='%s' 
		and zt_testreport.objectType='testtask'
		and zt_testreport.createdDate >= '%s'
		and zt_testreport.createdDate <= '%s';		
	`, account, startTime, endTime)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var result TestReport
		err := rows.Scan(&result.ReportId, &result.TaskId, &result.ReportTitle, &result.TaskName, &result.ReportCreatedAt, &result.ReportCreator, &result.TaskBegin, &result.TaskEnd)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		results = append(results, result)
	}
	return results
}

func TestBugs(db *sql.DB, account, startTime, endTime string) []TestBug {
	result := []TestBug{}
	sqlCmd := fmt.Sprintf(`
		select
		zt_bug.id as bugId,
		zt_bug.title as bugTitle,
		zt_bug.openedBy as bugCreator,
		zt_bug.status as bugStatus,
		zt_bug.resolution as bugResolution
		from zt_bug
		where zt_bug.openedBy='%s'
		and zt_bug.openedDate >= '%s'
		and zt_bug.openedDate <= '%s';
	`, account, startTime, endTime)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var res TestBug
		err := rows.Scan(&res.BugId, &res.BugTitle, &res.BugCreator, &res.BugStatus, &res.BugResolution)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result = append(result, res)
	}
	return result
}

func StoryBugs(db *sql.DB, account, startTime, endTime string) []TestBug {
	result := []TestBug{}
	sqlCmd := fmt.Sprintf(`
		select
		zt_bug.id as bugId,
		zt_bug.title as bugTitle,
		zt_bug.openedBy as bugCreator,
		zt_bug.status as bugStatus,
		zt_bug.resolution as bugResolution
		from zt_bug
		where zt_bug.openedBy='%s'
		and zt_bug.resolution='tostory'
		and zt_bug.resolvedDate >= '%s'
		and zt_bug.resolvedDate <= '%s';
	`, account, startTime, endTime)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var res TestBug
		err := rows.Scan(&res.BugId, &res.BugTitle, &res.BugCreator, &res.BugStatus, &res.BugResolution)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result = append(result, res)
	}
	return result
}



package rd

import (
	"database/sql"
	"fmt"
	"kpi-bot/common"
	dbQuery "kpi-bot/db"

	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

type (
	RdKpi2 struct {
		Account string
		Db      *sql.DB
		// 起始时间
		StartTime string
		// 结束时间
		EndTime string
		Coefficient RdCoefficient
	}

	RdCoefficient struct {
		PROJECT_PROGRESS_STANDARD float64
		DELAY_DAYS_SCORE 	 float64
		STORY_STANDARD 	 float64
		STORY_BASE_TIME 	 float64
		STORY_BASE_SCORE 	 float64
		BUG_CARRY_OVER_STANDARD float64
		BUG_ONE_SCORE 	 float64
		TOP_COEFFICIENT 	 float64
		SECOND_COEFFICIENT 	 float64
		THIRD_COEFFICIENT 	 float64
	}

	RdKpiResult2 struct {
		ProjectDetail string
		ProjectGrade  float64
		ProjectDatas  []ProjectData
		StoryDetail   string
		StoryGrade    float64
		StoryDatas    []StoryData
		BugDetail     string
		BugGrade      float64
		BugDatas      []BugData
		TimeDetail    string
		TimeGrade     float64
		TotalGrade    float64
		StartTime     string
		EndTime       string
		AccountName   string
		Coefficient   float64
	}

	ProjectData struct {
		Root      int
		Name      string
		RealEnd   *string
		End       *string
		DelayDays int
	}

	StoryData struct {
		StoryId       int
		StoryTitle    string
		StoryEstimate float64
		TaskConsumed  float64
		StoryBase     float64
	}

	BugData struct {
		BugId         int
		BugTitle      string
		BugStatus     string
		BugResolution string
	}
)

// NewRdKpi 创建一个研发KPI对象
func NewRdKpi2(db *sql.DB, account, startTime, endTime string, coefficient RdCoefficient) *RdKpi2 {
	return &RdKpi2{
		Account:   account,
		Db:        db,
		StartTime: startTime,
		EndTime:   endTime,
		Coefficient: coefficient,
	}
}

// GetRdKpiGrade 获取研发KPI信息
func (l *RdKpi2) GetRdKpiGrade2() (result RdKpiResult2) {
	result.StartTime = l.StartTime
	result.EndTime = l.EndTime
	result.AccountName = l.Account
	account := dbQuery.QueryAccount(l.Db, l.Account)
	if account.RealName != "" {
		result.AccountName = account.RealName
	}
	projects := dbQuery.QueryRdProjects(l.Db, l.Account, l.StartTime, l.EndTime)
	delayDays := 0
	// 项目进度完成情况
	for _, project := range projects {
		ddays := common.CalculateDelayDays(*project.RealEnd, *project.End)
		delayDays += ddays
		result.ProjectDatas = append(result.ProjectDatas, ProjectData{
			Root:      project.Root,
			Name:      project.Name,
			RealEnd:   project.RealEnd,
			End:       project.End,
			DelayDays: ddays,
		})
		// result.ProjectDetail += fmt.Sprintf("项目名称: %s 延时天数: %d\n\n",project.Name, ddays)
	}
	// 项目进度分数
	result.ProjectGrade = float64(int(l.Coefficient.PROJECT_PROGRESS_STANDARD) - delayDays* int(l.Coefficient.DELAY_DAYS_SCORE))
	if result.ProjectGrade < 0 {
		result.ProjectGrade = 0
	}
	result.ProjectDetail = fmt.Sprintf("總延時天數: %d\n\n", delayDays) + result.ProjectDetail

	// 需求完成情况
	tasks := dbQuery.QueryRdTasks(l.Db, l.Account, l.StartTime, l.EndTime)
	type (
		StoryMap struct {
			StoryId       int
			StoryTitle    string
			StoryEstimate float64
			TaskConsumed  float64
		}
	)
	storyMap := map[int]StoryMap{}
	// 将任务的需求整理出来
	for _, task := range tasks {
		if story, ok := storyMap[task.StoryId]; !ok {
			storyMap[task.StoryId] = StoryMap{
				StoryId:       task.StoryId,
				StoryTitle:    task.StoryTitle,
				StoryEstimate: task.StoryEstimate,
				TaskConsumed:  task.TaskConsumed,
			}
		} else {
			story.TaskConsumed += task.TaskConsumed
			storyMap[task.StoryId] = story
		}
	}
	totalStoryHours := 0.0 // 预估工时
	totalTaskHours := 0.0 // 实际工时
	totalStoryCount := 0
	for _, story := range storyMap {
		// 需求基础分
		storyBase := l.GetStoryBase2(story.StoryEstimate, story.TaskConsumed)
		result.StoryDatas = append(result.StoryDatas, StoryData{
			StoryId:       story.StoryId,
			StoryTitle:    story.StoryTitle,
			StoryEstimate: story.StoryEstimate,
			TaskConsumed:  story.TaskConsumed,
			StoryBase:     storyBase,
		})
		// result.StoryDetail += fmt.Sprintf("需求id: %d, 需求名称: %s, 预估工时: %f, 实际工时: %f, 需求基础分: %f\n\n", story.StoryId, story.StoryTitle, story.StoryEstimate, story.TaskConsumed, storyBase)
		result.StoryGrade += storyBase
		totalStoryHours += story.StoryEstimate
		totalTaskHours += story.TaskConsumed
		totalStoryCount++
	}
	result.StoryDetail = fmt.Sprintf("总需求数: %d\n\n总预估工时: %.2f\n\n总实际工时: %.2f\n\n", totalStoryCount, totalStoryHours, totalTaskHours)

	// bug遗留率
	bugs := dbQuery.RdBugs(l.Db, l.Account, l.EndTime)
	deleteBugScore := 0.0
	for _, bug := range bugs {
		// result.BugDetail += fmt.Sprintf("bug id: %d, bug标题: %s, bug状态: %s, bug解决情况: %s\n\n", bug.BugId, bug.BugTitle, bug.BugStatus, bug.BugResolution)
		result.BugDatas = append(result.BugDatas, BugData{
			BugId:         bug.BugId,
			BugTitle:      bug.BugTitle,
			BugStatus:     bug.BugStatus,
			BugResolution: bug.BugResolution,
		})
		deleteBugScore += l.Coefficient.BUG_ONE_SCORE
	}
	result.BugDetail = fmt.Sprintf("总bug数: %d\n\n", len(bugs))
	result.BugGrade = float64(int(l.Coefficient.BUG_CARRY_OVER_STANDARD) - int(deleteBugScore))
	if result.BugGrade < 0 {
		result.BugGrade = 0
	}

	result.TotalGrade = result.ProjectGrade + result.StoryGrade + result.BugGrade

	result.Coefficient = l.GetKpiGradeStandard2(result.TotalGrade)
	return result
}

func (l *RdKpi2) MakeRdReport(department, career, dir, boss, path string) error {
	data := l.GetRdKpiGrade2()

	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("open file fail: %v", err)
	}
	defer f.Close()

	// Parse the date string into a time.Time object
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, data.StartTime)
	if err != nil {
		return fmt.Errorf("parse time err: %v", err)
	}

	// Extract the year and month
	year := t.Year()
	month := t.Month()

	// Sheet1 A1. 标题
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("%s %s岗%v年%v月绩效考核表", department, career, year, int(month)))

	// Sheet1 A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门："+department)

	// Sheet1 E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", data.AccountName))

	// Sheet1 F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", boss))

	// Sheet1 G4. 项目进度延时率 完成情况

	f.SetCellValue("Sheet1", "G4", data.ProjectDetail)

	// Sheet1 H4. 项目进度延时率 最终得分
	f.SetCellValue("Sheet1", "H4", data.ProjectGrade)

	// Sheet1 G5. 需求达成率 完成情况
	f.SetCellValue("Sheet1", "G5", data.StoryDetail)

	// Sheet1 H5. 需求达成率 最终得分
	f.SetCellValue("Sheet1", "H5", data.StoryGrade)

	// Sheet1 G6. bug遗留率 完成情况
	f.SetCellValue("Sheet1", "G6", data.BugDetail)

	// Sheet1 H6. bug遗留率 最终得分
	f.SetCellValue("Sheet1", "H6", data.BugGrade)

	// Sheet1 H9. 总分数
	f.SetCellValue("Sheet1", "H9", data.TotalGrade)

	// Sheet1 G11. 绩效基数
	f.SetCellValue("Sheet1", "G11", data.Coefficient)

	
	// Shee2 A1 Project
	f.NewSheet("Sheet2")
	f.SetColWidth("Sheet2", "A", "E", 20)
	f.SetColWidth("Sheet2", "B", "B", 100)
	f.SetCellValue("Sheet2", "A1", "项目")
	rowNum := 2
	f.SetCellValue("Sheet2", "A2", "项目id")
	f.SetCellValue("Sheet2", "B2", "项目名称")
	f.SetCellValue("Sheet2", "C2", "实际结束时间")
	f.SetCellValue("Sheet2", "D2", "预计结束时间")
	f.SetCellValue("Sheet2", "E2", "延时天数")
	for _, project := range data.ProjectDatas {
		rowNum++
		f.SetCellValue("Sheet2", fmt.Sprintf("A%v", rowNum), project.Root)
		f.SetCellValue("Sheet2", fmt.Sprintf("B%v", rowNum), project.Name)
		if project.RealEnd != nil {
			f.SetCellValue("Sheet2", fmt.Sprintf("C%v", rowNum), *project.RealEnd)
		}
		if project.End != nil {
			f.SetCellValue("Sheet2", fmt.Sprintf("D%v", rowNum), *project.End)
		}
		f.SetCellValue("Sheet2", fmt.Sprintf("E%v", rowNum), project.DelayDays)
	}

	rowNum++
	f.SetCellValue("Sheet2", fmt.Sprintf("A%v", rowNum), "需求")
	rowNum++
	f.SetCellValue("Sheet2", fmt.Sprintf("A%v", rowNum), "需求id")
	f.SetCellValue("Sheet2", fmt.Sprintf("B%v", rowNum), "需求名称")
	f.SetCellValue("Sheet2", fmt.Sprintf("C%v", rowNum), "预估工时")
	f.SetCellValue("Sheet2", fmt.Sprintf("D%v", rowNum), "实际工时")
	f.SetCellValue("Sheet2", fmt.Sprintf("E%v", rowNum), "需求分")

	for _, story := range data.StoryDatas {
		rowNum++
		f.SetCellValue("Sheet2", fmt.Sprintf("A%v", rowNum), story.StoryId)
		f.SetCellValue("Sheet2", fmt.Sprintf("B%v", rowNum), story.StoryTitle)
		f.SetCellValue("Sheet2", fmt.Sprintf("C%v", rowNum), story.StoryEstimate)
		f.SetCellValue("Sheet2", fmt.Sprintf("D%v", rowNum), story.TaskConsumed)
		f.SetCellValue("Sheet2", fmt.Sprintf("E%v", rowNum), story.StoryBase)
	}

	rowNum++
	f.SetCellValue("Sheet2", fmt.Sprintf("A%v", rowNum), "Bug")
	rowNum++
	f.SetCellValue("Sheet2", fmt.Sprintf("A%v", rowNum), "Bug id")
	f.SetCellValue("Sheet2", fmt.Sprintf("B%v", rowNum), "Bug 标题")
	f.SetCellValue("Sheet2", fmt.Sprintf("C%v", rowNum), "Bug 状态")
	f.SetCellValue("Sheet2", fmt.Sprintf("D%v", rowNum), "Bug 解决方案")

	for _, bug := range data.BugDatas {
		rowNum++
		f.SetCellValue("Sheet2", fmt.Sprintf("A%v", rowNum), bug.BugId)
		f.SetCellValue("Sheet2", fmt.Sprintf("B%v", rowNum), bug.BugTitle)
		f.SetCellValue("Sheet2", fmt.Sprintf("C%v", rowNum), bug.BugStatus)
		f.SetCellValue("Sheet2", fmt.Sprintf("D%v", rowNum), bug.BugResolution)
	}
	

	

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v/%s", year, int(month), dir)
	filePath := fmt.Sprintf("./export/%v-%v/%s/%v-%v-绩效考核模板-研发-%v.xlsx", year, int(month), dir, year, int(month), data.AccountName)
	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Create the folder if it does not exist
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating folder: %v", err)
		}
		fmt.Println("Folder created successfully.")
	} else {
		fmt.Println("Folder already exists.")
	}

	// Save the modified file
	if err = f.SaveAs(filePath); err != nil {
		return fmt.Errorf("save as %v, err: %v", filePath, err)
	}
	return nil
}

// GetStoryBase 获取需求基础分
func (l *RdKpi2) GetStoryBase2(estimate, consumed float64) float64 {
	return (estimate / l.Coefficient.STORY_BASE_TIME) * l.Coefficient.STORY_BASE_SCORE
}

func (l *RdKpi2) GetKpiGradeStandard2(totalGrade float64) float64 {
	if totalGrade >= 90 {
		return l.Coefficient.TOP_COEFFICIENT
	} else if totalGrade < 90 && totalGrade >= 70 {
		return l.Coefficient.SECOND_COEFFICIENT
	} else if totalGrade < 70 && totalGrade >= 60 {
		return l.Coefficient.THIRD_COEFFICIENT
	}
	return 0
}

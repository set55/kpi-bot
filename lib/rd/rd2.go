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

const (
	// 项目进度延时率 分值
	PROJECT_PROGRESS_STANDARD2 = 40
	DELAY_DAYS_SCORE           = 4

	// 需求完成率 分值
	STORY_STANDARD2 = 30

	// 需求基础 分值
	STORY_BASE_TIME2  = 0.1   // 小时
	STORY_BASE_SCORE2 = 0.021 // 分值

	// bug遗留率 分值
	BUG_CARRY_OVER_STANDARD2 = 30
	BUG_ONE_SCORE            = 1

	// 系数
	TOP_COEFFICIENT2    = 1.2
	SECOND_COEFFICIENT2 = 1.0
	THIRD_COEFFICIENT2  = 0.8
)

type (
	RdKpi2 struct {
		Account string
		Db      *sql.DB
		// 起始时间
		StartTime string
		// 结束时间
		EndTime string
	}

	RdKpiResult2 struct {
		ProjectDetail string
		ProjectGrade  float64
		StoryDetail   string
		StoryGrade    float64
		BugDetail     string
		BugGrade      float64
		TimeDetail    string
		TimeGrade     float64
		TotalGrade    float64
		StartTime     string
		EndTime       string
		AccountName string
		Coefficient float64
	}
)

// NewRdKpi 创建一个研发KPI对象
func NewRdKpi2(db *sql.DB, account, startTime, endTime string) *RdKpi2 {
	return &RdKpi2{
		Account:   account,
		Db:        db,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// GetRdKpiGrade 获取研发KPI信息
func (l *RdKpi2) GetRdKpiGrade2() (result RdKpiResult2) {
	result.StartTime = l.StartTime
	result.EndTime = l.EndTime
	result.AccountName = common.AccountToName(l.Account)
	projects := dbQuery.QueryRdProjects(l.Db, l.Account, l.StartTime, l.EndTime)
	delayDays := 0
	// 项目进度完成情况
	for _, project := range projects {
		delayDays += common.CalculateDelayDays(*project.RealEnd, *project.End)
		result.ProjectDetail += fmt.Sprintf("id: %d 项目名称: %s, 实际结束时间: %s, 计划结束时间: %s, 延时天数: %d\n\n", project.Root, project.Name, *project.RealEnd, *project.End, common.CalculateDelayDays(*project.RealEnd, *project.End))
	}
	// 项目进度分数
	result.ProjectGrade = float64(PROJECT_PROGRESS_STANDARD2 - delayDays*DELAY_DAYS_SCORE)
	if result.ProjectGrade < 0 {
		result.ProjectGrade = 0
	}

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

	for _, story := range storyMap {
		// 需求基础分
		storyBase := GetStoryBase2(story.StoryEstimate, story.TaskConsumed)
		result.StoryDetail += fmt.Sprintf("需求id: %d, 需求名称: %s, 预估工时: %f, 实际工时: %f, 需求基础分: %f\n\n", story.StoryId, story.StoryTitle, story.StoryEstimate, story.TaskConsumed, storyBase)
		result.StoryGrade += storyBase
	}

	// bug遗留率
	bugs := dbQuery.RdBugs(l.Db, l.Account)
	deleteBugScore := 0.0
	for _, bug := range bugs {
		result.BugDetail += fmt.Sprintf("bug id: %d, bug标题: %s, bug状态: %s, bug解决情况: %s\n\n", bug.BugId, bug.BugTitle, bug.BugStatus, bug.BugResolution)
		deleteBugScore += BUG_ONE_SCORE
	}
	result.BugGrade = float64(BUG_CARRY_OVER_STANDARD2 - int(deleteBugScore))

	result.TotalGrade = result.ProjectGrade + result.StoryGrade + result.BugGrade

	result.Coefficient = GetKpiGradeStandard2(result.TotalGrade)
	return result
}

func (l *RdKpi2) MakeRdReport(path string) error {
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

	// A1. 标题
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 研发工程师岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", data.AccountName))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度延时率 完成情况
	
	f.SetCellValue("Sheet1", "G4", data.ProjectDetail)

	// H4. 项目进度延时率 最终得分
	f.SetCellValue("Sheet1", "H4", data.ProjectGrade)

	// G5. 需求达成率 完成情况
	f.SetCellValue("Sheet1", "G5", data.StoryDetail)

	// H5. 需求达成率 最终得分
	f.SetCellValue("Sheet1", "H5", data.StoryGrade)

	// G6. bug遗留率 完成情况
	f.SetCellValue("Sheet1", "G6", data.BugDetail)

	// H6. bug遗留率 最终得分
	f.SetCellValue("Sheet1", "H6", data.BugGrade)

	// H11. 总分数
	f.SetCellValue("Sheet1", "H11", data.TotalGrade)

	// G13. 绩效基数
	f.SetCellValue("Sheet1", "G13", data.Coefficient)

	// G14. 最终得分系数

	// G17. 当月绩效奖金

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板-研发-%v.xlsx", year, int(month), year, int(month), data.AccountName)
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
func GetStoryBase2(estimate, consumed float64) float64 {
	return (estimate / STORY_BASE_TIME2) * STORY_BASE_SCORE2
}

func GetKpiGradeStandard2(totalGrade float64) float64 {
	if totalGrade >= 90 {
		return TOP_COEFFICIENT2
	} else if totalGrade < 90 && totalGrade >= 70 {
		return SECOND_COEFFICIENT2
	} else if totalGrade < 70 && totalGrade >= 60 {
		return THIRD_COEFFICIENT2
	}
	return 0
}

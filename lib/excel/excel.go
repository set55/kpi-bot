package excel

import (
	"fmt"
	"kpi-bot/common"
	"kpi-bot/lib/pm"
	"kpi-bot/lib/rd"
	"kpi-bot/lib/test"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)




func MakeRdExcel(path string, data rd.RdKpiGrade) error {
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("云平台组 研发工程师岗%v年%v月绩效考核表", year, int(month)))
	
	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：云平台组")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	projectDetail := "项目名称/冲刺名称/预期结束时间/实际结束时间/差值\n\n"
	for _, v := range data.ProjectProgressList {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.ProjectName, v.SprintName, v.End, v.RealEnd, v.DiffDays)
	}
	f.SetCellValue("Sheet1", "G4", projectDetail)

	// H4. 项目进度达成率 最终得分
	f.SetCellValue("Sheet1", "H4", data.AvgProgressStandardGrade)

	// G5. 需求达成率 完成情况
	storyDetail := "需求id/需求标题/预估工时/需求分数\n\n"
	for _, v := range data.StoryList {
		storyDetail += fmt.Sprintf("%v/%v/%v/%v\n\n", v.Id, v.Title, v.Estimate, v.Score)
	}
	f.SetCellValue("Sheet1", "G5", storyDetail)

	// H5. 需求达成率 最终得分
	f.SetCellValue("Sheet1", "H5", data.TotalStoryScore)

	// G6. bug遗留率 完成情况
	bugDetail := "测试报告/bug id/bug标题/bug解决方案/bug状态\n\n"
	for _, v := range data.BugInfoList {
		bugDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.TestReport, v.BugId, v.BugTitle, v.BugResolution, v.BugStatus)
	}
	f.SetCellValue("Sheet1", "G6", bugDetail)

	// H6. bug遗留率 最终得分
	f.SetCellValue("Sheet1", "H6", data.BugCarryStandardGrade)

	// G7. 工时预估达成比 完成情况
	f.SetCellValue("Sheet1", "G7", fmt.Sprintf("工时预估达成比：%v", data.TimeEstimateRate))

	// H7. 工时预估达成比 最终得分
	f.SetCellValue("Sheet1", "H7", data.TimeEstimateStandardGrade)

	// G8. 版本提测次数 完成情况
	pubDetail := "项目类型/项目名称/发版次数/最后一次提测时间\n\n"
	for _, v := range data.PubTimeList {
		pubDetail += fmt.Sprintf("%v/%v/%v/%v\n\n", v.ProjectType, v.ProjectName, v.PubTimes, v.LastPubTime)
	}
	f.SetCellValue("Sheet1", "G8", pubDetail)

	// H8. 版本提测次数 最终得分
	f.SetCellValue("Sheet1", "H8", data.AvgPubTimesStandardGrade)

	// H11. 总分数
	f.SetCellValue("Sheet1", "H11", data.TotalGrade)

	// G13. 绩效基数
	f.SetCellValue("Sheet1", "G13", data.TotalGradeStandard)

	// G14. 最终得分系数
	f.SetCellValue("Sheet1", "G14", data.TotalGradeStandard)

	// G17. 当月绩效奖金
	f.SetCellValue("Sheet1", "G17", data.TotalGradeStandard * common.GetRewardByAccount(data.Account))


	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板-研发-%v.xlsx", year, int(month), year, int(month), common.AccountToName(data.Account))
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


func MakeRdWithoutTestreportExcel(path string, data rd.RdWithoutTestReportKpiGrade) error {
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("云平台组 研发工程师岗(无测试报告)%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：云平台组")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	projectDetail := "项目名称/冲刺名称/预期结束时间/实际结束时间/差值\n\n"
	for _, v := range data.ProjectProgressList {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.ProjectName, v.SprintName, v.End, v.RealEnd, v.DiffDays)
	}
	f.SetCellValue("Sheet1", "G4", projectDetail)

	// H4. 项目进度达成率 最终得分
	f.SetCellValue("Sheet1", "H4", data.AvgProgressStandardGrade)

	// G5. 需求达成率 完成情况
	storyDetail := "需求id/需求标题/预估工时/需求分数\n\n"
	for _, v := range data.StoryList {
		storyDetail += fmt.Sprintf("%v/%v/%v/%v\n\n", v.Id, v.Title, v.Estimate, v.Score)
	}
	f.SetCellValue("Sheet1", "G5", storyDetail)

	// H5. 需求达成率 最终得分
	f.SetCellValue("Sheet1", "H5", data.TotalStoryScore)

	// G6. bug遗留率 完成情况
	bugDetail := "项目/bug id/bug标题/bug解决方案/bug状态\n\n"
	for _, v := range data.BugInfoList {
		bugDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.ProjectName, v.BugId, v.BugTitle, v.BugResolution, v.BugStatus)
	}
	f.SetCellValue("Sheet1", "G6", bugDetail)

	// H6. bug遗留率 最终得分
	f.SetCellValue("Sheet1", "H6", data.BugCarryStandardGrade)

	// G7. 工时预估达成比 完成情况
	f.SetCellValue("Sheet1", "G7", fmt.Sprintf("工时预估达成比：%v", data.TimeEstimateRate))

	// H7. 工时预估达成比 最终得分
	f.SetCellValue("Sheet1", "H7", data.TimeEstimateStandardGrade)

	// H10. 总分数
	f.SetCellValue("Sheet1", "H10", data.TotalGrade)

	// G12. 绩效基数
	f.SetCellValue("Sheet1", "G12", data.TotalGradeStandard)

	// G13. 最终得分系数
	f.SetCellValue("Sheet1", "G13", data.TotalGradeStandard)

	// G16. 当月绩效奖金
	f.SetCellValue("Sheet1", "G16", data.TotalGradeStandard * common.GetRewardByAccount(data.Account))

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板(无测试报告)-研发-%v.xlsx", year, int(month), year, int(month), common.AccountToName(data.Account))
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


func MakePmExcel(path string, data pm.PmKpiGrade) error {
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("云平台组 项目经理岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：云平台组")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	f.SetCellValue("Sheet1", "G4", fmt.Sprintf("平均差值天数: %v", data.ProgressAvgDiffDays))

	// H4. 项目进度达成率 最终得分
	f.SetCellValue("Sheet1", "H4", data.ProgressStandardGrade)

	// G5. 项目成果完整率 完成情况
	projectCompleteDetail := "项目名称/需求完整率\n\n"
	for _, v := range data.ProjectCompleteList {
		projectCompleteDetail += fmt.Sprintf("%v/%v\n\n", v.ProjectName, v.CompleteRate)
	}
	f.SetCellValue("Sheet1", "G5", projectCompleteDetail)

	// H5. 项目成果完整 最终得分
	f.SetCellValue("Sheet1", "H5", data.CompleteRateStandardGrade)

	// G6. 项目规划需求数 完成情况
	storyNumDetail := "Projected/Developed/Closed\n\n"
	storyNumDetail += fmt.Sprintf("%v/%v/%v\n\n", data.ProjectedStoryNum, data.DevelopedStoryNum, data.ClosedStoryNum)
	f.SetCellValue("Sheet1", "G6", storyNumDetail)

	// H6. bug遗留率 最终得分
	f.SetCellValue("Sheet1", "H6", data.StoryNumGrade)

	// G7. 预估承诺完成率 完成情况
	f.SetCellValue("Sheet1", "G7", fmt.Sprintf("平均承诺天数差值：%v", data.PromiseDiffDays))

	// H7. 预估承诺完成率 最终得分
	f.SetCellValue("Sheet1", "H7", data.PromiseStandardGrade)

	// H10. 总分数
	f.SetCellValue("Sheet1", "H10", data.TotalGrade)

	// G12. 绩效基数
	f.SetCellValue("Sheet1", "G12", data.TotalGradeStandard)

	// G13. 最终得分系数
	f.SetCellValue("Sheet1", "G13", data.TotalGradeStandard)

	// G16. 当月绩效奖金
	f.SetCellValue("Sheet1", "G16", data.TotalGradeStandard * common.GetRewardByAccount(data.Account))


	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板-项目-%v.xlsx", year, int(month), year, int(month), common.AccountToName(data.Account))
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


func MakeTestExcel(path string, data test.TestKpiGrade) error {
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("云平台组 测试工程师岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：云平台组")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	f.SetCellValue("Sheet1", "G4", fmt.Sprintf("平均差值天数: %v", data.TestProgressAvgDiffDays))

	// H4. 项目进度达成率 最终得分
	f.SetCellValue("Sheet1", "H4", data.TestProgressAvgDiffDaysStandardGrade)

	// G5. 项目有效bug率  完成情况
	f.SetCellValue("Sheet1", "G5", fmt.Sprintf("有效bug率: %v", data.ValidateBugRate))

	// H5. 项目有效bug率  最终得分
	f.SetCellValue("Sheet1", "H5", data.ValidateBugRateStandardGrade)

	// G6. bug转需求数 完成情况
	f.SetCellValue("Sheet1", "G6", fmt.Sprintf("bug转需求数: %v", data.BugToStoryNum))

	// H6. bug转需求数 最终得分
	f.SetCellValue("Sheet1", "H6", data.BugToStoryGrade)

	// G7. 用例发现Bug率 完成情况
	f.SetCellValue("Sheet1", "G7", fmt.Sprintf("用例发现Bug率：%v", data.CaseBugRate))

	// H7. 用例发现Bug率 最终得分
	f.SetCellValue("Sheet1", "H7", data.CaseBugRateStandardGrade)

	// H10. 总分数
	f.SetCellValue("Sheet1", "H10", data.TotalGrade)

	// G12. 绩效基数
	f.SetCellValue("Sheet1", "G12", data.TotalGradeStandard)

	// G13. 最终得分系数
	f.SetCellValue("Sheet1", "G13", data.TotalGradeStandard)

	// G16. 当月绩效奖金
	f.SetCellValue("Sheet1", "G16", data.TotalGradeStandard * common.GetRewardByAccount(data.Account))


	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板-测试-%v.xlsx", year, int(month), year, int(month), common.AccountToName(data.Account))
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
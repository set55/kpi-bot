package excel

import (
	"fmt"
	"kpi-bot/common"
	"kpi-bot/lib/deveops"
	"kpi-bot/lib/pm"
	"kpi-bot/lib/rd"
	"kpi-bot/lib/test"
	"math/rand"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

// 绩效考核-研发
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 研发工程师岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度延时率 完成情况
	projectDetail := fmt.Sprintf("平均项目延时率：%v\n\n", data.AvgDiffRate)
	projectDetail += "项目id/项目名称/计划开始时间/计划结束时间/实际结束时间\n\n"

	for _, v := range data.ProjectProgressList {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.ProjectId, v.ProjectName, v.Begin, v.End, v.RealEnd)
	}
	f.SetCellValue("Sheet1", "G4", projectDetail)

	// H4. 项目进度延时率 最终得分
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
	bugDetail := fmt.Sprintf("测试单数量: %v, bug遗留率: %v\n\n", data.TestTaskCount, data.BugCarryOverRate)
	bugDetail += "项目名称/bug id/bug标题/bug解决方案/bug状态\n\n"
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

	// G8. 版本提测次数 完成情况
	pubDetail := fmt.Sprintf("平均提测次数：%v\n\n", data.AvgPubTimes)
	pubDetail += "项目类型/项目名称/发版次数/最后一次提测时间\n\n"
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
	f.SetCellValue("Sheet1", "G17", data.TotalGradeStandard*common.GetRewardByAccount(data.Account))

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

// 绩效考核-研发(无测试报告)
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 研发工程师岗(无测试报告)%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	projectDetail := fmt.Sprintf("平均项目延时率：%v\n\n", data.AvgDiffRate)
	projectDetail += "项目id/项目名称/计划开始时间/计划结束时间/实际结束时间\n\n"

	for _, v := range data.ProjectProgressList {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.ProjectId, v.ProjectName, v.Begin, v.End, v.RealEnd)
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
	bugDetail := fmt.Sprintf("bug遗留率: %v\n\n", data.BugCarryOverRate)
	bugDetail += "项目/bug id/bug标题/bug解决方案/bug状态\n\n"
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
	f.SetCellValue("Sheet1", "G16", data.TotalGradeStandard*common.GetRewardByAccount(data.Account))

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

// 绩效考核-项目
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 项目经理岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	projectDetail := fmt.Sprintf("延時率：%v\n\n", data.DiffRate)
	projectDetail += "项目id/项目标题/项目类型/项目开始时间/项目预估结束时间/项目实际结束时间/测试开始时间/测试预估结束时间/测试实际结束时间\n\n"
	for _, v := range data.ProjectProgressList {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v/%v/%v/%v/%v\n\n", v.ProjectId, v.ProjectName, v.ProjectType, v.ProjectBegin, v.ProjectEnd, v.ProjectRealEnd, v.TestStart, v.TestEnd, v.TestRealEnd)
	}
	f.SetCellValue("Sheet1", "G4", projectDetail)

	// H4. 项目进度达成率 最终得分
	f.SetCellValue("Sheet1", "H4", data.ProgressStandardGrade)

	// G5. 项目成果完整率 完成情况
	projectCompleteDetail := fmt.Sprintf("平均完成率: %v\n\n", data.CompleteRate)
	projectCompleteDetail += "项目名称/需求完整率\n\n"
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

	// H6. 项目规划需求数 最终得分
	f.SetCellValue("Sheet1", "H6", data.StoryNumGrade)

	// // G7. 预估承诺完成率 完成情况
	// f.SetCellValue("Sheet1", "G7", fmt.Sprintf("平均承诺天数差值：%v", data.PromiseDiffDays))

	// // H7. 预估承诺完成率 最终得分
	// f.SetCellValue("Sheet1", "H7", data.PromiseStandardGrade)

	// G7. 预估工时准确率 完成情况
	timeEstimateDetail := fmt.Sprintf("平均预估工时准确率: %v\n\n", data.TimeEstimateRate)
	timeEstimateDetail += "需求id/需求名称/预估工时/实际工时/准确率\n\n"
	for _, v := range data.TimeEstimateList {
		timeEstimateDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.StoryId, v.Title, v.Estimate, v.StoryConsumed, v.EstimateRate)
	}
	f.SetCellValue("Sheet1", "G7", timeEstimateDetail)

	// H7. 预估工时准确率 最终得分
	f.SetCellValue("Sheet1", "H7", data.TimeEstimateGrade)

	// H10. 总分数
	f.SetCellValue("Sheet1", "H10", data.TotalGrade)

	// G12. 绩效基数
	f.SetCellValue("Sheet1", "G12", data.TotalGradeStandard)

	// G13. 最终得分系数
	f.SetCellValue("Sheet1", "G13", data.TotalGradeStandard)

	// G16. 当月绩效奖金
	f.SetCellValue("Sheet1", "G16", data.TotalGradeStandard*common.GetRewardByAccount(data.Account))

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

// 绩效考核-项目(无测试报告)
func MakePmExcelWithoutTestReport(path string, data pm.PmKpiGradeWithoutTestReport) error {
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 项目经理岗(无测试报告)%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	projectDetail := fmt.Sprintf("延時率：%v\n\n", data.DiffRate)
	projectDetail += "项目id/项目标题/项目类型/项目开始时间/项目预估结束时间/项目实际结束时间\n\n"
	for _, v := range data.ProjectProgressList {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v/%v\n\n", v.ProjectId, v.ProjectName, v.ProjectType, v.ProjectBegin, v.ProjectEnd, v.ProjectRealEnd)
	}
	f.SetCellValue("Sheet1", "G4", projectDetail)

	// H4. 项目进度达成率 最终得分
	f.SetCellValue("Sheet1", "H4", data.ProgressStandardGrade)

	// G5. 项目成果完整率 完成情况
	projectCompleteDetail := fmt.Sprintf("平均完成率: %v\n\n", data.CompleteRate)
	projectCompleteDetail += "项目名称/需求完整率\n\n"
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

	// H6. 项目规划需求数 最终得分
	f.SetCellValue("Sheet1", "H6", data.StoryNumGrade)

	// // G7. 预估承诺完成率 完成情况
	// f.SetCellValue("Sheet1", "G7", fmt.Sprintf("平均承诺天数差值：%v", data.PromiseDiffDays))

	// // H7. 预估承诺完成率 最终得分
	// f.SetCellValue("Sheet1", "H7", data.PromiseStandardGrade)

	// G7. 预估工时准确率 完成情况
	timeEstimateDetail := fmt.Sprintf("平均预估工时准确率: %v\n\n", data.TimeEstimateRate)
	timeEstimateDetail += "需求id/需求名称/预估工时/实际工时/准确率\n\n"
	for _, v := range data.TimeEstimateList {
		timeEstimateDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.StoryId, v.Title, v.Estimate, v.StoryConsumed, v.EstimateRate)
	}
	f.SetCellValue("Sheet1", "G7", timeEstimateDetail)

	// H7. 预估工时准确率 最终得分
	f.SetCellValue("Sheet1", "H7", data.TimeEstimateGrade)

	// H10. 总分数
	f.SetCellValue("Sheet1", "H10", data.TotalGrade)

	// G12. 绩效基数
	f.SetCellValue("Sheet1", "G12", data.TotalGradeStandard)

	// G13. 最终得分系数
	f.SetCellValue("Sheet1", "G13", data.TotalGradeStandard)

	// G16. 当月绩效奖金
	f.SetCellValue("Sheet1", "G16", data.TotalGradeStandard*common.GetRewardByAccount(data.Account))

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板(无测试报告)-项目-%v.xlsx", year, int(month), year, int(month), common.AccountToName(data.Account))
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

// 绩效考核-测试
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 测试工程师岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	projectDetail := fmt.Sprintf("延時率: %v\n\n", data.DiffRate)
	projectDetail += "测试任务/测试报告/测试任务开始时间/测试任务计划结束时间/测试报告生成时间\n\n"
	for _, v := range data.TestProgressInfos {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.TestTaskName, v.TestReportTitle, v.TestTaskBegin, v.TestTaskEnd, v.TestReportEnd)
	}
	f.SetCellValue("Sheet1", "G4", projectDetail)

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
	f.SetCellValue("Sheet1", "G16", data.TotalGradeStandard*common.GetRewardByAccount(data.Account))

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

// kpi统计
func MakeKpiStatisticsExcel(path string, startTime string, rdGrades map[string]rd.RdKpiGrade, rdWithoutGrades map[string]rd.RdWithoutTestReportKpiGrade,
	pmGrades map[string]pm.PmKpiGrade, pmWithout map[string]pm.PmKpiGradeWithoutTestReport, testsGrade map[string]test.TestKpiGrade) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("open file fail: %v", err)
	}
	defer f.Close()

	// Parse the date string into a time.Time object
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, startTime)
	if err != nil {
		return fmt.Errorf("parse time err: %v", err)
	}

	// Extract the year and month
	year := t.Year()
	month := t.Month()

	statisticMap := map[string]int64{
		"lower60":   0,
		"60-69":     0,
		"70-79":     0,
		"80-89":     0,
		"90-100":    0,
		"101-110":   0,
		"higher110": 0,
	}

	for _, v := range rdGrades {
		grade := v.TotalGrade
		if grade < 60 {
			statisticMap["lower60"]++
		} else if grade >= 60 && grade < 70 {
			statisticMap["60-69"]++
		} else if grade >= 70 && grade < 80 {
			statisticMap["70-79"]++
		} else if grade >= 80 && grade < 90 {
			statisticMap["80-89"]++
		} else if grade >= 90 && grade <= 100 {
			statisticMap["90-100"]++
		} else if grade > 100 && grade <= 110 {
			statisticMap["101-110"]++
		} else if grade > 110 {
			statisticMap["higher110"]++
		}
	}

	for _, v := range rdWithoutGrades {
		grade := v.TotalGrade
		if grade < 60 {
			statisticMap["lower60"]++
		} else if grade >= 60 && grade < 70 {
			statisticMap["60-69"]++
		} else if grade >= 70 && grade < 80 {
			statisticMap["70-79"]++
		} else if grade >= 80 && grade < 90 {
			statisticMap["80-89"]++
		} else if grade >= 90 && grade <= 100 {
			statisticMap["90-100"]++
		} else if grade > 100 && grade <= 110 {
			statisticMap["101-110"]++
		} else if grade > 110 {
			statisticMap["higher110"]++
		}
	}

	for _, v := range pmGrades {
		grade := v.TotalGrade
		if grade < 60 {
			statisticMap["lower60"]++
		} else if grade >= 60 && grade < 70 {
			statisticMap["60-69"]++
		} else if grade >= 70 && grade < 80 {
			statisticMap["70-79"]++
		} else if grade >= 80 && grade < 90 {
			statisticMap["80-89"]++
		} else if grade >= 90 && grade <= 100 {
			statisticMap["90-100"]++
		} else if grade > 100 && grade <= 110 {
			statisticMap["101-110"]++
		} else if grade > 110 {
			statisticMap["higher110"]++
		}
	}

	for _, v := range pmWithout {
		grade := v.TotalGrade
		if grade < 60 {
			statisticMap["lower60"]++
		} else if grade >= 60 && grade < 70 {
			statisticMap["60-69"]++
		} else if grade >= 70 && grade < 80 {
			statisticMap["70-79"]++
		} else if grade >= 80 && grade < 90 {
			statisticMap["80-89"]++
		} else if grade >= 90 && grade < 100 {
			statisticMap["90-100"]++
		} else if grade >= 100 && grade <= 110 {
			statisticMap["101-110"]++
		} else if grade > 110 {
			statisticMap["higher110"]++
		}
	}

	for _, v := range testsGrade {
		grade := v.TotalGrade
		if grade < 60 {
			statisticMap["lower60"]++
		} else if grade >= 60 && grade < 70 {
			statisticMap["60-69"]++
		} else if grade >= 70 && grade < 80 {
			statisticMap["70-79"]++
		} else if grade >= 80 && grade < 90 {
			statisticMap["80-89"]++
		} else if grade >= 90 && grade <= 100 {
			statisticMap["90-100"]++
		} else if grade > 100 && grade <= 110 {
			statisticMap["101-110"]++
		} else if grade > 110 {
			statisticMap["higher110"]++
		}
	}

	// A2 小于60
	f.SetCellValue("Sheet1", "A2", statisticMap["lower60"])

	// B2 60-69
	f.SetCellValue("Sheet1", "B2", statisticMap["60-69"])

	// C2 70-79
	f.SetCellValue("Sheet1", "C2", statisticMap["70-79"])

	// D2 80-89
	f.SetCellValue("Sheet1", "D2", statisticMap["80-89"])

	// E2 90-100
	f.SetCellValue("Sheet1", "E2", statisticMap["90-100"])

	// F2 101-110
	f.SetCellValue("Sheet1", "F2", statisticMap["101-110"])

	// G2 大于110
	f.SetCellValue("Sheet1", "G2", statisticMap["higher110"])

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-kpi统计.xlsx", year, int(month), year, int(month))
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

// 绩效考核-运维
func MakeDeveopsExcel(path string, data deveops.DeveopsKpiGrade) error {
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
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("软件服务中心 运维工程师岗%v年%v月绩效考核表", year, int(month)))

	// A2. 被考评人员部门：XXXX
	f.SetCellValue("Sheet1", "A2", "被考评人员部门：软件服务中心")

	// E2. 被考评人员：XXXX
	f.SetCellValue("Sheet1", "E2", fmt.Sprintf("被考评人员：%v", common.AccountToName(data.Account)))

	// F2. 考评人：xxxx
	f.SetCellValue("Sheet1", "F2", fmt.Sprintf("考评人：%v", "Set"))

	// G4. 项目进度达成率 完成情况
	projectDetail := fmt.Sprintf("平均项目延时率：%v\n\n", data.AvgDiffRate)
	projectDetail += "项目id/项目名称/计划开始时间/计划结束时间/实际结束时间\n\n"

	for _, v := range data.ProjectProgressList {
		projectDetail += fmt.Sprintf("%v/%v/%v/%v/%v\n\n", v.ProjectId, v.ProjectName, v.Begin, v.End, v.RealEnd)
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

	// G7. 工时预估达成比 完成情况
	f.SetCellValue("Sheet1", "G7", fmt.Sprintf("工时预估达成比：%v", data.TimeEstimateRate))

	// H7. 工时预估达成比 最终得分
	f.SetCellValue("Sheet1", "H7", data.TimeEstimateStandardGrade)

	// H10. 总分数
	f.SetCellValue("Sheet1", "H11", data.TotalGrade)

	// G12. 绩效基数
	f.SetCellValue("Sheet1", "G13", data.TotalGradeStandard)

	// G13. 最终得分系数
	f.SetCellValue("Sheet1", "G14", data.TotalGradeStandard)

	// G16. 当月绩效奖金
	f.SetCellValue("Sheet1", "G17", data.TotalGradeStandard*common.GetRewardByAccount(data.Account))

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/%v-%v", year, int(month))
	filePath := fmt.Sprintf("./export/%v-%v/%v-%v-绩效考核模板-运维-%v.xlsx", year, int(month), year, int(month), common.AccountToName(data.Account))
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

// 
func MakeWhatEverExcel(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("open file fail: %v", err)
	}
	defer f.Close()

	columns := []string{"J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"}

	// for _, col := range columns {
	// 	for i := 5; i <= 41; i++ {
	// 		randomValue := rand.Intn(5) + 6 // Generate a random value between 70 and 90
	// 		f.SetCellValue("管理岗汇总", fmt.Sprintf("%v%v", col, i), fmt.Sprintf("%v", randomValue))
	// 	}
	// }

	for i := 5; i <= 41; i++ {
		tmpTotalGrade := 0
		for _, col := range columns {
			if col != "T" && col != "U" && col != "V" {
				randomValue := rand.Intn(4) + 6 // Generate a random value between 5 and 10
				tmpTotalGrade += randomValue
				f.SetCellValue("管理岗汇总", fmt.Sprintf("%v%v", col, i), fmt.Sprintf("%v", randomValue))
			}

			if col == "T" {
				f.SetCellValue("管理岗汇总", fmt.Sprintf("%v%v", col, i), fmt.Sprintf("%v", tmpTotalGrade))
			}

			if col == "U" {
				tmpAbility := "低"
				if tmpTotalGrade >= 70 && tmpTotalGrade <= 89 {
					tmpAbility = "中"
				}

				if tmpTotalGrade >= 90 && tmpTotalGrade <= 100 {
					tmpAbility = "高"
				}
				f.SetCellValue("管理岗汇总", fmt.Sprintf("%v%v", col, i), tmpAbility)
			}

			if col == "V" {
				abilities := []string{"中", "高"}
                randomAbility := abilities[rand.Intn(len(abilities))]
				f.SetCellValue("管理岗汇总", fmt.Sprintf("%v%v", col, i), randomAbility)
			}

		}

	}

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/管理岗人才胜任盘点2025")
	filePath := fmt.Sprintf("./export/管理岗人才胜任盘点2025/管理岗人才胜任盘点2025.xlsx")
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


func MakeWhatEverExcelNormal(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("open file fail: %v", err)
	}
	defer f.Close()

	columns := []string{"J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"}

	// for _, col := range columns {
	// 	for i := 5; i <= 41; i++ {
	// 		randomValue := rand.Intn(5) + 6 // Generate a random value between 70 and 90
	// 		f.SetCellValue("管理岗汇总", fmt.Sprintf("%v%v", col, i), fmt.Sprintf("%v", randomValue))
	// 	}
	// }
	countPotentialLow := 10
	for i := 5; i <= 254; i++ {
		tmpTotalGrade := 0
		for _, col := range columns {
			if col != "T" && col != "U" && col != "V" {
				base := 6
				if i >= 141 && i <= 145 {
					base = 7
				}

				randomValue := rand.Intn(4) + base // Generate a random value between 5 and 10
				tmpTotalGrade += randomValue
				f.SetCellValue("非管理岗汇总", fmt.Sprintf("%v%v", col, i), fmt.Sprintf("%v", randomValue))
			}

			if col == "T" {
				f.SetCellValue("非管理岗汇总", fmt.Sprintf("%v%v", col, i), fmt.Sprintf("%v", tmpTotalGrade))
			}

			if col == "U" {
				tmpAbility := "低"
				if tmpTotalGrade >= 70 && tmpTotalGrade <= 89 {
					tmpAbility = "中"
				}

				if tmpTotalGrade >= 90 && tmpTotalGrade <= 100 {
					tmpAbility = "高"
				}
				f.SetCellValue("非管理岗汇总", fmt.Sprintf("%v%v", col, i), tmpAbility)
			}

			if col == "V" {
				abilities := []string{"中", "高"}
				if !(i >= 141 && i <= 145) && countPotentialLow <= 10 {
					abilities = []string{"低", "中", "高"}
				}
                randomAbility := abilities[rand.Intn(len(abilities))]
				if randomAbility == "低" {
					countPotentialLow++
				}
				f.SetCellValue("非管理岗汇总", fmt.Sprintf("%v%v", col, i), randomAbility)
			}

		}

	}

	// 建立资料夹
	folderPath := fmt.Sprintf("./export/一般员工人才胜任力盘点2025")
	filePath := fmt.Sprintf("./export/一般员工人才胜任力盘点2025/一般员工人才胜任力盘点2025.xlsx")
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

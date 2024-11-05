package common

import (
	"fmt"
	"strings"
	"time"
)


func AccountArrayToString(accounts []string) string {
	arr := []string{}
	for _, v := range accounts {
        arr = append(arr, fmt.Sprintf("\"%s\"", v))
    }

    // Join the array elements into a single string
    return strings.Join(arr, ",")
}


func AccountToName(account string) string {
	switch account {
	case "set.su":
		return "Set"
	case "paul.gao":
		return "高长荣"
	case "justin.lee":
		return "李玠廷"
	case "shawn.wang":
		return "汪晓航"
	case "samy.gou":
		return "缑富永"
	case "champion.fu":
		return "付庆平"
	case "alan.tin":
		return "田佳发"
	case "shiwen.tin":
		return "方世文"
	case "guoqiao.chen":
		return "陈国桥"
	case "xiechen":
		return "谢晨"
	case "zouyanling":
		return "邹燕玲"
	case "ruanbanyong":
		return "阮班勇"
	case "zhouyao":
		return "周尧"
	case "liuxiaoyan":
		return "刘晓彦"
	case "linyanhai":
		return "林焰海"
	case "jihuaqing":
		return "吉桦庆"
	case "liuhongtao":
		return "刘洪涛"
	case "wangtuhe":
		return "王土何"
	case "deakin.han":
		return "韩象金"
	case "qixiaofeng":
		return "祁晓锋"
	case "wangshaoyu":
		return "王少宇"
	case "simon.chen":
		return "陈熙存"
	default:
		return account
	}
}

func ProjectTypeTransform(projectType string) string {
	switch projectType {
	case "project":
		return "项目"
	case "sprint":
		return "冲刺"
	default:
		return "unknown"
	}
}

func GetRewardByAccount(account string) float64 {
	switch account {
	case "set.su":
		return 0
	case "paul.gao":
		return 0
	case "justin.lee":
		return 0
	case "shawn.wang":
		return 2000
	case "samy.gou":
		return 0
	case "champion.fu":
		return 2000
	case "alan.tin":
		return 1000
	case "shiwen.tin":
		return 2000
	case "guoqiao.chen":
		return 3000
	case "xiechen":
		return 1000
	case "zouyanling":
		return 0
	case "ruanbanyong":
		return 1000
	case "zhouyao":
		return 1000
	case "liuxiaoyan":
		return 2500
	case "linyanhai":
		return 0
	case "jihuaqing":
		return 2000
	case "liuhongtao":
		return 2000
	case "wangtuhe":
		return 2000
	case "deakin.han":
		return 3500
	case "qixiaofeng":
		return 3500
	case "wangshaoyu":
		return 3000
	case "simon.chen":
		return 2000
	default:
		return 0
	}
}

func GetProjectProgressExpectRate(planDiff, realDiff float64) (float64) {
	fmt.Printf("planDiff: %v, realDiff: %v\n", planDiff, realDiff)
	if planDiff == 0 {
		return 2 // 若计划天数为0, rate视为大于1.2 给最低0分
	}

	planSubstractDays := int(planDiff / 7) // 每7天 -1天
	realSubstractDays := int(realDiff / 7) // 每7天 -1天

	finalPlanDiff := planDiff - float64(planSubstractDays)
	finalRealDiff := realDiff - float64(realSubstractDays)
	fmt.Printf("finalPlanDiff: %v, finalRealDiff: %v\n", finalPlanDiff, finalRealDiff)

	return (finalRealDiff / finalPlanDiff) - 1
}


func CountWeekends(start, end string) (int, int) {
    saturdays := 0
    sundays := 0

	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, start)
	endDate, _ := time.Parse(layout, end)



    // Loop through each day in the date range
    for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
        switch d.Weekday() {
        case time.Saturday:
            saturdays++
        case time.Sunday:
            sundays++
        }
    }
    return saturdays, sundays
}

func GetBugStandard(bugRate float64) float64 {
	if bugRate == 0 {
		return 1
	} else if bugRate <= 0.1 {
		return 0.9
	} else if bugRate <= 0.2 {
		return 0.8
	} else if bugRate <= 0.3 {
		return 0.7
	} else if bugRate <= 0.4 {
		return 0.6
	} else if bugRate <= 0.5 {
		return 0.5
	} else {
		return 0
	}
}
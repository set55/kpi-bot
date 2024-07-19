package common

import (
	"fmt"
	"strings"
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
		return 1000
	case "wangtuhe":
		return 2000
	case "deakin.han":
		return 0
	case "qixiaofeng":
		return 0
	case "wangshaoyu":
		return 0
	default:
		return 0
	}
}
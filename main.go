package main

import (
	"database/sql"
	// "kpi-bot/lib/excel"
	// "encoding/json"
	"fmt"
	// "kpi-bot/lib/pm"
	// "kpi-bot/lib/rd"
	// "kpi-bot/lib/test"
	"log"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/xuri/excelize/v2"
)

func main() {
	// Define the data source name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "developer", "Mh-mJ?sp.G\"43*_HrCXRP9+^QS%3Et2yZE", "192.168.2.8", "32606", "zentao")

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	fmt.Println("Connected to MariaDB successfully!")

	// // all team members
	// rds := []string{"set.su", "paul.gao", "justin.lee", "samy.gou", "champion.fu", "alan.tin", "jihuaqing", "liuhongtao", "deakin.han"}
	// rdsWithoutTest := []string{"shiwen.tin", "xiechen", "zouyanling", "ruanbanyong", "zhouyao", "liuxiaoyan", "wangtuhe"}
	// tests := []string{"linyanhai", "wangshaoyu"}
	// pms := []string{"shawn.wang", "qixiaofeng"}

	// // new kpi instance
	// rdKpi := rd.NewRdKpi(db, rds)
	// rdKpiWithoutTest := rd.NewRdKpiWithoutTestReport(db, rdsWithoutTest)
	// testKpi := test.NewTestKpi(db, tests)
	// pmKpi := pm.NewPmKpi(db, pms)
	
	
	// rdKpiGrades := rdKpi.GetRdKpiGrade()
	// rdKpiWithoutTestGrades := rdKpiWithoutTest.GetRdKpiWithoutTestReportGrade()
	// testKpiGrades := testKpi.GetTestKpiGrade()
	// pmKpiGrades := pmKpi.GetPmKpiGrade()

	// rdJson, _ := json.Marshal(rdKpiGrades)
	// rdWithoutTestJson, _ := json.Marshal(rdKpiWithoutTestGrades)
	// testJson, _ := json.Marshal(testKpiGrades)
	// pmJson, _ := json.Marshal(pmKpiGrades)

	// fmt.Printf("rd: %v\n\n", string(rdJson))
	// fmt.Printf("rd no test: %v\n\n", string(rdWithoutTestJson))
	// fmt.Printf("test: %v\n\n", string(testJson))
	// fmt.Printf("pm: %v\n\n", string(pmJson))

	// Open the Excel file







    // f, err := excelize.OpenFile("./excel/绩效考核模板-研发.xlsx")
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // defer f.Close()
	// err = f.SetCellValue("Sheet1", "A1", "tttttttt")
	// if err != nil {
	// 	fmt.Println(err)
	// }


	// Iterate over all the sheets
    // for _, sheetName := range f.GetSheetList() {
    //     // Get the dimensions of the sheet
    //     maxRow, maxCol := getSheetDimensions(f, sheetName)
    //     fmt.Printf("SheetName: %s, maxRow: %v, maxCol: %v\n", sheetName, maxRow, maxCol)
    //     for rowIndex := 1; rowIndex <= maxRow; rowIndex++ {
    //         for colIndex := 1; colIndex <= maxCol; colIndex++ {
    //             // Get the cell name
    //             cellName, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
    //             // Get the cell value
    //             cellValue, err := f.GetCellValue(sheetName, cellName)
    //             if err != nil {
    //                 fmt.Println(err)
    //                 return
    //             }
                
    //             // Print the content of each cell
    //             fmt.Printf("Sheet: %s, RowIndex: %v, ColIndex: %v, Cell: %s, Content: %s\n", sheetName, rowIndex, colIndex, cellName, cellValue)
                
    //             if cellValue == "" {
    //                 // Insert the word "Filled" into empty cells
    //                 if err := f.SetCellValue(sheetName, cellName, "Filled"); err != nil {
    //                     fmt.Println(err)
    //                     return
    //                 }
    //             }
    //         }
    //     }
    // }

	

	// // Save the modified file
    // if err = f.SaveAs("./export/modified_绩效考核模板-研发.xlsx"); err != nil {
    //     fmt.Println(err)
    // }


	// err = excel.MakeRdExcel("./excel/绩效考核模板-研发.xlsx", "2024-07-01 00:00:00", "2024-07-31 23:59:59")
	// if err != nil {
	// 	fmt.Println(err)
	// }

}

// func getSheetDimensions(f *excelize.File, sheetName string) (int, int) {
//     maxRow, maxCol := 0, 0
//     rows, err := f.GetRows(sheetName)
//     if err != nil {
//         return 0, 0
//     }
//     for rowIndex, row := range rows {
//         if len(row) > maxCol {
//             maxCol = len(row)
//         }
//         if rowIndex+1 > maxRow {
//             maxRow = rowIndex + 1
//         }
//     }
//     return maxRow, maxCol
// }

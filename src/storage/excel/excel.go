package Excel

import (
	"github.com/tealeg/xlsx"
	"log"
)

// arg sheetNum - count starts from 0
func ReadSheet(fileName string, sheetNum int) [][]string {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
//	for _, sheet := range xlFile.Sheets {
		sheetsAmount := len(xlFile.Sheets)

		if sheetNum >= sheetsAmount {
			return [][]string{}
		}

		sheet := xlFile.Sheets[sheetNum]
		bd := [][]string{}
		for _, row := range sheet.Rows {
			part := []string{}
			for _, cell := range row.Cells {
				part = append(part, cell.String())
			}
			bd = append(bd, part)
//			fmt.Println(part)
		}
//		fmt.Println(bd)
		return bd
//	}

}
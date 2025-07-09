package excel

import (
	"errors"
	"github.com/miiy/goc/utils/slice"
	"github.com/xuri/excelize/v2"
	"log"
	"reflect"
)

// WriteExcelData 写excel文件
// header 表头
// tag struct tag
// data slice 行
func WriteExcelData(header []string, tag string, data interface{}) ([]byte, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Close err:", err)
		}
	}()

	// 表头
	if err := file.SetSheetRow("Sheet1", "A1", &header); err != nil {
		log.Println("SetSheetRow err:", err)
	}

	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	if rt.Kind() != reflect.Slice {
		return nil, errors.New("导出错误")
	}
	// 行
	for rowNum := 0; rowNum < rv.Len(); rowNum++ {
		row := rv.Index(rowNum)
		rowSlice := slice.StructToSliceByTagValues(row.Interface(), tag, header)
		for colNum, v := range rowSlice {
			if err := setCellValue(file, "Sheet1", colNum+1, rowNum+2, v); err != nil {
				log.Println("setCellValue err:", err)
			}
		}
	}
	buf, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// SetCellValue 设置单元格值
func setCellValue(f *excelize.File, sheetName string, col, row int, value interface{}) error {
	cellName, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return err
	}
	err = f.SetCellValue(sheetName, cellName, value)
	if err != nil {
		return err
	}
	return nil
}

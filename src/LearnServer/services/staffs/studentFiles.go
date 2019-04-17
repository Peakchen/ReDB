package staffs

import (
	"log"
	"net/http"

	"LearnServer/utils"
	"github.com/labstack/echo"

	"github.com/tealeg/xlsx"
)

func deleteStudentTmpFileHandler(c echo.Context) error {
	uid := c.Param("uid")
	utils.Manager.Delete(uid)
	return c.JSON(http.StatusOK, "successfully deleted")
}

func createXlsErrorMessage(errMessage string) echo.Map {
	type errorType struct {
		Error string `json:"error"`
	}
	return echo.Map{
		"uid": "-1",
		"columns": utils.MakeColumns(
			utils.StrList{"错误", "error"},
		),
		"data": []errorType{
			errorType{
				Error: errMessage,
			},
		},
	}
}

func previewStudentFileHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	xls, err := xlsx.OpenReaderAt(src, file.Size)
	if err != nil {
		return err
	}
	columns := utils.MakeColumns(
		utils.StrList{"姓名", "name"},
		utils.StrList{"性别", "gender"},
	)

	if xls.Sheets[0].Rows[0].Cells[0].String() != "姓名" || xls.Sheets[0].Rows[0].Cells[1].String() != "性别" {
		return c.JSON(http.StatusOK, createXlsErrorMessage("表格列名称不正确！"))
	}

	data := []studentType{}
	for _, row := range xls.Sheets[0].Rows[1:] {
		s := studentType{}
		err = row.ReadStruct(&s)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, s)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"uid":     utils.Manager.Insert(data),
		"columns": columns,
		"data":    data,
	})
}

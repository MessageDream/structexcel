/**
 * @file: csv_render.go
 * @created: 2021-08-04 05:30:50
 * @author: jayden (jaydenzhao@outlook.com)
 * @description: code for project
 * @last modified: 2021-08-04 05:30:59
 * @modified by: jayden (jaydenzhao@outlook.com>)
 */
package excel

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CSVRender struct {
}

var _ IORender = &CSVRender{}

func (r *CSVRender) Render(to io.Writer, source interface{}) error {
	encodedWriter := transform.NewWriter(to, simplifiedchinese.GBK.NewEncoder())

	writer := csv.NewWriter(encodedWriter)
	defer writer.Flush()

	iter, errCh := iterator(source)
	batch := 0
	for row := range iter {
		if err := writer.Write(row); err != nil {
			return err
		}
		batch++
		if batch == 100 {
			writer.Flush()
			batch = 0
			if err := writer.Error(); err != nil {
				return err
			}
		}
	}

	if err := <-errCh; err != nil {
		return err
	}
	return writer.Error()
}

func NewHttpCSVRender(w http.ResponseWriter, filename string) Render {

	w.Header().Set("Content-type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fmtFileName(filename, ".csv")))
	w.Header().Set("Cache-Control", "must-revalidate,post-check=0,pre-check=0")
	w.Header().Set("Expires", "0")

	return &defaultRender{
		w,
		&CSVRender{},
	}
}

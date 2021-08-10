/**
 * @file: xls_render.go
 * @created: 2021-08-05 09:38:17
 * @author: jayden (jaydenzhao@outlook.com)
 * @description: code for project
 * @last modified: 2021-08-05 09:38:20
 * @modified by: jayden (jaydenzhao@outlook.com>)
 */
package excel

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTMLRender struct {
}

var _ IORender = &HTMLRender{}

func (r *HTMLRender) Render(to io.Writer, source interface{}) error {

	if _, err := to.Write([]byte{0xef, 0xbb, 0xbf}); err != nil {
		return err
	}
	if _, err := to.Write([]byte("<html><head><meta http-equiv='Content-Type' content='application/vnd.ms-excel; charset=utf-8' /></head><table border=1>")); err != nil {
		return err
	}
	iter, errCh := iterator(source)
	batch := 0
	for row := range iter {
		var builder strings.Builder
		if batch == 0 {
			for _, col := range row {
				builder.WriteString(fmt.Sprintf("<th>%s</th>", col))
			}
		} else {
			builder.WriteString("<tr>")
			for _, col := range row {
				builder.WriteString(fmt.Sprintf("<td>%s</td>", col))
			}
			builder.WriteString("</tr>")
		}
		if _, err := to.Write([]byte(builder.String())); err != nil {
			return err
		}
		batch++
	}
	if _, err := to.Write([]byte("</table></html>")); err != nil {
		return err
	}

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}

func NewHttpHTMLXLSRender(w http.ResponseWriter, filename string) Render {
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Add("Content-type", "application/force-download")
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fmtFileName(filename, ".xls")))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Cache-Control", "must-revalidate,post-check=0,pre-check=0")
	w.Header().Set("Expires", "0")
	w.Header().Set("Pragma", "public")

	return &defaultRender{
		w,
		&HTMLRender{},
	}
}

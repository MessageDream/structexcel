/**
 * @file: render.go
 * @created: 2021-08-04 05:21:50
 * @author: jayden (jaydenzhao@outlook.com)
 * @description: code for project
 * @last modified: 2021-08-04 05:21:54
 * @modified by: jayden (jaydenzhao@outlook.com>)
 */
package excel

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	codec "github.com/MessageDream/structcodec"
)

type IORender interface {
	Render(to io.Writer, source interface{}) error
}

type Render interface {
	Render(source interface{}) error
}

type defaultRender struct {
	w        io.Writer
	ioRender IORender
}

var _ Render = &defaultRender{}

func (r *defaultRender) Render(source interface{}) error {
	return r.ioRender.Render(r.w, source)
}

func iterator(source interface{}) (<-chan []string, <-chan error) {
	ch := make(chan []string, 0)
	errCh := make(chan error, 0)
	go func() {
		var headers []string

		gen := func(item interface{}) error {
			rows, err := decode(item, &headers)
			if err != nil {
				return err
			}
			for _, row := range rows {
				ch <- row
			}
			return nil
		}

		switch s := source.(type) {
		case <-chan interface{}:
			for item := range s {
				if err := gen(item); err != nil {
					errCh <- err
					break
				}
			}
		default:
			sType := reflect.TypeOf(s)
			switch sType.Kind() {
			case reflect.Array, reflect.Slice:
				sValue := reflect.ValueOf(s)
				sLen := sValue.Len()
				for i := 0; i < sLen; i++ {
					v := sValue.Index(i)
					var item interface{}

					if v.CanInterface() {
						item = v.Interface()
					} else {
						item = reflect.Indirect(v).Interface()
					}
					if err := gen(item); err != nil {
						errCh <- err
						break
					}
				}
			default:
				if err := gen(s); err != nil {
					errCh <- err
					break
				}
			}
		}
		close(ch)
		close(errCh)
	}()
	return ch, errCh
}

func decode(item interface{}, headers *[]string) (results [][]string, err error) {
	var result map[string]interface{}
	decoder, err := codec.NewDecoder(&codec.DecoderConfig{
		TagName: tagName,
		Result:  &result,
		Squash:  true,
	})
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(item); err != nil {
		return nil, err
	}
	var row []string
	hds := *headers
	if len(hds) == 0 {
		if orderly, ok := item.(HeaderExport); ok {
			hds = orderly.HeadersToExport()
		}
		if len(hds) == 0 {
			for k, v := range result {
				hds = append(hds, k)
				row = append(row, fmtString(v))
			}
			results = append(results, hds, row)
			*headers = hds
			return
		}
		*headers = hds
		results = append(results, hds)
	}
	//for sort
	for _, k := range hds {
		v, ok := result[k]
		if ok {
			row = append(row, fmtString(v))
		} else {
			row = append(row, "")
		}
	}
	results = append(results, row)
	return
}

func fmtString(v interface{}) string {
	if sv, ok := v.(string); ok {
		return strings.TrimLeft(sv, "+-=@ ")
	}
	return fmt.Sprintf("%v", v)
}

func fmtFileName(fileName, ext string) string {
	fExt := path.Ext(fileName)
	if fExt == ext {
		return fileName
	}
	return strings.TrimSuffix(filepath.Base(fileName), fExt) + ext
}

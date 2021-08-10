/**
 * @file: render_test.go
 * @created: 2021-08-04 06:56:04
 * @author: jayden (jaydenzhao@outlook.com)
 * @description: code for project
 * @last modified: 2021-08-04 06:56:08
 * @modified by: jayden (jaydenzhao@outlook.com>)
 */
package excel

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

type Source uint8

func (s Source) String() string {
	switch s {
	case 1:
		return "本地"
	default:
		return "远程"
	}
}

type SourceDisplay string

func (s SourceDisplay) Raw() Source {
	switch s {
	case "本地":
		return 1
	default:
		return 2
	}
}

type StampTime time.Time

type TestT struct {
	ID         string    `excel:"-"`
	M2         string    `excel:"设备唯一标识"`
	Source     Source    `excel:"来源,convert_method=SourceConvert"`
	CreateTime StampTime `excel:"创建时间,time_format=2006-01-02 15:04:05"`
}

func (t TestT) HeadersToExport() []string {
	return []string{
		"设备唯一标识",
		"来源",
		"创建时间",
	}
}

func (t TestT) SourceConvert(fieldName string, val interface{}) (interface{}, error) {
	switch s := val.(type) {
	case Source:
		return s.String(), nil
	case string:
		return SourceDisplay(s).Raw(), nil
	}
	return val, nil
}

var (
	commands = map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}

	list = []*TestT{
		{
			ID:         "xxx",
			M2:         "123",
			Source:     2,
			CreateTime: StampTime(time.Now()),
		},
		{
			M2:         "456",
			Source:     1,
			CreateTime: StampTime(time.Now()),
		},
		{
			M2:         "789",
			Source:     2,
			CreateTime: StampTime(time.Now()),
		},
	}
)

func server(path string, t *testing.T) {
	server := &http.Server{Addr: ":12345"}

	time.AfterFunc(time.Second*1, func() {
		run, ok := commands[runtime.GOOS]
		if !ok {
			t.Fatalf("don't know how to open things on %s platform", runtime.GOOS)
		}
		cmd := exec.Command(run, fmt.Sprintf("http://localhost:12345/%s", path))
		err := cmd.Start()
		if err != nil {
			t.Fatal("cmd.Start:", err)
		}
		time.AfterFunc(time.Second*2, func() {
			server.Shutdown(context.Background())
		})
	})
	err := server.ListenAndServe()
	if err != nil {
		t.Fatal("ListenAndServe:", err)
	}
}

func TestHttpCSV(t *testing.T) {
	http.HandleFunc("/csv", func(rw http.ResponseWriter, r *http.Request) {
		render := NewHttpCSVRender(rw, "test_http.csv")
		err := render.Render(list)
		if err != nil {
			rw.Write([]byte(err.Error()))
		}
	})
	server("/csv", t)
}

func TestHttpHTMLXLS(t *testing.T) {
	http.HandleFunc("/xls", func(rw http.ResponseWriter, r *http.Request) {
		render := NewHttpHTMLXLSRender(rw, "test_http.xls")
		err := render.Render(list)
		if err != nil {
			rw.Write([]byte(err.Error()))
		}
	})
	server("/xls", t)
}

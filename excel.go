/**
 * @file: excel.go
 * @created: 2021-08-03 05:03:08
 * @author: jayden (jaydenzhao@outlook.com)
 * @description: code for project
 * @last modified: 2021-08-03 05:03:14
 * @modified by: jayden (jaydenzhao@outlook.com>)
 */
package excel

import (
	codec "github.com/MessageDream/structcodec"
)

const (
	defaultTagName = "excel"
)

var (
	tagName = defaultTagName
)

func SetTagName(tag string) {
	tagName = tag
}

type Export = codec.IEncode
type Import = codec.IDecode

type HeaderExport interface {
	HeadersToExport() []string
}

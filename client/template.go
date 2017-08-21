// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"html/template"
	"path/filepath"
	"regexp"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

// 模板文件的扩展名
const templateExtension = ".html"

// 编译主题的模板。
func (client *Client) compileTemplate() error {
	funcMap := template.FuncMap{
		"strip":    stripTags,
		"html":     htmlEscaped,
		"ldate":    client.longDateFormat,
		"sdate":    client.shortDateFormat,
		"rfc3339":  rfc3339DateFormat,
		"themeURL": func(p string) string { return vars.ThemesURL(p) },
	}

	tpl, err := template.New("client").
		Funcs(funcMap).
		ParseGlob(filepath.Join(client.data.Theme.Path, "*"+templateExtension))
	if err != nil {
		return err
	}
	client.template = tpl

	return client.checkPostTemplate()
}

// 检测文章中的模板名称是否在模板中真实存在
func (client *Client) checkPostTemplate() error {
	for _, post := range client.data.Posts {
		if nil != client.template.Lookup(post.Template) {
			continue
		}

		return &data.FieldError{
			Message: "不存在",
			Field:   "template",
			File:    post.Slug,
		}
	}

	return nil
}

func rfc3339DateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(time.RFC3339)
}

func (client *Client) longDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(client.data.Config.LongDateFormat)
}

func (client *Client) shortDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(client.data.Config.ShortDateFormat)
}

// 将内容显示为 HTML 内容
func htmlEscaped(html string) interface{} {
	return template.HTML(html)
}

// 去掉所有的标签信息
var stripExpr = regexp.MustCompile("</?[^</>]+/?>")

// 过滤标签。
func stripTags(html string) string {
	return stripExpr.ReplaceAllString(html, "")
}

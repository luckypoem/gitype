// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client ...
package client

import (
	"html/template"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/mux"
)

// Client 所有数据的缓存，每次更新数据时，
// 直接声明一个新的 Client 实例，丢弃原来的 Client 即可。
type Client struct {
	path       *vars.Path
	info       *info
	mux        *mux.Mux
	etag       string
	template   *template.Template // 主题编译后的模板
	rss        []byte
	atom       []byte
	sitemap    []byte
	opensearch []byte
	patterns   []string // 记录所有的路由项，方便翻译时删除

	Created int64 // 当前数据的加载时间
	Data    *data.Data
}

// New 声明一个新的 Client 实例
func New(path *vars.Path, mux *mux.Mux) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	client := &Client{
		path:    path,
		mux:     mux,
		etag:    strconv.FormatInt(now, 10),
		Created: now,
		Data:    d,
	}
	client.info = client.newInfo()

	errFilter := func(fn func() error) {
		if err == nil {
			err = fn()
		}
	}

	errFilter(client.compileTemplate)
	errFilter(client.buildRSS)
	errFilter(client.buildAtom)
	errFilter(client.buildSitemap)
	errFilter(client.buildOpensearch)
	errFilter(client.initRoutes)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Free 释放 Client 内容
func (client *Client) Free() {
	for _, pattern := range client.patterns {
		client.mux.Remove(pattern)
	}
}

func formatUnix(unix int64, format string) string {
	t := time.Unix(unix, 0)
	return t.Format(format)
}

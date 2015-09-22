// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 用于执行typing的安装。
// 输出默认的配置文件；
// 输出默认的日志配置文件；
// 填充默认的数据到数据库；
package install

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

// OutputConfigFile 用于输出配置文件到指定的位置。
// 目前包含了日志配置文件和程序本身的配置文件。
func OutputConfigFile(logsConfigPath, configPath string) error {
	if err := ioutil.WriteFile(logsConfigPath, logFile, os.ModePerm); err != nil {
		return err
	}

	cfg := &core.Config{
		Core: &web.Config{
			HTTPS:      false,
			CertFile:   "",
			KeyFile:    "",
			Port:       "8080",
			ServerName: "typing",
			Static: map[string]string{
				"/admin": "./static/admin/",
			},
		},

		DBDSN:    "./output/main.db",
		DBPrefix: "typing_",
		DBDriver: "sqlite3",

		FrontAPIPrefix: "/api",
		AdminAPIPrefix: "/admin/api",

		ThemeDir: "./static/front/",
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, os.ModePerm)
}

// FillDB 向数据库写入初始内容。
func FillDB(db *orm.DB) error {
	if db == nil {
		return errors.New("db==nil")
	}

	// 创建表
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = tx.MultCreate(
		&models.Option{},
		&models.Comment{},
		&models.Meta{},
		&models.Post{},
		&models.Relationship{},
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	// option
	if err := fillOptions(db); err != nil {
		return err
	}

	// meta
	if err = fillMetas(db); err != nil {
		return err
	}

	// post
	post := &models.Post{
		Title:    "第一篇日志",
		Content:  "<p>这是你的第一篇日志</p>",
		Created:  time.Now().Unix(),
		Modified: time.Now().Unix(),
	}
	if _, err := db.Insert(post); err != nil {
		return err
	}

	// comment
	comment := &models.Comment{
		PostID:     1,
		Content:    "<p>沙发</p>",
		AuthorName: "游客",
		State:      models.CommentStateWaiting,
	}
	if _, err := db.Insert(comment); err != nil {
		return err
	}

	// relationship
	if _, err := db.Insert(&models.Relationship{MetaID: 1, PostID: 1}); err != nil {
		return err
	}
	if _, err := db.Insert(&models.Relationship{MetaID: 2, PostID: 1}); err != nil {
		return err
	}

	return nil
}

func fillMetas(db *orm.DB) error {
	metas := []*models.Meta{
		// cats
		{
			Name:        "default",
			Title:       "默认分类",
			Type:        models.MetaTypeCat,
			Order:       10,
			Parent:      models.MetaNoParent,
			Description: "所有添加的文章，默认添加此分类下。",
		},

		// tags
		{Name: "tag1", Title: "标签一", Type: models.MetaTypeTag, Description: "<h5>tag1</h5>"},
		{Name: "tag2", Title: "标签二", Type: models.MetaTypeTag, Description: "<h5>tag2</h5>"},
		{Name: "tag3", Title: "标签三", Type: models.MetaTypeTag, Description: "<h5>tag3</h5>"},
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if err := tx.InsertMany(metas); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return err
}

func fillOptions(db *orm.DB) error {
	opt := &core.Options{
		SiteName:    "typing blog",
		SecondTitle: "副标题",
		ScreenName:  "typing",
		Password:    core.HashPassword("123"),
		Keywords:    "typing",
		Description: "typing-极简的博客系统",
		Suffix:      ".html",

		PageSize:     20,
		DateFormat:   "2006-01-02 15:04:05",
		SidebarSize:  10,
		CommentOrder: core.CommentOrderDesc,

		PostsChangefreq: "never",
		CatsChangefreq:  "daily",
		TagsChangefreq:  "daily",
		PostsPriority:   0.9,
		CatsPriority:    0.6,
		TagsPriority:    0.4,

		Theme: "default",
	}

	maps, err := opt.ToMaps()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	sql := "INSERT INTO #options ({key},{group},{value}) VALUES(?,?,?)"
	for _, item := range maps {
		_, err := tx.Exec(true, sql, item["key"], item["group"], item["value"])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return err
}

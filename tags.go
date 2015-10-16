// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// @api put /admin/api/tags/{id}/merge 将指定的标签合并到当前标签
// @apiGroup admin
//
// @apiRequest json
// @apiParam tags array 所有需要合并的标签ID列表。
// @apiExample json
// {"tags": [1,2,3] }
//
// @apiSuccess 204 修改完成
func adminPutTagMerge(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// @api get /admin/api/tags 获取所有标签信息
// @apiGroup admin
//
// @apiRequest json
// @apiheader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam tags array 所有标签的列表
func adminGetTags(w http.ResponseWriter, r *http.Request) {
	sql := `SELECT m.{name},m.{title},m.{description},m.{id},count(r.{tagID}) AS {count}
			FROM #tags AS m
			LEFT JOIN #relationships AS r ON m.{id}=r.{tagID}
			GROUP BY m.{id}`
	rows, err := db.Query(true, sql)
	if err != nil {
		logs.Error("getTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		logs.Error("getTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"tags": maps}, nil)
}

// @api put /admin/api/tags/{id} 修改某id的标签内容
// @apiParam id int 需要修改的标签id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name        string 唯一名称
// @apiParam title       string 显示的标题
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 204 no content
//
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail  array  说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPutTag(w http.ResponseWriter, r *http.Request) {
	t := &models.Tag{}
	if !core.ReadJSON(w, r, t) {
		return
	}

	// 检测是否为空
	errs := &core.ErrorResult{Message: "格式错误", Detail: map[string]string{}}
	if len(t.Name) == 0 {
		errs.Add("name", "不能为空")
	}
	if len(t.Title) == 0 {
		errs.Add("title", "不能为空")
	}
	if errs.HasErrors() {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	var ok bool
	t.ID, ok = core.ParamID(w, r, "id")
	if !ok {
		return
	}

	// 检测是否存在同名
	titleExists, nameExists, err := tagIsExists(t)
	if err != nil {
		logs.Error("adminPutTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if titleExists {
		errs.Add("title", "与已有标签同名")
	}
	if nameExists {
		errs.Add("name", "与已有标签同名")
	}
	if errs.HasErrors() {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	if _, err := db.Update(t); err != nil {
		logs.Error("adminPutTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api post /admin/api/tags 添加新标签
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name        string 唯一名称
// @apiParam title       string 显示的标题
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 201 created
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail  array  说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPostTag(w http.ResponseWriter, r *http.Request) {
	t := &models.Tag{}
	if !core.ReadJSON(w, r, t) {
		return
	}

	errs := &core.ErrorResult{Message: "格式错误"}
	if t.ID != 0 {
		errs.Add("id", "不允许的字段")
	}
	if len(t.Title) == 0 {
		errs.Add("title", "不能为空")
	}
	if len(t.Name) == 0 {
		errs.Add("name", "不能为空")
	}
	if errs.HasErrors() {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	t.ID = 0
	titleExists, nameExists, err := tagIsExists(t)
	if err != nil {
		logs.Error("adminPostTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if titleExists {
		errs.Add("title", "已有同名字体段")
	}
	if nameExists {
		errs.Add("name", "已有同名字体段")
	}
	if errs.HasErrors() {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	if _, err := db.Insert(t); err != nil {
		logs.Error("adminPostTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, "{}", nil)
}

// @api delete /admin/api/tags/{id} 删除该id的标签
// @apiParam id int 需要删除的标签id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 204 no content
func adminDeleteTag(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("adminDeleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if _, err := tx.Delete(&models.Tag{ID: id}); err != nil {
		logs.Error("adminDeleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 删除与之对应的关联数据。
	sql := "DELETE FROM #relationships WHERE {tagID}=?"
	if _, err := tx.Exec(true, sql, id); err != nil {
		logs.Error("adminDeleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if err := tx.Commit(); err != nil {
		logs.Error("adminDeleteMeta:", err)
		tx.Rollback()
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// 是否存在相同name或是title的标签
// title返回参数表示是否有title字段相同，name返回参数表示是否有name字段相同。
func tagIsExists(t *models.Tag) (title bool, name bool, err error) {
	sql := db.Where("({name}=? OR {title}=?) AND {id}<>?", t.Name, t.Title, t.ID).
		Table("#tags")

	maps, err := sql.SelectMapString(true, "name", "title")
	if err != nil {
		return false, false, err
	}

	if len(maps) == 0 {
		return false, false, nil
	}

	for _, record := range maps {
		if record["name"] == t.Name {
			name = true
		}
		if record["title"] == t.Title {
			title = true
		}
	}

	return title, name, nil
}

// @api get /admin/api/tags/{id} 获取指定id的标签内容
// @apiParam id int 标签的id
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam id          int 	标签的id
// @apiParam name        string 标签的唯一名称，可能为空
// @apiParam title       string 标签名称
// @apiParam description string 对标签的详细描述
func adminGetTag(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	t := &models.Tag{ID: id}
	if err := db.Select(t); err != nil {
		logs.Error("adminGetTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, t, nil)
}

// 将一串标签名转换成id
// names为一种由标签名组成的字符串，名称之间由逗号分隔。
func getTagsID(names string) ([]int64, error) {
	name := strings.Split(names, ",")
	if len(name) == 0 {
		return nil, nil
	}

	cond := strings.Repeat("?,", len(name))
	sql := "SELECT {id} FROM #tags WHERE {name} IN(" + cond[:len(cond)-1] + ")"
	args := make([]interface{}, 0, len(name))
	for _, v := range name {
		args = append(args, v)
	}
	rows, err := db.Query(true, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maps, err := fetch.ColumnString(false, "id", rows)
	if err != nil {
		return nil, err
	}

	ret := make([]int64, 0, len(maps))
	for _, str := range maps {
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, num)
	}

	return ret, nil
}

// 获取与某post相关联的标签
func getPostTags(postID int64) ([]int64, error) {
	sql := `SELECT rs.{tagID} FROM #relationships AS rs
	LEFT JOIN #tags AS m ON m.{id}=rs.{tagID}
	WHERE rs.{postID}=?`
	rows, err := db.Query(true, sql, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maps, err := fetch.ColumnString(false, "tagID", rows)
	if err != nil {
		return nil, err
	}

	ret := make([]int64, 0, len(maps))
	for _, v := range maps {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, num)
	}
	return ret, nil
}
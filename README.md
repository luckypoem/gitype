gitype [![Build Status](https://travis-ci.org/caixw/gitype.svg?branch=nosql)](https://travis-ci.org/caixw/gitype)
[![Go version](https://img.shields.io/badge/Go-1.10-brightgreen.svg?style=flat)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/caixw/gitype)](https://goreportcard.com/report/github.com/caixw/gitype)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
======

基于 Git 的博客系统，具有以下特性：

1. 无数据库，通过 Git 管理发布的内容；
1. 无分类，通过标签来归类；
1. 不区分页面和普通文章；
1. 可以实时搜索内容；
1. 自动生成 RSS、Atom、Sitemap 和 Opensearch 等内容；
1. 支持 PWA；
1. 自定义主题。



演示地址： https://caixw.io



### 使用

1. 下载代码：`go get github.com/caixw/gitype`；
1. 运行 `scripts/build.sh` 编译代码（也可以直接执行 `go build` 编译，除了版本号，并无其它差别。）；
1. 执行 `gitype -init=/to/path` 输出初始的数据内容；
1. 运行 `gitype -appdir=/to/path`。

*./scripts 目录下包含了部分平台下的转换成守护进程的脚本*
*./testdata 也是一个完整的工作目录，如果不想执行 `-init` 命令初始化的话，也可以直接复制 ./testdata 的内容。*



#### 主要参数


参数名      | 值类型     | 描述
|:----------|:-----------|:-----
| preview   | bool       | 预览模式，指定此参数，可以监视用户的数据目录，一旦有更改，就会自动重新加载数据，当用户在本地创作文章时，可以指这此值。
| appdir    | string     | 指定数据目录

*具体可通过运行 `gitype -h` 查看所有的命令*



### 目录结构

appdir 的目录结构是固定的。
其中 conf 为程序的配置相关内容，包含了后台更新界面的密码，不能对外公开；
data 为博客的实际内容，包含了文章，标签，友链以及网站名称等相关的配置，
所有针对博客内容的相关设置和内容发布，都直接体现在此目录下。

```
|--- conf 程序的配置文件
|     |
|     |--- logs.xml 日志的配置文件
|     |
|     |--- web.yaml 程序的配置文件
|     |
|     |--- webhook.yaml webhook 的配置文件
|
|--- data 程序的数据目录
      |
      |--- meta 博客的一些设置项
      |     |
      |     |--- config.yaml 基本设置项，比如网站名称等
      |     |
      |     |--- tags.yaml 标签的定义
      |     |
      |     |--- links.yaml 友情链接
      |
      |--- posts 文章所在的目录
      |
      |--- raws 其它可直接通过地址访问的内容可直接放在此处
      |
      |--- themes 自定义的主题目录
            |
            |--- default 默认的主题
```



#### conf 目录下内容

conf 目录下的为程序级别的配置文件，需要重启才能使更改生效。其中：
- web.yaml 网站的启动数据信息；
- webhook.yaml 自动更新的触发条件；
- logs.xml 定义了日志的输出形式和保存路径，具体配置可参考 [logs](https://github.com/issue9/logs) 的相关文档。


##### web.yaml

参考 https://github.com/issue9/web 中的配置文件内容



##### webhook.yaml

名称        | 类型          | 描述
:-----------|:--------------|:------
url         | string        | webhooks 的接收地址
frequency   | time.Duration | webhooks 的最小更新频率
method      | string        | webhooks 接收地址的接收方法，默认为 POST
repoURL     | string        | 远程仓库的地址


#### data 目录下内容


涉及的时间均为 RFC3339 格式：2006-01-02T15:04:05Z07:00。


##### meta/config.yaml

config.yaml 指定了网站的一些基本配置情况：

名称            | 类型            | 描述
:---------------|:----------------|:------
title           | string          | 网站标题
subtitle        | string          | 网站副标题
beian           | string          | 备案号
uptime          | string          | 上线时间，符合 rfc 3339 标准的时间字符串
pageSize        | int             | 每页显示的数量
longDateFormat  | string          | 长时间的显示格式，Go 的时间格式化方式
shortDateFormat | string          | 短时间的显示格式，Go 的时间格式化方式
theme           | string          | 默认主题
type            | string          | 所有 HTML 页面的 mimetype，默认使用 text/html
icon            | Icon            | 网站的图标
menus           | []Link          | 菜单内容，格式与 links.yaml 的相同
author          | Author          | 文章的默认作者信息
license         | Link            | 文章的默认版权信息
archive         | Archive         | 存档页的相关配置
outdated        | time.Duration   | 超过此时间值，文章被标记为过时内容，显示一些提示信息
rss             | RSS             | rss 配置，若不需要，则不指定该值即可
atom            | RSS             | atom 配置，若不需要，则不指定该值即可
sitemap         | Sitemap         | sitemap 相关配置，若不需要，则不指定该值即可
opensearch      | Opensearch      | opensearch 相关配置，若不需要，则不指定该值即可
pwa             | PWA             | PWA 的相关配置，不指定，则不支持该功能
pages           | map[string]Page | 各个类型页面的一些自定义项


###### Author

名称      | 类型        | 描述
:---------|:------------|:----------
name      | string      | 名称
url       | string      | 网站地址
email     | string      | 邮箱
avatar    | string      | 头像


###### Archive

名称      | 类型        | 描述
:---------|:------------|:----------
order     | string      | 存档的排序方式，可以是：desc(默认) 和 month
type      | string      | 存档的分类方式，可以是按年：year(默认) 或是按月：month
format    | string      | 标题的格式


###### RSS

名称      | 类型     | 描述
:---------|:---------|:----------
title     | string   | 标题
size      | int      | 显示数量
url       | string   | 地址
type      | string   | 当前文件的 mimetype


###### Sitemap

名称           | 类型     | 描述
:--------------|:---------|:----------
url            | string   | Sitemap 的地址
xslURL         | string   | 为 sitemap.xml 配置的 xsl，可以为空
enableTag      | bool     | 是否把标签放到 Sitemap 中
priority       | float    | 标签页的权重
changefreq     | string   | 标签页的修改频率
postPriority   | float    | 文章页的权重
postChangefreq | string   | 文章页的修改频率
type           | string   | 当前文件的 mimetype，默认为 application/atom+xml 或是 applicatin/rss+xml


###### Opensearch

名称        | 类型     | 描述
:-----------|:---------|:----------
url         | string   | opensearch 的地址
title       | string   | 出现于 html>head>link.title 属性中
shortName   | string   | shortName 值
description | string   | description 值
longName    | string   | longName 值
image       | Icon     | image 值
type        | string   | 当前文件的 mimetype 若不指定，则使用 application/opensearchdescription+xml


###### PWA

有关 pwa 的说明，可以参考以下内容：

*https://developer.mozilla.org/zh-CN/docs/Web/Manifest*

名称            | 类型     | 描述
:---------------|:---------|:----------
serviceWorkers  | string   | 指定 service worker 文件的地址，相对于根
manifest        | Manifest | manifest.json 的相关配置


###### Manifest

名称            | 类型     | 描述
:---------------|:---------|:----------
url             | string   | manifest 的地址
type            | string   | 当前文件的 mimetype 若不指定，则使用 application/manifest+json
lang            | string   | name 值所使用的语言
name            | string   | name 值
shortName       | string   | short_name 值
startURL        | string   | start_url 值
display         | string   | display 值
description     | string   | description 值
dir             | string   | dir 值，表示文字方向
orientation     | string   | orientation 值
scope           | string   | scope 值
themeColor      | string   | theme_color 值
backgroundColor | string   | background_color 值
longName        | string   | longName 值
icons           | []Icon   | icons 值


###### Icon

名称      | 类型     | 描述
:---------|:---------|:----------
type      | string   | 图标的 mimetype
sizes     | string   | 图标的大小
url       | string   | 图标地址


###### Link

名称      | 类型     | 描述
:---------|:---------|:----------
text      | string   | 字面文字，可以不唯一
url       | string   | 对应的链接地址
title     | string   | a 标签的 title 属性。可以为空
icon      | string   | 一个 URL 或是 fontawesome 图标名称
rel       | string   | a 标签的 rel 属性
type      | string   | 指向内容的类型


###### Page
名称         | 类型     | 描述
:------------|:---------|:----------
title        | string   | 页面的 html>head>title
keywords     | string   | 页面的 html>head>meta.keywords
description  | string   | 页面的 html>head>meta.description

**部分页面可使用 %content% 占位符，分别是表示网站名称和可自由取代的内容，比如 tag 页面，%content 会用标签名代替**

**tag 和 post 页面的 keywords 和 description 是不可更改的**



##### meta/links.yaml

links.yaml 用于指定友情链接，为一个数组。每个元素均为一个 `Link`。
每个元素可以使用 [XFN](https://gmpg.org/xfn/)。



##### meta/tags.yaml

tags.yaml 用于指定所有的标签内容。为一个数组，每个元素包含以下字段：

名称      | 类型     | 描述
:---------|:---------|:----------
slug      | string   | 唯一名称，文章引用此值，地址中也使用此值
title     | string   | 字面文字，可以不唯一
color     | string   | 颜色值，在展示所有标签的页面，会以此颜色显示
content   | string   | 用于描述该标签的详细内容，可以是 **HTML**
series    | bool     | 是否为一个专题



##### posts

data/posts 为文章目录，目录层次可以按自己的习惯进行分类，系统根据是否包含 `meta.yaml`
和 `content.html` 来区分当前目录是否为一篇文章内容。比如：
```
--- posts
      +--- about
      |      |
      |      +--- meta.yaml
      |      |
      |      +--- content.html
      |
      +--- 2016
      |
      +--- 2017
            |
            +--- post1
            |      |
            |      +--- meta.yaml
            |      |
            |      +--- content.html
            |
            +--- post2
                   |
                   +--- meta.yaml
                   |
                   +--- content.html
```
其中 `/posts/about`、`/posts/2017/post2` 和 `/posts/2017/post2` 均被判断为文章。


###### meta.yaml

meta.yaml 包含了当前文章的一些细节信息。

名称      | 类型      | 描述
:---------|:----------|:----------
title     | string    | 标题
created   | string    | 创建时间，符合 rfc 3339 标准的时间字符串
modified  | string    | 修改时间，符合 rfc 3339 标准的时间字符串
tags      | string    | 关联的标签，以逗号分隔多个字符串，标签名为 meta/tags.yaml 中的 slug
summary   | string    | 摘要，同时也作为 html>head>meta.description 的内容
content   | string    | 内容
outdated  | string    | 已过时文章的提示信息
state     | string    | 状态，可以是 top、last、draft 和 default，默认为 default
image     | string    | 封面图片
author    | Author    | 作者，默认为 meta/config.yaml 中的 author 内容
license   | Link      | 版本信息，默认为 meta/config.yaml 中的 license 内容
template  | string    | 使用的模板，默认为 post
keywords  | string    | html>head>meta.keywords 标签的内容，如果为空，使用 tags
language  | string    | 语言标签
assets    | array     | 需要 PWA 缓存的信息，如果系统未启用，这些内容不启作用。



##### themes

data/themes 下为主题文件，可定义多个主题，通过 config 中的 theme 指定当前使用的主题。
主题模板语法为 [html/template](https://golang.org/pkg/html/template/)。


单一主题下，可以为文章详细页定义多个模板，通过每篇文章的 meta.yaml 可以自定义当前文章使用的模板，
默认情况下，使用 post 模板。


###### theme.yaml

定义了主题相关的一些属性。

名称         | 类型      | 描述
:------------|:----------|:----------
name         | string    | 名称
version      | string    | 版本号
description  | string    | 主题的详细描述信息
author       | Author    | 作者信息
assets       | []string  | 需要 PWA 缓存的信息，如果系统未启用，这些内容不启作用。

*如果指定了 assets 的内容，则每次更新主题内容时，必须改变版本号，PWA 根据版本号确定是否需要更新缓存的内容*



###### 错误模板

400 及以上的错误信息，均可以自定义，方式为在当前主题目录下，新建一个与错误代码相对应的 HTML 文件，
比如 400 错误，会读取 400.html 文件，以此类推。但是只能是纯 HTML 文本，不能包含模板代码。


##### raws

当访问的页面不存在时，会尝试从 raws 下访问相关内容。比如 `/abc.html`，会尝试在查找 `raws/abc.html`
文件是否存在；甚至 `/post/2016/about.htm` 这样标准的文章路由，如果文章不存在，会也访问 `raws`
目录，查看其下是否在正好相同的文件。 




### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。

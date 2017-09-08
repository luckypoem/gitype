// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

var robots = `User-agent:*
Disallow:/themes/`

// 从 app/testdata/conf.logs.xml 而来
var defaultLogsXML = `<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <info prefix="[INFO]" flag="">
        <console output="stderr" foreground="green" background="black" />
        <rotate prefix="info-" dir="./testdata/logs/" size="5M" />
    </info>
 
    <debug prefix="[DEBUG]">
        <console output="stderr" foreground="yellow" background="blue" />
        <rotate prefix="debug-" dir="./testdata/logs/" size="5M" />
    </debug>
 
    <trace prefix="[TRACE]">
        <console output="stderr" foreground="yellow" background="blue" />
        <rotate prefix="trace-" dir="./testdata/logs/" size="5M" />
    </trace>
 
    <warn prefix="[WARNNING]">
        <console output="stderr" foreground="yellow" background="blue" />
        <rotate prefix="info-" dir="./testdata/logs/" size="5M" />
    </warn>
 
    <error prefix="[ERROR]" flag="log.llongfile">
        <console output="stderr" foreground="red" background="blue" />
        <rotate prefix="error-" dir="./testdata/logs/" size="5M" />
    </error>
 
    <critical prefix="[CRITICAL]" flag="log.llongfile">
        <console output="stderr" foreground="red" background="blue" />
        <rotate prefix="critical-" dir="./testdata/logs/" size="5M" />
    </critical>
</logs>
`
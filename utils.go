package golaravel

import (
	"regexp"
	"runtime"
	"time"
)

func (g *Golaravel) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(),"$1")

	g.InfoLog.Printf("Load Time: %s took %s", name, elapsed)
}
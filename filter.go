// 过滤器

package filter

import (
	"time"
)

type Filter interface {
	Run(paramName string, paramValue interface{}) (interface{}, *Error)
}

var timeLoc *time.Location

func init() {
	timeLoc, _ = time.LoadLocation("Asia/Shanghai")
}

package aliyun

import "go.gh.ink/smsutils/v3/driver"

func init() {
	driver.Register(Name, Driver{})
}

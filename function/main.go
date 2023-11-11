package main

import (
	"fmt"

	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
	huaweicontext "huaweicloud.com/go-runtime/go-api/context"
	huaweiruntime "huaweicloud.com/go-runtime/pkg/runtime"
)

func handler(event event.ConfigEvent, ctx huaweicontext.RuntimeContext) (interface{}, error) {
	fmt.Println(event.DomainId)
	return "test", nil
}

func main() {
	huaweiruntime.Register(handler)
}

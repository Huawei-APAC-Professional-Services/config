package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
	"github.com/Huawei-APAC-Professional-Services/config-rules/service"
	huaweicontext "huaweicloud.com/go-runtime/go-api/context"
	huaweiruntime "huaweicloud.com/go-runtime/pkg/runtime"
)

const (
	CheckOnlyOneEnterpriseAdministrator string = "one_enterprise_administrator"
)

var iamService *service.ConfigIAMClient

func handler(payload []byte, ctx huaweicontext.RuntimeContext) error {
	var et event.ConfigEvent
	err := json.Unmarshal(payload, &et)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info(string(payload))
	ak := ctx.GetAccessKey()
	sk := ctx.GetSecretKey()
	token := ctx.GetToken()
	region := et.InvokingEvent.RegionId
	switch provider := et.InvokingEvent.Provider; provider {
	case "iam":
		iamService = service.NewIAMClient(ak, sk, region)
	default:
		fmt.Println("test")
	}
	switch ruletype := et.RuleParameter["rule"].(type) {
	case string:
		fmt.Printf("String: %v\n", ruletype)
	default:
		fmt.Printf("Yes: %v\n", reflect.TypeOf(ruletype))
	}
	switch rule := et.RuleParameter["rule"].(string); rule {
	case CheckOnlyOneEnterpriseAdministrator:
		result, err := iamService.HasOnlyOneEnterpriseAdministrator()
		if err != nil {
			return err
		}
		if result {
			return et.ReportComplianceStatus(event.CompliantResult, token)
		} else {
			return et.ReportComplianceStatus(event.NonCompliantResult, token)
		}
	default:
		slog.Error("unkonw customized rule")
		return errors.New("unknow customized rule")
	}
}

func main() {
	huaweiruntime.Register(handler)
}

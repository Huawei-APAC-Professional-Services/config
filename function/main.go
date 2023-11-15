package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
	"github.com/Huawei-APAC-Professional-Services/config-rules/service"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	huaweicontext "huaweicloud.com/go-runtime/go-api/context"
	huaweiruntime "huaweicloud.com/go-runtime/pkg/runtime"
)

var configClient *service.ConfigClient

const (
	CheckOnlyOneEnterpriseAdministrator string = "one_enterprise_administrator"
)

func handler(payload []byte, ctx huaweicontext.RuntimeContext) (interface{}, error) {
	var et event.ConfigEvent
	var region string
	err := json.Unmarshal(payload, &et)
	if err != nil {
		return nil, err
	}
	ak := ctx.GetAccessKey()
	sk := ctx.GetSecretKey()
	domainId := et.DomainId
	if *et.TriggerType == "period" {
		region = os.Getenv("region")
	} else {
		region = *et.InvokingEvent.RegionId
	}
	globalAuth := global.NewCredentialsBuilder().WithAk(ak).WithSk(sk).WithDomainId(*domainId).Build()
	configClient = service.NewConfigClient(globalAuth, region)
	switch rule := et.RuleParameter["rule"]["value"]; rule {
	case CheckOnlyOneEnterpriseAdministrator:
		result, err := configClient.EnsureHasOnlyOneEnterpriseAdministratorPeriodCheck(&et, region)
		if err != nil {
			return nil, err
		}
		err = result.UpdatePolicyState(ctx.GetToken())
		if err != nil {
			return nil, err
		}
		return "ok", nil
	default:
		slog.Error("unkonw customized rule")
		return nil, errors.New("unknow customized rule")
	}
}

func main() {
	huaweiruntime.Register(handler)
}

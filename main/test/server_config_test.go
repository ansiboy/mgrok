package mgrok_test

import (
	"encoding/json"
	"fmt"
	"mgrok/httpProxy"
	"testing"
)

func Test_Server_Temp(t *testing.T) {
	actionData := httpProxy.ActionData{
		Action: httpProxy.ActionDelete,
		Data: httpProxy.HTTPRedirect{
			SourceAddr: "mgrok.cn mgrok.cn",
			TargetAddr: "mgrok.cn mgrok.cn mgrok.cn mgrok.cn",
		},
	}
	data, _ := json.Marshal(actionData)
	fmt.Println(string(data))
}

package controllers

import (
	"github.com/revel/revel"
)

type MfApi struct {
	*revel.Controller
}

func (c MfApi) GetStatus() revel.Result {
	data := make(map[string]interface{})
	data["ServerStatus"] = "Running"
	data["success"] = true

	return c.RenderJSON(data)
}

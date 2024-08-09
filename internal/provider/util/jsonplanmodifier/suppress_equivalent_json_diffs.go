package jsonplanmodifier

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func SuppressEquivalentJsonDiffs() planmodifier.String {
	return suppressEquivalentJsonDiffs{}
}

type suppressEquivalentJsonDiffs struct{}

func (m suppressEquivalentJsonDiffs) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// if req.StateValue is equivalent to req.PlanValue, use the state value for resp.PlanValue
	// else, use plan value for resp.PlanValue
	resp.PlanValue = req.PlanValue

	if IsKnown(req.PlanValue) && IsKnown(req.StateValue) {
		isSame, err := isSameJson(req.PlanValue.ValueString(), req.StateValue.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error comparing JSON values", err.Error())
			return
		}

		if isSame {
			resp.PlanValue = req.StateValue
		}
	}

	return
}

func (m suppressEquivalentJsonDiffs) Description(_ context.Context) string {
	return "The value of this attribute only changes if the json string fails a deep equal check."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m suppressEquivalentJsonDiffs) MarkdownDescription(_ context.Context) string {
	return "The value of this attribute only changes if the json string fails a deep equal check."
}

func isSameJson(a string, b string) (bool, error) {
	var aobj interface{}
	var bobj interface{}

	if err := json.Unmarshal([]byte(a), &aobj); err != nil {
		return false, fmt.Errorf("error unmarshalling string: %s => %v", a, err.Error())
	}

	if err := json.Unmarshal([]byte(b), &bobj); err != nil {
		return false, fmt.Errorf("error unmarshalling string: %s => %v", b, err.Error())
	}

	return reflect.DeepEqual(aobj, bobj), nil
}

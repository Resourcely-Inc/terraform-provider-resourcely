package jsonplanmodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSuppressEquivalentJsonDiffs(t *testing.T) {
	tests := []struct {
		name              string
		planValue         string
		stateValue        string
		expectedPlanValue types.String
	}{
		{
			name:              "Test Case 1: Equivalent JSON Strings",
			planValue:         `{"foo": "bar"}`,
			stateValue:        `{"foo": "bar"}`,
			expectedPlanValue: types.StringValue(`{"foo": "bar"}`),
		},
		{
			name:              "Test Case 2: Non-equivalent JSON Strings",
			planValue:         `{"foo": "bar"}`,
			stateValue:        `{"foo": "baz"}`,
			expectedPlanValue: types.StringValue(`{"foo": "bar"}`),
		},
		{
			name:              "Test Case 3: Functionally Equivalent JSON Strings",
			planValue:         ` {"bar": "foo", "foo":   "bar"}  `,
			stateValue:        `{"foo": "bar", "bar": "foo"}`,
			expectedPlanValue: types.StringValue(`{"foo": "bar", "bar": "foo"}`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := planmodifier.StringRequest{
				PlanValue:  types.StringValue(test.planValue),
				StateValue: types.StringValue(test.stateValue),
			}
			var resp planmodifier.StringResponse
			supp := SuppressEquivalentJsonDiffs()
			supp.PlanModifyString(context.Background(), req, &resp)
			if test.expectedPlanValue != resp.PlanValue {
				t.Errorf("Expected PlanValue '%s' but got '%s'", test.expectedPlanValue.ValueString(), resp.PlanValue.ValueString())
			}
		})
	}
}

func TestIsSameJson(t *testing.T) {
	tests := []struct {
		name          string
		str1          string
		str2          string
		expectedResp  bool
		expectedError bool
	}{
		{
			name:          "Test Case 1: Equivalent JSON Strings",
			str1:          `{"foo": "bar"}`,
			str2:          `{"foo": "bar"}`,
			expectedResp:  true,
			expectedError: false,
		},
		{
			name:          "Test Case 2: Non-equivalent JSON Strings",
			str1:          `{"foo": "bar"}`,
			str2:          `{"foo": "baz"}`,
			expectedResp:  false,
			expectedError: false,
		},
		{
			name:          "Test Case 3: Functionally Equivalent JSON Strings",
			str1:          ` {"bar": "foo", "foo":   "bar"}  `,
			str2:          `{"foo": "bar", "bar": "foo"}`,
			expectedResp:  true,
			expectedError: false,
		},
		{
			name:          "Test Case 4: Invalid JSON Strings",
			str1:          `invalid json`,
			str2:          `{"foo": "bar", "bar": "foo"}`,
			expectedResp:  false,
			expectedError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := isSameJson(test.str1, test.str2)
			if test.expectedError {
				if err == nil {
					t.Errorf("Expected error parsing %s or %s but did not get one", test.str1, test.str2)
				}
			}
			if test.expectedResp != resp {
				t.Errorf("Expected response '%v' but got '%v' when evaluating isSameJSON for str1: %s and str2: %s", test.expectedResp, resp, test.str1, test.str2)
			}
		})
	}
}

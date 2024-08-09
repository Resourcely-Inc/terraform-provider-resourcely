package jsonvalidator_test

import (
	"context"
	"testing"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/provider/util/jsonvalidator"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringIsJsonValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.String
		validator validator.String
		expErrors int
	}

	testCases := map[string]testCase{
		"simple-string": {
			in:        types.StringValue("\"foo\""),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 0,
		},
		"simple-number": {
			in:        types.StringValue("1.1"),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 0,
		},
		"simple-bool": {
			in:        types.StringValue("true"),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 0,
		},
		"simple-tuple": {
			in:        types.StringValue("[1,2,3]"),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 0,
		},
		"simple-object": {
			in:        types.StringValue("{\"a\": 1,\"b\":2}"),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 0,
		},
		"missing-quotes-around-string": {
			in:        types.StringValue("foo"),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 1,
		},
		"missing-comma-in-object": {
			in:        types.StringValue("{\"a\": 1\n\"b\":2}"),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in:        types.StringNull(),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in:        types.StringUnknown(),
			validator: jsonvalidator.StringIsJSON(),
			expErrors: 0,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				ConfigValue: test.in,
			}
			res := validator.StringResponse{}
			test.validator.ValidateString(context.TODO(), req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatalf("expected %d error(s), got none", test.expErrors)
			}

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}

func TestOneOfValidator_Description(t *testing.T) {
	v := jsonvalidator.StringIsJSON()

	expected := "string must be valid JSON"
	got := v.MarkdownDescription(context.Background())

	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

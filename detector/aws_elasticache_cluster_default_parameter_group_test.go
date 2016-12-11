package detector

import (
	"reflect"
	"testing"

	"github.com/wata727/tflint/config"
	"github.com/wata727/tflint/issue"
)

func TestDetectAwsElastiCacheClusterDefaultParameterGroup(t *testing.T) {
	cases := []struct {
		Name   string
		Src    string
		Issues []*issue.Issue
	}{
		{
			Name: "default.redis3.2 is default parameter group",
			Src: `
resource "aws_elasticache_cluster" "cache" {
    parameter_group_name = "default.redis3.2"
}`,
			Issues: []*issue.Issue{
				&issue.Issue{
					Type:    "NOTICE",
					Message: "\"default.redis3.2\" is default parameter group. You cannot edit it.",
					Line:    3,
					File:    "test.tf",
				},
			},
		},
		{
			Name: "application3.2 is not default parameter group",
			Src: `
resource "aws_elasticache_cluster" "cache" {
    parameter_group_name = "application3.2"
}`,
			Issues: []*issue.Issue{},
		},
	}

	for _, tc := range cases {
		var issues = []*issue.Issue{}
		TestDetectByCreatorName(
			"CreateAwsElastiCacheClusterDefaultParameterGroupDetector",
			tc.Src,
			config.Init(),
			config.Init().NewAwsClient(),
			&issues,
		)

		if !reflect.DeepEqual(issues, tc.Issues) {
			t.Fatalf("Bad: %s\nExpected: %s\n\ntestcase: %s", issues, tc.Issues, tc.Name)
		}
	}
}

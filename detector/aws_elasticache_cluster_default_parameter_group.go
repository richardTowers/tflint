package detector

import (
	"fmt"
	"regexp"

	"github.com/wata727/tflint/issue"
)

type AwsElastiCacheClusterDefaultParameterGroupDetector struct {
	*Detector
}

func (d *Detector) CreateAwsElastiCacheClusterDefaultParameterGroupDetector() *AwsElastiCacheClusterDefaultParameterGroupDetector {
	return &AwsElastiCacheClusterDefaultParameterGroupDetector{d}
}

func (d *AwsElastiCacheClusterDefaultParameterGroupDetector) Detect(issues *[]*issue.Issue) {
	for filename, list := range d.ListMap {
		for _, item := range list.Filter("resource", "aws_elasticache_cluster").Items {
			parameterGroupToken, err := hclLiteralToken(item, "parameter_group_name")
			if err != nil {
				d.Logger.Error(err)
				continue
			}
			parameterGroup, err := d.evalToString(parameterGroupToken.Text)
			if err != nil {
				d.Logger.Error(err)
				continue
			}

			if d.isDefaultCacheParameterGroup(parameterGroup) {
				issue := &issue.Issue{
					Type:    "NOTICE",
					Message: fmt.Sprintf("\"%s\" is default parameter group. You cannot edit it.", parameterGroup),
					Line:    parameterGroupToken.Pos.Line,
					File:    filename,
				}
				*issues = append(*issues, issue)
			}
		}
	}
}

func (d *AwsElastiCacheClusterDefaultParameterGroupDetector) isDefaultCacheParameterGroup(s string) bool {
	return regexp.MustCompile("^default").Match([]byte(s))
}

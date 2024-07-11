package phpCsFixer

import (
	"errors"

	"github.com/bensaufley/toSarif/formats/sarif"
	"github.com/bensaufley/toSarif/util"
)

func (s *PhpCsFixerSchemaJson) ToSarif() (*sarif.Sarif22SchemaJson, error) {
	if s == nil {
		return nil, errors.New("nil schema")
	}

	results := []sarif.Result{}
	ruleSlice := []string{}
	fileSlice := []string{}
	for _, file := range s.Files {
		fileName := "file://" + file.Name
		var fileIndex int
		fileSlice, fileIndex = util.AddUnique(fileSlice, fileName)
		for _, fixer := range file.AppliedFixers {
			rule := fixer
			var ruleIndex int
			ruleSlice, ruleIndex = util.AddUnique(ruleSlice, rule)
			results = append(results, sarif.Result{
				RuleId:    &rule,
				RuleIndex: ruleIndex,
				Locations: []sarif.Location{
					{
						PhysicalLocation: &sarif.PhysicalLocation{
							ArtifactLocation: &sarif.ArtifactLocation{
								Uri:   &fileName,
								Index: fileIndex,
							},
						},
					},
				},
			})
		}
	}

	files := make([]sarif.Artifact, len(fileSlice))
	for i, file := range fileSlice {
		files[i] = sarif.Artifact{
			Location: &sarif.ArtifactLocation{
				Uri: &file,
			},
		}
	}

	rules := make([]sarif.ReportingDescriptor, len(ruleSlice))
	for i, rule := range ruleSlice {
		var namespace *string
		var message *string
		for ns, rules := range Rules {
			if msg, ok := rules[rule]; ok {
				namespace = &ns
				message = &msg
				break
			}
		}
		var ruleUri *string
		if namespace != nil {
			uri := "https://github.com/PHP-CS-Fixer/PHP-CS-Fixer/blob/master/doc/rules/" + *namespace + "/" + rule + ".rst"
			ruleUri = &uri
		}
		var desc *sarif.MultiformatMessageString
		if message != nil {
			desc = &sarif.MultiformatMessageString{
				Text: *message,
			}
		}
		rules[i] = sarif.ReportingDescriptor{
			Id:               rule,
			HelpUri:          ruleUri,
			ShortDescription: desc,
		}
	}

	schema := "http://json.schemastore.org/sarif-2.1.0-rtm.5"
	toolVersion := "unknown"
	infoUri := "https://cs.symfony.com/"
	sarif := &sarif.Sarif22SchemaJson{
		Schema: &schema,
		Runs: []sarif.Run{
			{
				Artifacts: files,
				Tool: sarif.Tool{
					Driver: sarif.ToolComponent{
						Name:           "php-cs-fixer",
						Version:        &toolVersion,
						InformationUri: &infoUri,
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}
	return sarif, nil
}

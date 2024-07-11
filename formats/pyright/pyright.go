package pyright

import (
	"errors"

	"github.com/bensaufley/toSarif/formats/sarif"
	"github.com/bensaufley/toSarif/util"
)

func severityToSarif(severity Severity) sarif.ResultLevel {
	switch severity {
	case Error:
		return sarif.ResultLevelError
	case Warning:
		return sarif.ResultLevelWarning
	case Information:
		return sarif.ResultLevelNote
	}
	return sarif.ResultLevelNone
}

func (s *Schema) ToSarif() (*sarif.Sarif22SchemaJson, error) {
	if s == nil {
		return nil, errors.New("nil schema")
	}

	schema := "http://json.schemastore.org/sarif-2.1.0-rtm.5"

	pyrightVersion := s.Version

	results := make([]sarif.Result, len(s.GeneralDiagnostics))
	artifacts := []string{}
	rules := []string{}

	for i, diag := range s.GeneralDiagnostics {
		file := "file://" + diag.File
		message := diag.Message
		startLine := diag.Range.Start.Line
		startColumn := diag.Range.Start.Character
		endLine := diag.Range.End.Line
		endColumn := diag.Range.End.Character

		var fileIndex int
		var ruleIndex int
		artifacts, fileIndex = util.AddUnique(artifacts, file)
		if diag.Rule != nil {
			rules, ruleIndex = util.AddUnique(rules, *diag.Rule)
		}

		results[i] = sarif.Result{
			Level: severityToSarif(diag.Severity),
			Message: sarif.Message{
				Text: &message,
			},
			RuleId:    diag.Rule,
			RuleIndex: ruleIndex,
			Locations: []sarif.Location{
				{
					PhysicalLocation: &sarif.PhysicalLocation{
						ArtifactLocation: &sarif.ArtifactLocation{
							Uri:   &file,
							Index: fileIndex,
						},
						Region: &sarif.Region{
							StartLine:   &startLine,
							StartColumn: &startColumn,
							EndLine:     &endLine,
							EndColumn:   &endColumn,
						},
					},
				},
			},
		}
	}

	arts := make([]sarif.Artifact, len(artifacts))
	for i, artifact := range artifacts {
		arts[i] = sarif.Artifact{
			Location: &sarif.ArtifactLocation{
				Uri: &artifact,
			},
		}
	}

	rulesArr := make([]sarif.ReportingDescriptor, len(rules))
	helpUri := "https://github.com/microsoft/pyright/blob/main/docs/configuration.md"
	for i, rule := range rules {
		desc, ok := Rules[rule]
		var description *sarif.MultiformatMessageString
		if ok {
			description = &sarif.MultiformatMessageString{
				Text: desc,
			}
		}
		rulesArr[i] = sarif.ReportingDescriptor{
			Id:               rule,
			HelpUri:          &helpUri,
			ShortDescription: description,
		}
	}

	toolUri := "https://github.com/microsoft/pyright"
	sarif := &sarif.Sarif22SchemaJson{
		Version: "2.1.0",
		Schema:  &schema,
		Runs: []sarif.Run{
			{
				Artifacts: arts,
				Tool: sarif.Tool{
					Driver: sarif.ToolComponent{
						Name:           "pyright",
						Version:        &pyrightVersion,
						Rules:          rulesArr,
						InformationUri: &toolUri,
					},
				},
				Results: results,
			},
		},
	}
	return sarif, nil
}

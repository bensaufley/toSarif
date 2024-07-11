package pyright

type Summary struct {
	FilesAnalyzed    int
	ErrorCount       int
	WarningCount     int
	InformationCount int
	TimeInSec        float32
}

type RangeEnd struct {
	Line      int
	Character int
}

type Range struct {
	Start RangeEnd
	End   RangeEnd
}

type Severity = string

const (
	Error       Severity = "error"
	Warning     Severity = "warning"
	Information Severity = "information"
)

type Diagnostic struct {
	File     string
	Severity Severity
	Message  string
	Rule     *string
	Range    Range
}

type Schema struct {
	Version            string
	Time               string
	GeneralDiagnostics []Diagnostic
	Summary            Summary
}

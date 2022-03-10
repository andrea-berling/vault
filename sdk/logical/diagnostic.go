package logical

// These severities map to the tfdiags.Severity values, plus an explicit
// unknown in case that enum grows without us noticing here.
const (
	DiagnosticSeverityUnknown = "unknown"
	DiagnosticSeverityError   = "error"
	DiagnosticSeverityWarning = "warning"
)

// DiagnosticResponse is a wrapper to return diagnostics to HTTP calls
// and a total hack
type DiagnosticResponse struct {
	FormatVersion string `json:"format_version"`

	// We include some summary information that is actually redundant
	// with the detailed diagnostics, but avoids the need for callers
	// to re-implement our logic for deciding these.
	Valid        bool         `json:"valid"`
	ErrorCount   int          `json:"error_count"`
	WarningCount int          `json:"warning_count"`
	Diagnostics  []Diagnostic `json:"diagnostics"`
}

// Diagnostics is a list of diagnostics. Diagnostics is intended to be used
// where a Go "error" might normally be used, allowing richer information
// to be conveyed (more context, support for warnings).
//
// A nil Diagnostics is a valid, empty diagnostics list, thus allowing
// heap allocation to be avoided in the common case where there are no
// diagnostics to report at all.
type Diagnostics []Diagnostic

// Diagnostic represents any tfdiags.Diagnostic value. The simplest form has
// just a severity, single line summary, and optional detail. If there is more
// information about the source of the diagnostic, this is represented in the
// range field.
type Diagnostic struct {
	Severity string             `json:"severity"`
	Summary  string             `json:"summary"`
	Detail   string             `json:"detail"`
	Address  string             `json:"address,omitempty"`
	Range    *DiagnosticRange   `json:"range,omitempty"`
	Snippet  *DiagnosticSnippet `json:"snippet,omitempty"`
}

// Pos represents a position in the source code.
type Pos struct {
	// Line is a one-based count for the line in the indicated file.
	Line int `json:"line"`

	// Column is a one-based count of Unicode characters from the start of the line.
	Column int `json:"column"`

	// Byte is a zero-based offset into the indicated file.
	Byte int `json:"byte"`
}

// DiagnosticRange represents the filename and position of the diagnostic
// subject. This defines the range of the source to be highlighted in the
// output. Note that the snippet may include additional surrounding source code
// if the diagnostic has a context range.
//
// The Start position is inclusive, and the End position is exclusive. Exact
// positions are intended for highlighting for human interpretation only and
// are subject to change.
type DiagnosticRange struct {
	Filename string `json:"filename,omitempty"`
	Start    *Pos   `json:"start,omitempty"`
	End      *Pos   `json:"end,omitempty"`
}

// DiagnosticSnippet represents source code information about the diagnostic.
// It is possible for a diagnostic to have a source (and therefore a range) but
// no source code can be found. In this case, the range field will be present and
// the snippet field will not.
type DiagnosticSnippet struct {
	// Context is derived from HCL's hcled.ContextString output. This gives a
	// high-level summary of the root context of the diagnostic: for example,
	// the resource block in which an expression causes an error.
	Context *string `json:"context"`

	// Code is a possibly-multi-line string of Terraform configuration, which
	// includes both the diagnostic source and any relevant context as defined
	// by the diagnostic.
	Code string `json:"code"`

	// StartLine is the line number in the source file for the first line of
	// the snippet code block. This is not necessarily the same as the value of
	// Range.Start.Line, as it is possible to have zero or more lines of
	// context source code before the diagnostic range starts.
	StartLine int `json:"start_line"`

	// HighlightStartOffset is the character offset into Code at which the
	// diagnostic source range starts, which ought to be highlighted as such by
	// the consumer of this data.
	HighlightStartOffset int `json:"highlight_start_offset"`

	// HighlightEndOffset is the character offset into Code at which the
	// diagnostic source range ends.
	HighlightEndOffset int `json:"highlight_end_offset"`

	// Values is a sorted slice of expression values which may be useful in
	// understanding the source of an error in a complex expression.
	Values []DiagnosticExpressionValue `json:"values"`
}

// DiagnosticExpressionValue represents an HCL traversal string (e.g.
// "var.foo") and a statement about its value while the expression was
// evaluated (e.g. "is a string", "will be known only after apply"). These are
// intended to help the consumer diagnose why an expression caused a diagnostic
// to be emitted.
type DiagnosticExpressionValue struct {
	Traversal string `json:"traversal"`
	Statement string `json:"statement"`
}

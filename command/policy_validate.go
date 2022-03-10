package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/posener/complete"
)

type Output struct {
	FormatVersion string `json:"format_version"`

	// We include some summary information that is actually redundant
	// with the detailed diagnostics, but avoids the need for callers
	// to re-implement our logic for deciding these.
	Valid        bool               `json:"valid"`
	ErrorCount   int                `json:"error_count"`
	WarningCount int                `json:"warning_count"`
	Diagnostics  []vault.Diagnostic `json:"diagnostics"`
}

var (
	_ cli.Command             = (*PolicyValidateCommand)(nil)
	_ cli.CommandAutocomplete = (*PolicyValidateCommand)(nil)
)

type PolicyValidateCommand struct {
	*BaseCommand
}

func (c *PolicyValidateCommand) Synopsis() string {
	return "Validates the syntax of a policy on disk"
}

func (c *PolicyValidateCommand) Help() string {
	helpText := `
Usage: vault policy validate [options] PATH

  Verifies that a local policy file is syntactically valid. If a
  problem is found, this command will output any discovered errors and
  their line numbers. The output will be in a human-readable format 
  unless the -json flag is used, which provides a more detailed 
  diagnostic in JSON format.

  Validate the local policy file "my-policy.hcl":

      $ vault policy validate my-policy.hcl

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PolicyValidateCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetOutputFormat)
}

func (c *PolicyValidateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictFiles("*.hcl")
}

func (c *PolicyValidateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PolicyValidateCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	}

	var diags vault.Diagnostics

	switch {
	case args[0] == "-":
		b, err := readFromStdin()
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		_, diags = vault.ParseACLPolicyReturnDiagnostics(namespace.RootNamespace, string(b), "")
	default:
		for _, path := range args {
			b, err := readFromFile(path)
			if err != nil {
				c.UI.Error(err.Error())
				return 1
			}
			_, d := vault.ParseACLPolicyReturnDiagnostics(namespace.RootNamespace, string(b), filepath.Base(path))
			diags = append(diags, d...)
		}
	}

	// Write out the diagnostics
	if len(diags) != 0 {
		OutputDiagnostics(c, formOutput(diags))
		return 1 // we got a diagnostic, so return 1 regardless of the outcome of OutputDiagnostics above
	}

	c.UI.Output("Success! The policy is valid.")
	return 0
}

func formOutput(diagnostics vault.Diagnostics) Output {
	var (
		errorCount   int
		warningCount int
	)

	for _, d := range diagnostics {
		if d.Severity == vault.DiagnosticSeverityError {
			errorCount++
		} else {
			warningCount++
		}
	}

	return Output{
		FormatVersion: "1",
		Valid:         errorCount == 0,
		ErrorCount:    errorCount,
		WarningCount:  warningCount,
		Diagnostics:   diagnostics,
	}
}

func readFromFile(path string) ([]byte, error) {
	// Get the filepath, accounting for ~ and stuff
	p, err := homedir.Expand(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path: %w", err)
	}

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("error reading source file: %w", err)
	}

	return b, nil
}

func readFromStdin() ([]byte, error) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("couldn't read from stdin: %w", err)
	}

	return b, nil
}

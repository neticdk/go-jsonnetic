package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/neticdk/go-jsonnetic/internal/cmd/utils"
	"github.com/neticdk/go-jsonnetic/internal/jsonneticcli"
	"github.com/neticdk/go-jsonnetic/pkg/jsonnetic"
	"github.com/neticdk/go-stdlib/file"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	AppName   = "jsonnetic"
	ShortDesc = "Jsonnetic is a jsonnet implementation with a extra native functions"
)

var LongDesc = fmt.Sprintf("Jsonnetic adds the following native functions to jsonnet:\n\n%s", utils.PrettyFuncList())

func newRootCmd(ac *jsonneticcli.Context) *cobra.Command {
	o := &rootOptions{}
	c := cmd.NewRootCommand(ac.EC).
		Build()

	c.Args = cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)
	c.RunE = func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		if err := o.Complete(ctx, ac); err != nil {
			return err
		}
		if err := o.Validate(ctx, ac); err != nil {
			return err
		}
		return o.Run(ctx, ac)
	}

	// Temporarily hide the output flag, until we can use go-common v0.13.0 where it is optional
	_ = c.PersistentFlags().MarkHidden("output")
	o.bindFlags(c.Flags(), ac)

	return c
}

type rootOptions struct {
	// filename is the file to generate the output from
	filename string
	// Specify an additional library search dir (right-most wins)
	jpath []string
	// Specify the number of allowed stack frames, if not set, the default is used from the go-jsonnet implementation
	maxStack int
	// output-file is the file to write the output to
	outputFile string
	// outputMulti is the directory to write the output to for multi-file output
	outputMulti string
	// createDirs is a flag to create the output directories if they do not exist
	createDirs bool
	// StringOutput is used to configure jsonnet VM to return string output
	stringOutput bool
}

func (o *rootOptions) bindFlags(f *pflag.FlagSet, _ *jsonneticcli.Context) {
	f.StringSliceVarP(&o.jpath, "jpath", "J", nil, "Add a library search directory (rightmost takes precedence)")
	f.IntVarP(&o.maxStack, "max-stack", "s", 0, "Set the maximum number of stack frames. Uses the go-jsonnet default if not specified")
	f.StringVarP(&o.outputFile, "output-file", "", "", "Write output to the specified file. Defaults to stdout")
	f.StringVarP(&o.outputMulti, "multi", "m", "", "Write multi-file output to the specified directory.")
	f.BoolVarP(&o.createDirs, "create-output-dirs", "c", false, "Create output directories if they do not exist.")
	f.BoolVarP(&o.stringOutput, "string", "S", false, "Expect a string, manifest as plain text")
}

func (o *rootOptions) Complete(_ context.Context, ac *jsonneticcli.Context) error {
	// Set the filename
	o.filename = ac.EC.CommandArgs[0]

	if o.outputMulti != "" {
		if !strings.HasSuffix(o.outputMulti, "/") {
			o.outputMulti += "/"
		}
	}
	return nil
}

func (o *rootOptions) Validate(_ context.Context, _ *jsonneticcli.Context) error {
	// Check if the filename exists
	if exists, _ := file.Exists(o.filename); !exists {
		return &cmd.InvalidArgumentError{
			Val:     o.filename,
			Context: "filename does not exist, or is not accessible",
		}
	}
	// Check if jpath directories exist
	for _, jpath := range o.jpath {
		if exists, _ := file.Exists(jpath); !exists {
			return &cmd.InvalidArgumentError{
				Flag:    "jpath",
				Val:     jpath,
				Context: "jpath does not exist, or is not accessible",
			}
		}
	}

	return nil
}

func (o *rootOptions) Run(_ context.Context, _ *jsonneticcli.Context) error {
	// Get the jsonnet VM
	vm := jsonnetic.MakeVM(o.jpath, o.maxStack)

	// Set the string output flag
	vm.StringOutput = o.stringOutput

	var err error
	var output string
	var outputDict map[string]string

	if o.outputMulti != "" {
		outputDict, err = vm.EvaluateFileMulti(o.filename)
	} else {
		output, err = vm.EvaluateFile(o.filename)
	}
	if err != nil {
		return &cmd.GeneralError{
			Message: "failed to evaluate jsonnet file",
			Err:     err,
		}
	}

	// Write the output to multi files or a single file
	if o.outputMulti != "" {
		err = utils.WriteMultiOutputFiles(outputDict, o.outputMulti, o.outputFile, o.createDirs)
	} else {
		err = utils.WriteOutputFile(output, o.outputFile, o.createDirs)
	}
	if err != nil {
		return &cmd.GeneralError{Err: err}
	}

	return nil
}

// Execute runs the root command and returns the exit code
func Execute(version string) int {
	ec := cmd.NewExecutionContext(AppName, ShortDesc, version)
	ac := jsonneticcli.NewContext(ec)
	ec.LongDescription = LongDesc
	// rootCmd represents the base command when called without any subcommands
	rootCmd := newRootCmd(ac)
	err := rootCmd.Execute()
	_ = ec.Spinner.Stop()
	if err != nil {
		ec.ErrorHandler.HandleError(err)
		return 1
	}
	return 0
}

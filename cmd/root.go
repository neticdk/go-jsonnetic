package cmd

import (
	"context"
	"fmt"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	clierrors "github.com/neticdk/go-common/pkg/cli/errors"
	"github.com/neticdk/go-common/pkg/file"
	"github.com/neticdk/go-jsonnetic/internal/cmd/utils"
	"github.com/neticdk/go-jsonnetic/internal/jsonneticcli"
	"github.com/neticdk/go-jsonnetic/pkg/jsonnetic"
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
	evalMulti   bool
	// createDirs is a flag to create the output directories if they do not exist
	createDirs bool
}

func (o *rootOptions) bindFlags(f *pflag.FlagSet, _ *jsonneticcli.Context) {
	f.StringSliceVarP(&o.jpath, "jpath", "J", nil, "Add a library search directory (rightmost takes precedence)")
	f.IntVarP(&o.maxStack, "max-stack", "s", 0, "Set the maximum number of stack frames. Uses the go-jsonnet default if not specified")
	f.StringVarP(&o.outputFile, "output-file", "o", "", "Write output to the specified file. Defaults to stdout")
	f.StringVarP(&o.outputMulti, "multi", "m", "", "Write multi-file output to the specified directory.")
	f.BoolVarP(&o.createDirs, "create-output-dirs", "c", false, "Create output directories if they do not exist.")
}

func (o *rootOptions) Complete(_ context.Context, ac *jsonneticcli.Context) error {
	// Set the filename
	o.filename = ac.EC.CommandArgs[0]

	if o.outputMulti != "" {
		if o.outputMulti[len(o.outputMulti)-1] != '/' {
			o.outputMulti += "/"
		}
	}
	return nil
}

func (o *rootOptions) Validate(_ context.Context, _ *jsonneticcli.Context) error {
	// Check if the filename exists
	if exists, _ := file.Exists(o.filename); !exists {
		return &clierrors.InvalidArgumentError{
			Val:     o.filename,
			Context: "filename does not exist, or is not accessible",
		}
	}
	// Check if jpath directories exist
	for _, jpath := range o.jpath {
		if exists, _ := file.Exists(jpath); !exists {
			return &clierrors.InvalidArgumentError{
				Flag:    "jpath",
				Val:     jpath,
				Context: "jpath does not exist, or is not accessible",
			}
		}
	}

	if o.outputMulti != "" {
		o.evalMulti = true
	}

	return nil
}

func (o *rootOptions) Run(_ context.Context, _ *jsonneticcli.Context) error {
	// Get the jsonnet VM
	vm := jsonnetic.MakeVM(o.jpath, o.maxStack)

	var err error
	var output string
	var outputDict map[string]string

	if o.evalMulti {
		outputDict, err = vm.EvaluateFileMulti(o.filename)
	} else {
		output, err = vm.EvaluateFile(o.filename)
	}
	if err != nil {
		return &clierrors.GeneralError{
			Message: "failed to evaluate jsonnet file",
			Err:     err,
		}
	}

	// Write output JSON.
	if o.evalMulti {
		err = utils.WriteMultiOutputFiles(outputDict, o.outputMulti, o.outputFile, o.createDirs)
		if err != nil {
			return &clierrors.GeneralError{
				Err: err,
			}
		}
	} else {
		err = utils.WriteOutputFile(output, o.outputFile, o.createDirs)
		if err != nil {
			return &clierrors.GeneralError{
				Err: err,
			}
		}
	}

	return nil
}

// Execute runs the root command and returns the exit code
func Execute(version string) int {
	ec := cmd.NewExecutionContext(AppName, ShortDesc, version)
	ac := jsonneticcli.NewContext(ec)
	ec.LongDescription = LongDesc
	ec.PFlags.OutputFormatEnabled = false
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

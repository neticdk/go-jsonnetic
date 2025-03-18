package utils

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/neticdk/go-jsonnetic/pkg/jsonnetic/native"
)

const (
	FileModeNewDirectory = 0o750
	FileModeNewFile      = 0o640
)

// PrettyFuncList returns a pretty list of all available functions.
func PrettyFuncList() string {
	funcList := ""
	for _, f := range native.Funcs() {
		funcList += fmt.Sprintf("- %s\n", f.Name)
	}
	return funcList
}

// WriteMultiOutputFiles writes the output to multiple files and if enabled creates directories.
func WriteMultiOutputFiles(output map[string]string, outputDir, outputFile string, createDirs bool) (err error) { //nolint:revive
	// If multiple file output is used, then iterate over each string from
	// the sequence of strings returned by jsonnet_evaluate_snippet_multi,
	// construct pairs of filename and content, and write each output file.

	manifest := os.Stdout
	if outputFile != "" {
		manifest, err = os.Create(outputFile)
		if err != nil {
			return err
		}
		defer func() {
			if ferr := manifest.Close(); ferr != nil {
				err = ferr
			}
		}()
	}

	// Create a sorted list of outputPaths to ensure deterministic output.
	outputPaths := slices.Collect(maps.Keys(output))
	sort.Strings(outputPaths)
	for _, path := range outputPaths {
		newContent := output[path]
		filename := filepath.Join(outputDir, path)

		_, err = manifest.WriteString(filename)
		if err != nil {
			return err
		}

		_, err = manifest.WriteString("\n")
		if err != nil {
			return err
		}

		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			existingContent, err := os.ReadFile(filename)
			if err != nil {
				return err
			}
			if string(existingContent) == newContent {
				// Do not bump the timestamp on the file if its content is
				// the same. This may trigger other tools (e.g. make) to do
				// unnecessary work.
				continue
			}
		}
		if createDirs {
			if err = os.MkdirAll(filepath.Dir(filename), FileModeNewDirectory); err != nil {
				return err
			}
		}

		err = os.WriteFile(filename, []byte(newContent), FileModeNewFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteOutputFile writes the output to the given file, creating directories
// if requested, and printing to stdout instead if the outputFile is "".
func WriteOutputFile(output string, outputFile string, createDirs bool) (err error) { //nolint:revive
	if outputFile == "" {
		fmt.Print(output)
		return nil
	}

	if createDirs {
		if err = os.MkdirAll(filepath.Dir(outputFile), FileModeNewDirectory); err != nil {
			return err
		}
	}

	err = os.WriteFile(outputFile, []byte(output), FileModeNewFile)
	return err
}

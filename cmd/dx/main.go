package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var nameSpaces map[string]*cobra.Command
var DXHome string

var RootCmd = &cobra.Command{
	Use:     "dx",
	Version: Version,
	Long:    "Simple command runner/automation tool for a better developer experience",
}

func main() {

	nameSpaces = make(map[string]*cobra.Command)

	DXHome = GetDXHome()

	if DXHome == "" {
		fmt.Fprintln(os.Stderr, "ERROR: DXHome not found in any of ['HOME/DXHome', './DXHome', ENV DXHOME='']")
		os.Exit(1)
	}

	var commandDir string

	osCommandDir := filepath.Join(DXHome, fmt.Sprintf("%s_cmd", runtime.GOOS))
	crossPlatformCommandDir := filepath.Join(DXHome, "cmd")

	if fileExists(osCommandDir) {
		commandDir = osCommandDir
	} else if fileExists(crossPlatformCommandDir) {
		commandDir = crossPlatformCommandDir
	} else {
		fmt.Fprintln(os.Stderr, "ERROR: Command directory not found in either of: ", osCommandDir, crossPlatformCommandDir)
		os.Exit(1)
	}

	nameSpaces[commandDir] = RootCmd

	if err := createCommands(commandDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// I want the completion command hidden from usage
	RootCmd.CompletionOptions.HiddenDefaultCmd = true

	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func GetDXHome() string {

	// Use the ENV variable if set
	dxhome := os.Getenv("DXHOME")
	if dxhome != "" {
		return dxhome
	}

	// Otherwise prefer a DXHome in the current directory
	if fileExists("./DXHome") {
		return "./DXHome"
	}

	// Finally check to see if ~/DXHome exists and use that if it does.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if fileExists(filepath.Join(home, "DXHome")) {
		return filepath.Join(home, "DXHome")
	}

	return "" // Could not find a DXHome, let's default to nothing
}

func createCommands(dxhome string) error {
	err := filepath.Walk(dxhome,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Are we at the root or does the path (dir or file) begin with an underscore?
			if dxhome == path || strings.HasPrefix(filepath.Base(path), "_") {
				return nil // Skip it
			}

			var cmd *cobra.Command

			base := filepath.Base(path)
			namespace := filepath.Dir(path)

			log.Printf("base = '%s'\n", base)

			if info.IsDir() {
				// Directories, create name spaces in the command tree
				cmd = &cobra.Command{Use: base}
				nameSpaces[path] = cmd
			} else {
				cmd = &cobra.Command{
					Use: base,
					Run: func(cmd *cobra.Command, args []string) {

						// Run the command pointed at by path and pass any additional arguments to it
						c := exec.Command(path, args...)
						c.Stdout = os.Stdout
						c.Stderr = os.Stderr

						// make sure DXHOME is set so scripts can use it (for sourcing, docs, configs etc).
						c.Env = append(os.Environ(),
							fmt.Sprintf("DXHOME=%s", DXHome),
						)

						err := c.Run()

						if err != nil {
							fmt.Println(err.Error())
							return
						}
					},
				}
			}

			// Useful to override the built in help. I recommend only doing this at the top level
			// Since it has funny behaviour when doing this in a name space (it works though, but usage is broken)
			if base == "help" {
				nameSpaces[namespace].SetHelpCommand(cmd)
			}

			nameSpaces[namespace].AddCommand(cmd)

			return nil
		})

	return err
}

func fileExists(path string) bool {

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

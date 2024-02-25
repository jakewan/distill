package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/cbsinteractive/jakewan/altflag"
)

type cmdFSDirsize struct {
	flags       altflag.FlagSet
	deps        Dependencies
	startingDir string
	verbose     bool
}

func newCmdFSDirsize(deps Dependencies) runner {
	cmd := cmdFSDirsize{
		flags: altflag.NewFlagSet("dirsize"),
		deps:  deps,
	}
	stringTarget, err := cmd.flags.StringVar(
		&cmd.startingDir,
		string(argNameStartingDir),
		string(argShortFlagStartingDir),
		string(argUsageStartingDir))
	if err != nil {
		panic(err)
	}
	stringTarget.SetDefault("")
	boolTarget, err := cmd.flags.BoolVar(
		&cmd.verbose,
		string(argNameVerbose),
		string(argShortFlagVerbose),
		string(argUsageVerbose))
	if err != nil {
		panic(err)
	}
	boolTarget.SetDefault(false)
	return &cmd
}

func (cmd *cmdFSDirsize) name() string {
	return "dirsize"
}

func (cmd *cmdFSDirsize) init(args []string) error {
	if err := cmd.flags.Parse(args); err != nil {
		return err
	}
	return nil
}

func (cmd *cmdFSDirsize) run(ctx context.Context) error {
	const outputHeaderTemplate string = "%-36s  Size in bytes\n"
	const outputLineTemplate string = "%-36s  %d\n"

	type fileEntryInfo struct {
		absPath     string
		sizeInBytes int64
	}
	type dirEntryInfo struct {
		absPath     string
		sizeInBytes int64
	}
	rootFiles := map[string]fileEntryInfo{}
	rootDirs := map[string]*dirEntryInfo{}

	processDir := func(path string) error {
		parentDir := filepath.Dir(path)
		if abs, err := filepath.Abs(path); err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		} else {
			if parentDir == "." {
				rootDirs[path] = &dirEntryInfo{absPath: abs}
			}
		}
		return nil
	}

	processFile := func(path string, d fs.DirEntry) error {
		abs, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("error getting absolute path: %w", err)
		}
		fileType := d.Type()
		if fileType.IsRegular() {
			if info, err := d.Info(); err != nil {
				return fmt.Errorf("error getting file info: %w", err)
			} else if filepath.Dir(path) == "." {
				rootFiles[path] = fileEntryInfo{
					absPath:     abs,
					sizeInBytes: info.Size(),
				}
			} else {
				for k, v := range rootDirs {
					currentParent := filepath.Dir(path)
					for {
						if currentParent == k {
							if cmd.verbose {
								fmt.Printf("%s is in %s\n", path, k)
							}
							v.sizeInBytes += info.Size()
							return nil
						}
						currentParent = filepath.Dir(currentParent)
						if currentParent == "." {
							break
						}
					}
				}
			}
		}
		return nil
	}

	displayFiles := func() {
		sorted := []struct {
			n string
			s int64
		}{}
		for k, v := range rootFiles {
			sorted = append(sorted, struct {
				n string
				s int64
			}{
				n: k,
				s: v.sizeInBytes,
			})
		}
		slices.SortFunc(sorted, func(a, b struct {
			n string
			s int64
		}) int {
			if a.s > b.s {
				return -1
			} else if a.s < b.s {
				return 1
			}
			return 0
		})
		fmt.Printf(outputHeaderTemplate, "File")
		for _, v := range sorted {
			fmt.Printf(outputLineTemplate, v.n, v.s)
		}
	}

	displayDirectories := func() {
		sorted := []struct {
			n string
			s int64
		}{}
		for k, v := range rootDirs {
			sorted = append(sorted, struct {
				n string
				s int64
			}{
				n: k,
				s: v.sizeInBytes,
			})
		}
		slices.SortFunc(sorted, func(a, b struct {
			n string
			s int64
		}) int {
			if a.s > b.s {
				return -1
			} else if a.s < b.s {
				return 1
			}
			return 0
		})
		fmt.Printf(outputHeaderTemplate, "Directory")
		for _, v := range sorted {
			fmt.Printf(outputLineTemplate, v.n, v.s)
		}
	}

	if cmd.startingDir == "" {
		return newRequiredArgumentMissingError(argNameStartingDir)
	}
	if startingDirAbs, err := filepath.Abs(cmd.startingDir); err != nil {
		return fmt.Errorf("error getting absolute path of starting directory: %w", err)
	} else {
		fmt.Printf("Starting directory: %s\n", startingDirAbs)
		fileSystem := os.DirFS(startingDirAbs)
		if err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if cmd.verbose {
					fmt.Printf("Error walking filesystem: %s\n", err)
				}
				return nil
			}
			if path == "." {
				return nil
			} else if d.IsDir() {
				return processDir(path)
			} else {
				return processFile(path, d)
			}
		}); err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}
		fmt.Println()
		displayFiles()
		fmt.Println()
		displayDirectories()
		return nil
	}
}

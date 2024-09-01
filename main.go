package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	fmt.Fprintln(out)
	if !printFiles {
		printCataloges(out, path, "")
		return nil
	}
	printCatalogesAndFiles(out, path, "")
	return nil
}

func printCataloges(out io.Writer, path string, prefics string) {
	dirs := findDirs(path)
	amountDirs := len(dirs)
	if amountDirs == 0 {
		return
	}

	for idx, el := range dirs {
		fName := el.Name()
		fPath := filepath.Join(path, fName)
		nextDirs := findDirs(fPath)
		amountDirsInCurrDir := len(nextDirs)
		var newPrefics string

		if idx < amountDirs-1 {
			fmt.Fprintf(out, "%s├───%s\n", prefics, fName)
			if amountDirsInCurrDir > 0 {
				newPrefics = prefics + "│" + "\t"
			}
		}

		if idx == amountDirs-1 {
			fmt.Fprintf(out, "%s└───%s\n", prefics, fName)
			if amountDirsInCurrDir > 0 {
				newPrefics = prefics + "\t"
			}
		}
		printCataloges(out, fPath, newPrefics)
	}
}

func findDirs(path string) []fs.DirEntry {
	filesList, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	var dirs []fs.DirEntry
	for _, fileEntry := range filesList {
		if !fileEntry.IsDir() {
			continue
		}
		dirs = append(dirs, fileEntry)
	}
	return dirs
}

func printCatalogesAndFiles(out io.Writer, path string, prefix string) {
	filesList, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	amountFiles := len(filesList)

	for idx, el := range filesList {
		fName := dirEntryName(el)

		//если элемент - файл
		if !el.IsDir() {
			if idx < amountFiles-1 {
				printNotLast(out, prefix, fName)
			}
			if idx == amountFiles-1 {
				printLast(out, prefix, fName)
			}
			continue
		}

		//если элемент - каталог
		fPathCurrEl := filepath.Join(path, fName)
		elemsCurrDir, err := os.ReadDir(fPathCurrEl)
		if err != nil {
			log.Fatal(err)
		}
		amountElemsCurrDir := len(elemsCurrDir)
		var newPrefix string

		// Если элемент не последний в каталоге
		if idx < amountFiles-1 {
			printNotLast(out, prefix, fName)
			if amountElemsCurrDir > 0 {
				newPrefix = prefix + "│\t"
			}
			if amountElemsCurrDir == 0 {
				newPrefix = prefix + "\t"
			}
		}

		// Если элемент последний в каталоге
		if idx == amountFiles-1 {
			printLast(out, prefix, fName)
			newPrefix = prefix + "\t"
		}

		// Рекурсивный вызов
		printCatalogesAndFiles(out, fPathCurrEl, newPrefix)

	}
}

func dirEntryName(el fs.DirEntry) string {
	if el.IsDir() {
		return el.Name()
	}

	info, err := el.Info()
	if err != nil {
		log.Fatal(err)
	}
	fSize := info.Size()
	fSizeStr := strconv.Itoa(int(fSize)) + "b"
	if fSize == 0 {
		fSizeStr = "empty"
	}
	return el.Name() + " (" + fSizeStr + ")"
}

func printNotLast(out io.Writer, prefix, fName string) {
	fmt.Fprintf(out, "%s├───%s\n", prefix, fName)
}

func printLast(out io.Writer, prefix, fName string) {
	fmt.Fprintf(out, "%s└───%s\n", prefix, fName)
}

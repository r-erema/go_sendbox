package tree

import (
	"bytes"
	"os"
	"sort"
	"strconv"
)

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {
	outputThree(out, path, printFiles, 0, false, "") //todo: parentEnded default value?
	return nil
}

func outputThree(out *bytes.Buffer, path string, printFiles bool, nestingLevel int, parentEnded bool, childPrefix string) {

	currFile, _ := os.Open(path)
	files, _ := currFile.Readdir(0)

	if !printFiles {
		files = Filter(files, func(info os.FileInfo) bool {
			return info.IsDir()
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var fileName string
	if nestingLevel != 0 {
		childPrefix = chp(childPrefix, parentEnded)
	}

	for i, file := range files {

		if !file.IsDir() && !printFiles {
			continue
		}

		lastElem := i+1 == len(files)
		fileName = resolveFileName(file, lastElem, file.IsDir(), true)
		str := childPrefix + fileName
		out.Write([]byte(str))
		if file.IsDir() {
			outputThree(out, currFile.Name()+string(os.PathSeparator)+file.Name(), printFiles, nestingLevel+1, lastElem, childPrefix)
		}
	}

}

func Filter(vs []os.FileInfo, f func(os.FileInfo) bool) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func chp(existedChildPrefix string, parentEnded bool) string {
	if parentEnded {
		existedChildPrefix += "\t"
	} else {
		existedChildPrefix += "│\t"
	}
	return existedChildPrefix
}

func resolveFileName(info os.FileInfo, isLast bool, isDir, needNewLine bool) string {
	var sizeStr, prefix string

	if !isDir {
		if info.Size() > 0 {
			sizeStr = " (" + strconv.FormatInt(info.Size(), 10) + "b)"
		} else {
			sizeStr = " (empty)"
		}
	}

	if isLast {
		prefix = "└───"
	} else {
		prefix = "├───"
	}

	str := prefix + info.Name() + sizeStr
	if needNewLine {
		str += "\n"
	}

	return str
}

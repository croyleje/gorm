package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
)

// TODO: rewrite standardize nameing conversions filename and path.
// TODO: rewrite standardize variable names camelCase and PascalCase.

var (
	trashDir     = GetTrashDir() + "files/"
	trashInfoDir = GetTrashDir() + "info/"
	dirSizesFile = GetTrashDir() + "directorysizes"
)

type Item struct {
	Name         string
	Type         string
	Size         string
	Path         string
	DeletionDate string
	IsChecked    bool
}

func GetEntries() ([]Item, error) {
	if !checkFileExists(trashDir) {
		GetTrashDir()
	}
	file, err := os.Open(trashDir)
	if err != nil {
		log.Fatalf("error: %v %v", file, err)
	}

	defer file.Close()

	var list []Item

	fileList, _ := file.Readdir(0)

	for _, file := range fileList {
		list = append(list, Item{
			Name:         file.Name(),
			Type:         boolToString(file),
			Size:         getSize(file),
			Path:         getFullPath(file.Name()),
			DeletionDate: getDeletionDate(file.Name()),
			IsChecked:    false,
		})
	}

	return list, err
}

func checkFileExists(path string) bool {
	_, error := os.Stat(path)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func boolToString(item fs.FileInfo) string {
	if item.IsDir() {
		return "Directory"
	} else {
		return "File"
	}
}

func getSize(file fs.FileInfo) string {
	if file.IsDir() {
		size := getDirSize(file.Name())
		return size
	} else {
		size := byteSizeToDecimal(file.Size())
		return size
	}
}

func getDirSize(file string) string {
	for _, line := range lines(dirSizesFile) {
		if strings.Contains(line, file) {
			size, _, _ := strings.Cut(line, " ")
			// size = strings.Trim(line, " ")
			num, _ := strconv.ParseInt(size, 10, 64)
			size = byteSizeToDecimal(num)
			return size
		} else {
			return "\b"
			// return ""
		}
	}
	return "error"
}

func byteSizeToDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		// return fmt.Sprintf("%d B", b)
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	// return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func getFullPath(file string) string {
	for _, line := range lines(trashInfoDir + file + ".trashinfo") {
		if strings.Contains(line, "Path=") {
			line = strings.TrimLeft(line, "Path=")
			return line
		}
	}
	return "path not found"
}

func getDeletionDate(file string) string {
	trashInfoFile := trashInfoDir + file + ".trashinfo"
	for _, line := range lines(trashInfoFile) {
		if strings.Contains(line, "DeletionDate=") {
			line = strings.TrimLeft(line, "DeletionDate=")
			return line
		}
	}
	return "deletionDate not found"
}

func lines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error: %v %v", file, err)
	}
	scanner := bufio.NewScanner(file)
	result := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

func CliTrashList() {
	file, err := os.Open(trashDir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}

	defer file.Close()

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Size"},
			{Align: simpletable.AlignCenter, Text: "IsDirectory"},
			{Align: simpletable.AlignCenter, Text: "LastModTime"},
		},
	}

	var cells [][]*simpletable.Cell

	fileList, _ := file.Readdir(0)

	for _, files := range fileList {
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: files.Name()},
			{Text: byteSizeToDecimal(files.Size())},
			{Text: strconv.FormatBool(files.IsDir())},
			{Text: files.ModTime().String()},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.SetStyle(simpletable.StyleUnicode)

	fmt.Println(table.String())

}

func TrashPut(path []string) {
	files := path
	for _, v := range files {
		if isDirectory(v) == true {
			v = strings.ReplaceAll(v, "/", "")
			generateDirCache(v)
		}
		if checkFileExists(trashDir+v) == true {
			newName := incName(v)
			os.Rename(v, trashDir+newName)
			generateTrashInfo(v, newName)
		} else {
			os.Rename(v, trashDir+v)
			generateTrashInfo(v)
		}
	}
}

// TODO: rewrite path call to use `os` pachaage.
// variadic function if only one arg is recived file has not be incremented
func generateTrashInfo(arg ...string) {
	header := "[Trash Info]\n"
	path, _ := exec.Command("readlink", "-f", arg[0]).Output()
	date := time.Now().Format("2006-01-02T15:04:05")

	var trashfile string

	if len(arg) == 1 {
		trashfile = trashInfoDir + arg[0] + ".trashinfo"
	} else {
		trashfile = trashInfoDir + arg[1] + ".trashinfo"
	}

	// TODO: rewrite file write calls.

	// f, _ := os.OpenFile(trashfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f, _ := os.OpenFile(trashfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	defer f.Close()

	w := bufio.NewWriter(f)

	w.WriteString(header)
	w.WriteString("Path=" + string(path))
	w.WriteString("DeletionDate=" + date)

	w.Flush()

}

func isDirectory(path string) bool {
	info, _ := os.Stat(path)
	return info.IsDir()
}

func dirSize(path string) (string, error) {
	out, err := exec.Command("du", "-b", path).Output()
	match, _ := regexp.Compile(`\d+`)
	find := string(match.Find(out))

	return find, err

}

func generateDirCache(dir string) {
	size, _ := dirSize(dir)
	epoch := time.Now().Unix()
	time := strconv.FormatInt(epoch, 10)

	// f, _ := os.OpenFile(trashfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f, _ := os.OpenFile(dirSizesFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(size + " " + time + " " + dir + "\n")
	w.Flush()

}

// n variable used for the incName function to store iterations of loop
var n int = 1

// TODO: rewrite change variable names
func incName(filename string) string {
	var filepath string
	for checkFileExists(trashDir+filename) == true {
		filepath = fmt.Sprintf("%s_%d", filename, n)
		if checkFileExists(trashDir+filepath) == true {
			incName(filepath)
			n++
		}
		if checkFileExists(trashDir+filepath) == false {
			break
		}
	}
	return filepath
}

// TODO: add error handling and file checks
func RestoreItem(i string) {
	err := os.Rename(trashDir+i, getFullPath(i))
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(trashInfoDir + i + ".trashinfo")
}

// TODO: add error handling and file checks
func RestoreItemLocal(i string) {
	wd, _ := os.Getwd()
	wd = wd + "/"
	err := os.Rename(trashDir+i, wd+i)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(trashInfoDir + i + ".trashinfo")
}

// TODO: add error handling and checks
func DeleteItem(i string) {
	if isDirectory(trashDir + i) {
		os.RemoveAll(trashDir + i)
		removeDirSizesEntry(i)
	} else {
		os.Remove(trashDir + i)
	}
	os.Remove(trashInfoDir + i + ".trashinfo")
}

func IsHome(path string) bool {
	home, _ := os.UserHomeDir()
	abs, _ := filepath.Abs(path)

	if strings.Contains(abs, home) {
		return true
	} else {
		return false
	}
}

// func GetTrashDir() string {
// 	path, err := os.Getwd()
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatalf("unable to open trash file does not exist: %s", getMount(path))
// 	} else {
// 		if IsHome(path) {
// 			return home + "/.local/share/Trash/"
// 		}
// 	}
// 	return ""

// }

// func GetTrashDir() string {
// 	uid := os.Geteuid()
// 	uuid := fmt.Sprintf("%d", uid)
// 	path, err := os.Getwd()
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatalf("unable to open trash file does not exist: %s", getMount(path))
// 	} else {
// 		if IsHome(path) {
// 			return home + "/.local/share/Trash/"
// 		} else {
// 			return getMount(path) + "/.Trash-" + uuid + "/"
// 		}
// 	}
// 	return ""

// }

func GetTrashDir() string {
	uid := os.Geteuid()
	uuid := fmt.Sprintf("%d", uid)
	path, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	if IsHome(path) {
		if checkFileExists(home + "/.local/share/Trash/") {
			return home + "/.local/share/Trash/"
		} else {
			createTrashDir(home + "/.local/share/Trash/")
			return home + "/.local/share/Trash/"
		}
	}

	if !IsHome(path) {
		if checkFileExists(getMount(path) + "/.Trash-" + uuid + "/") {
			return getMount(path) + "/.Trash-" + uuid + "/"
		} else {
			createTrashDir(getMount(path) + "/.Trash-" + uuid + "/")
			return getMount(path) + "/.Trash-" + uuid + "/"
		}
	}
	return ""
}

// TODO: rewrite check and/or create Trash directories per mount point.
func createTrashDir(path string) {
	os.Mkdir(path, 0700)
	os.Mkdir(path+"files/", 0700)
	os.Mkdir(path+"info/", 0700)
}

func getMount(path string) string {
	pointRaw, err := exec.Command("df", "-l", "--output=target", path).Output()
	if err != nil {
		log.Fatalf("unable to find mount point: %v", err)
	}

	point := string(pointRaw)
	point = strings.TrimPrefix(point, "Mounted on")
	point = strings.Trim(point, "\n")

	return point

}

func removeDirSizesEntry(pattern string) {
	// NOTE: path to file to match on
	fpath := dirSizesFile

	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// NOTE: pattern to match string matched lines are removed

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// NOTE: the not operator
		if !strings.Contains(scanner.Text(), pattern) {
			_, err := buf.Write(scanner.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(fpath, buf.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

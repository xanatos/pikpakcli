package list

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/52funny/pikpakcli/conf"
	"github.com/52funny/pikpakcli/internal/pikpak"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var long bool
var recursive bool
var human bool
var path string
var parentId string

var ListCmd = &cobra.Command{
	Use:   "ls",
	Short: `Get the directory information under the specified folder`,
	Run: func(cmd *cobra.Command, args []string) {
		p := pikpak.NewPikPak(conf.Config.Username, conf.Config.Password)
		err := p.Login()
		if err != nil {
			logrus.Errorln("Login Failed:", err)
		}
		handle(&p, args)
	},
}

func init() {
	ListCmd.Flags().BoolVarP(&human, "human", "H", false, "display human readable format")
	ListCmd.Flags().BoolVarP(&long, "long", "l", false, "display long format")
	ListCmd.Flags().BoolVarP(&recursive, "recursive", "R", false, "display recursively")
	ListCmd.Flags().StringVarP(&path, "path", "p", "/", "display the specified path")
	ListCmd.Flags().StringVarP(&parentId, "parent-id", "P", "", "display the specified parent id")
}

func handle(p *pikpak.PikPak, args []string) {
	var currentPath string = ""

	if recursive && parentId == "" {
		currentPath = filepath.Clean(path)

		if currentPath != string(filepath.Separator) {
			currentPath = currentPath + string(filepath.Separator)
		}
	}

	var err error
	if parentId == "" {
		parentId, err = p.GetPathFolderId(path)
		if err != nil {
			logrus.Errorln("get path folder id error:", err)
			return
		}
	}

	listFolder(p, parentId, currentPath)
}

func listFolder(p *pikpak.PikPak, parentId string, currentPath string) {
	files, err := p.GetFolderFileStatList(parentId)
	if err != nil {
		logrus.Errorln("get folder file stat list error:", err)
		return
	}

	for _, file := range files {
		if long {
			if human {
				display(3, currentPath, &file)
			} else {
				display(2, currentPath, &file)
			}
		} else {
			display(0, currentPath, &file)
		}

		if recursive && file.Kind == "drive#folder" {
			var currentPath2 string = currentPath + file.Name + string(filepath.Separator)
			listFolder(p, file.ID, currentPath2)
		}
	}
}

// lH
// mode 0: normal print
// mode 2: long format
// mode 3: long format and human readable

func display(mode int, currentPath string, file *pikpak.FileStat) {
	var fileName = currentPath + file.Name

	switch mode {
	case 0:
		if file.Kind == "drive#folder" {
			fmt.Printf("%-20s\n", color.GreenString(fileName))
		} else {
			fmt.Printf("%-20s\n", fileName)
		}
	case 2:
		if file.Kind == "drive#folder" {
			fmt.Printf("%-26s d %-6s %-14s %s\n", file.ID, file.Size, file.CreatedTime.Format("2006-01-02 15:04:05"), color.GreenString(fileName))
		} else {
			fmt.Printf("%-26s f %-6s %-14s %s\n", file.ID, file.Size, file.CreatedTime.Format("2006-01-02 15:04:05"), fileName)
		}
	case 3:
		if file.Kind == "drive#folder" {
			fmt.Printf("%-26s d %-6s %-14s %s\n", file.ID, displayStorage(file.Size), file.CreatedTime.Format("2006-01-02 15:04:05"), color.GreenString(fileName))
		} else {
			fmt.Printf("%-26s f %-6s %-14s %s\n", file.ID, displayStorage(file.Size), file.CreatedTime.Format("2006-01-02 15:04:05"), fileName)
		}
	}
}

func displayStorage(s string) string {
	size, _ := strconv.ParseUint(s, 10, 64)
	cnt := 0
	for size > 1024 {
		cnt += 1
		if cnt > 5 {
			break
		}
		size /= 1024
	}
	res := strconv.Itoa(int(size))
	switch cnt {
	case 0:
		res += "B"
	case 1:
		res += "KB"
	case 2:
		res += "MB"
	case 3:
		res += "GB"
	case 4:
		res += "TB"
	case 5:
		res += "PB"
	}
	return res
}

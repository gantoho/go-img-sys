package logic

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type imgDataType struct {
	total int
	data  []string
}

func OpFiles(context *gin.Context) {
	host := context.Request.Host
	dst := "./files"
	files, err := os.Open(dst) //open the directory to read files in the directory
	if err != nil {
		fmt.Println("error opening directory:", err) //print error if directory is not opened
		return
	}
	defer files.Close()
	fileInfos, err := files.Readdir(-1)
	if err != nil {
		fmt.Println("error reading directory:", err) //if directory is not read properly print error message
		return
	}
	imgData := &imgDataType{}
	imgData.total = len(fileInfos)
	for _, fileInfos := range fileInfos {
		imgData.data = append(imgData.data, host+"/files/"+fileInfos.Name())
	}
	context.JSON(http.StatusOK, map[string]any{"total": imgData.total, "data": imgData.data})
}

func Bgimg(context *gin.Context) {
	dst := "./files"
	files, err := os.Open(dst) //open the directory to read files in the directory
	if err != nil {
		fmt.Println("error opening directory:", err) //print error if directory is not opened
		return
	}
	defer files.Close()
	fileInfos, err := files.Readdir(-1)
	if err != nil {
		fmt.Println("error reading directory:", err) //if directory is not read properly print error message
		return
	}
	context.String(http.StatusOK, fileInfos[rand.Intn(len(fileInfos))].Name())
}

func Upload(context *gin.Context) {
	host := context.Request.Host
	form, _ := context.MultipartForm()
	files := form.File["files"]
	return_ := map[int]string{}
	for i, file := range files {
		dst := "./files/" + file.Filename
		err := context.SaveUploadedFile(file, dst)
		if err != nil {
			return
		}
		return_[i] = host + "/files/" + file.Filename
	}
	context.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("%d个文件，上传成功", len(files)), "urls": return_})
}

func BgimgNum(context *gin.Context) {
	numStr := context.Param("number")
	num, _ := strconv.Atoi(numStr)
	dst := "./files"
	files, err := os.Open(dst) //open the directory to read files in the directory
	if err != nil {
		fmt.Println("error opening directory:", err) //print error if directory is not opened
		return
	}
	defer files.Close()
	fileInfos, err := files.Readdir(-1)
	if err != nil {
		fmt.Println("error reading directory:", err) //if directory is not read properly print error message
		return
	}
	group := make([]string, 0)
	if num <= 0 {
		num = 1
	}
	if num >= 10 {
		num = 10
	}
	for i := 0; i < num; i++ {
		group = append(group, fileInfos[rand.Intn(len(fileInfos))].Name())
	}
	context.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": group})
}

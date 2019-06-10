package main

import (
	"net/http"
	"strings"
	"crypto/md5"
	"time"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"encoding/json"
)

func sayHello(w http.ResponseWriter,r *http.Request)  {
	w.Write([]byte("hello world!!!"))
}

func main() {
	fileHandle := http.FileServer(http.Dir("./video"))
	//注册上传文件
	http.HandleFunc("/api/upload",uploadHandle)
	//注册获取视频列表
	http.HandleFunc("/api/list",getFileListHandle)
	http.HandleFunc("/sayHello", sayHello)
	http.Handle("/video/",http.StripPrefix("/video/",fileHandle))
	http.ListenAndServe(":8001",nil)
}

//uploadHandle
func uploadHandle(w http.ResponseWriter,r *http.Request)  {
//	限制客户端上传文件大小为10M
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.Body = http.MaxBytesReader(w, r.Body, 15*1024*1024)
	err := r.ParseMultipartForm(15 * 1024 * 1024)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	file, header, _ := r.FormFile("uploadFile")
	defer file.Close()
	ret := strings.HasSuffix(header.Filename, ".mp4")
	if ret==false {
		http.Error(w,"不是MP4",http.StatusInternalServerError)
		return
	}
//	获取随机文件名称
	sum := md5.Sum([]byte(header.Filename + time.Now().String()))
	//	转成16进制
	md5Str := fmt.Sprintf("%x", sum)
	newFileName := md5Str + ".mp4"
//	创建磁盘文件
	create, err := os.Create("./video/" + newFileName)
	defer create.Close()
	if err!= nil	 {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
//	拷贝文件
	if _, e := io.Copy(create, file);e!= nil {
		http.Error(w,e.Error(),http.StatusInternalServerError)
		return
	}
	w.Write([]byte("上传成功!"))
	return
}
func getFileListHandle(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	files, _:= filepath.Glob("video/*")
	var ret []string
	for _, file := range files {
		ret = append(ret,"http://"+r.Host+"/video/"+filepath.Base(file))

	}
	jsonRet, _ := json.Marshal(ret)
	w.Write(jsonRet)
	return

}
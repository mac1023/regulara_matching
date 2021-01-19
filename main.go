package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {

	webUrl, err := ioutil.ReadFile("source.txt")
	if err != nil {
		fmt.Println("文件读取错误", err)
		return
	}
	urlStr := string(webUrl)

	//获取网页内容
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("抓取网页报错 error: ", err.Error())
		return
	}
	data, err :=ioutil.ReadAll(resp.Body)
	webContent := string(data)

	reg := regexp.MustCompile(`Learn/.*(.mp3|.pdf)`)
	str := webContent

	var res [][]string
	res = reg.FindAllStringSubmatch(str, -1)

	fileName := "download.txt"
	file, err := PathExists(fileName)
	if file {
		fmt.Printf("path %s 已经存在，准备写入\n", fileName)
	} else {
		dstFile, err := os.Create(fileName)
		defer dstFile.Close()
		if err != nil {
			fmt.Printf("创建失败![%v]\n", err)
			return
		} else {
			fmt.Printf("创建成功!\n")

		}
	}

	fmt.Printf("开始写入...\n")
	for index, value := range res {
		for _, val := range value {

			if val == ".pdf" || val == ".mp3" {
				continue
			}
			enEscapeUrl, _ := url.QueryUnescape(val)
			fmt.Println(enEscapeUrl)

			str2 := strings.Split(enEscapeUrl, "/")
			str3 := str2[:len(str2) - 1]
			str4 := Implode("/", str3)
			ok, _ := PathExists(str4)
			if ok {
				//fmt.Printf("path %s 已经存在\n", str4)
			} else {
				err := os.MkdirAll(str4, 0777)
				//defer dstFile.Close()
				if err != nil {
					fmt.Printf("创建失败![%v]\n", err)
					return
				} else {
					fmt.Printf("创建成功!\n")

				}
			}

			ok, _ = PathExists(enEscapeUrl)
			if ok {
				continue
			}
			downloadUrl := "https://pan.uvooc.com/" + val
			// Get the data
			resp, err := http.Get(downloadUrl)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			// 创建一个文件用于保存 enEscapeUrl
			out, err := os.Create(enEscapeUrl)
			if err != nil {
				panic(err)
			}
			defer out.Close()

			// 然后将响应流和文件流对接起来
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				panic(err)
			}

			downloadUrl1 := downloadUrl + "\n"
			if index == len(res)-1 {
				downloadUrl1 = downloadUrl + "\n\n"
			}
			appendToFile(fileName, downloadUrl1)
		}

	}
	fmt.Printf("写入完成,本次写入%d条...\n", len(res))
}

func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("文件句柄打开失败. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

/*
   判断文件或文件夹是否存在
   如果返回的错误为nil,说明文件或文件夹存在
   如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
   如果返回的错误为其它类型,则不确定是否在存在
*/
func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Implode(glue string, pieces []string) string {
	return strings.Join(pieces, glue)
}
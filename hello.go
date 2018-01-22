package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	// ptnURL          = regexp.MustCompile(`<a href="(.*)" class="photolst_photo"`)
	ptnURL          = regexp.MustCompile(`<img(.*)src="(.*)view/photo/m/public/(.*)"`)
	ptnContentRough = regexp.MustCompile(`<img(.*)src="(.*)view/photo/l/public/(.*)"`)
	ptnImgName      = regexp.MustCompile(`public/(.*)\.jpg`)
	ptnHTMLTag      = regexp.MustCompile(`(?s)</?.*?>`)
	ptnSpace        = regexp.MustCompile(`(^\s+)|( )`)
)

/*
<div class="photo_wrap">
                <a href="https://www.douban.com/photos/photo/2510951297/" class="photolst_photo" title="">

                    <img width="37" src="https://img1.doubanio.com/view/photo/m/public/p2510951297.jpg" />
                </a><br/>
                <div class="pl"></div>
                <div style="color:#999">
                        <a href="https://www.douban.com/photos/photo/2510951297/#comments">2回应</a>
                </div>
						</div>
						https://img3.doubanio.com/view/photo/l/public/p2510951284.jpg

						<a class="mainphoto" href="https://www.douban.com/photos/photo/2510951279/#image" title="点击查看下一张">

                    <img width="502" src="https://img3.doubanio.com/view/photo/l/public/p2510951284.jpg">
                    </a>
*/
func Get(url string) (content string, statusCode int) {
	resp, err1 := http.Get(url)

	if err1 != nil {
		statusCode = -100
		return
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	statusCode = resp.StatusCode
	content = string(data)
	// fmt.Print("content:\n", content, "\n")

	return
}

type IndexItem struct {
	url     string
	title   string
	content string
}

func findIndex(content string) (index []string, err error) {
	matches := ptnURL.FindAllStringSubmatch(content, -1)
	fmt.Println(len(matches))
	fmt.Printf("%q", matches)
	index = make([]string, len(matches))
	for i, item := range matches {
		fmt.Println(item)
		fmt.Println(i)

		// index[i] = IndexItem{item[1], ""}
		index[i] = item[2] + "view/photo/l/public/" + item[3]
	}
	return
}

func readContent(url string) (content IndexItem) {

	// matches := ptnURL.FindAllStringSubmatch(content, -1)

	match := ptnImgName.FindStringSubmatch(url)
	if match != nil {
		content.url = url
		content.title = match[1] + ".jpg"
		raw, statusCode := Get(url)
		if statusCode != 200 {
			fmt.Print("Fail to get the raw data from", url, "\n")
			return
		}
		ioutil.WriteFile("./img/"+content.title, []byte(raw), 0644)
		// content = match[1]
	} else {
		return
	}

	return
}

func main() {
	fmt.Println(`Get index ...`)
	s, statusCode := Get("https://www.douban.com/photos/album/1613893843/?start=36")
	if statusCode != 200 {
		fmt.Printf("statusCode：%d\n", statusCode)
		return
	}
	index, _ := findIndex(s)

	fmt.Println(`Get contents and write to file ...`)
	for _, item := range index {
		fmt.Printf("Get content from %s and write to file.\n", item)
		// fileName := fmt.Sprintf("%s.txt", item.title)
		content := readContent(item)
		fmt.Println(content.title)

		// fmt.Printf("Finish writing to %s.\n", fileName)
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

//定义新的数据类型
type Spider struct {
	url    string
	header map[string]string
}

//定义 Spider get的方法
func (keyword Spider) get_html_header() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", keyword.url, nil)
	if err != nil {
	}
	for key, value := range keyword.header {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	return string(body)

}
func parse() {
	header := map[string]string{
		"Host":                      "movie.douban.com",
		"Connection":                "keep-alive",
		"Cache-Control":             "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/jpeg,*/*;q=0.8",
		"Referer":                   "https://movie.douban.com/top250",
	}

	//创建excel文件
	f, err := os.Create("./movies.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//写入标题
	f.WriteString("<meta http-equiv='content-type' content='text/html; charset=utf-8'><table>")

	//循环每页解析并把结果写入excel
	for i := 0; i < 10; i++ {
		fmt.Println("正在抓取第" + strconv.Itoa(i) + "页......")
		url := "https://movie.douban.com/top250?start=" + strconv.Itoa(i*25) + "&filter="
		spider := &Spider{url, header}
		html := spider.get_html_header()

		//链接
		pattern1 := `a href="https://movie.douban.com/subject/(.*?)/" class`
		rp1 := regexp.MustCompile(pattern1)
		find_txt1 := rp1.FindAllStringSubmatch(html, -1)

		//评价人数
		pattern2 := `<span>(.*?)评价</span>`
		rp2 := regexp.MustCompile(pattern2)
		find_txt2 := rp2.FindAllStringSubmatch(html, -1)

		//评分
		pattern3 := `property="v:average">(.*?)</span>`
		rp3 := regexp.MustCompile(pattern3)
		find_txt3 := rp3.FindAllStringSubmatch(html, -1)

		//电影名称
		pattern4 := `img(.*?)alt="(.*?)" src="(.*?)"`
		rp4 := regexp.MustCompile(pattern4)
		find_txt4 := rp4.FindAllStringSubmatch(html, -1)

		//电影名称
		pattern5 := `(.*)&nbsp;/&nbsp;(.*)&nbsp;/&nbsp;(.*)`
		rp5 := regexp.MustCompile(pattern5)
		find_txt5 := rp5.FindAllStringSubmatch(html, -1)

		// fmt.Println(find_txt1)
		// fmt.Println(find_txt2)
		// fmt.Println(find_txt3)
		// fmt.Println(find_txt4)
		// fmt.Println(find_txt5)
		// fmt.Println(find_txt6)

		//  打印全部数据和写入excel文件
		for i := 0; i < len(find_txt2); i++ {

			rp6 := regexp.MustCompile(`public/(.*)`)
			find_txt6 := rp6.FindAllStringSubmatch(find_txt4[i][3], -1)
			if _, err := os.Stat("./img/" + find_txt6[0][1]); err != nil {
				if os.IsNotExist(err) {
					// file does not exist
					url = find_txt4[i][3]
					spider = &Spider{url, header}
					raw := spider.get_html_header()
					ioutil.WriteFile("./img/"+find_txt6[0][1], []byte(raw), 0644)
				}
			}

			fmt.Printf("%s %s %s %s\n", find_txt6[0][1], find_txt4[i][2], find_txt3[i][1], find_txt2[i][1])
			f.WriteString("<tr><td rowspan='4'><img width=200 src='./img/" + find_txt6[0][1] + "'></td></tr><tr><td><a target='_blank' href='https://movie.douban.com/subject/" + find_txt1[i][1] + "/'>" + find_txt4[i][2] + "</a></td></tr><tr><td>" + find_txt3[i][1] + " " + find_txt2[i][1] + "评价</td></tr><tr><td>" + find_txt5[i][0] + "</td></tr>")

		}
	}
	f.WriteString("</table>")
}

func main() {
	t1 := time.Now() // get current time
	parse()
	elapsed := time.Since(t1)
	fmt.Println("爬虫结束,总共耗时: ", elapsed)

}

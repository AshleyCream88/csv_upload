package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const CheckImageURL = "https://aembed.com/adt8/chk_num.php"

func main() {
	records := openCsvGetRecords()
	cookies := loginGetCookie()
	u, err := url.Parse("https://aembed.com/adt8/video_get.php")
	checkError(err)
	for i := len(records) - 1; i >= 0; i-- {
		createVideo(struct {
			record []string
			url    *url.URL
			cookie []*http.Cookie
		}{record: records[i], url: u, cookie: cookies})
	}

}

func createVideo(opts struct {
	record []string
	url    *url.URL
	cookie []*http.Cookie
}) {
	query := opts.url.Query()
	query.Set("origurl", opts.record[3])
	opts.url.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", opts.url.String(), nil)
	req = reqAddAllCookie(req, opts.cookie)
	checkError(err)
	resp, err := http.DefaultClient.Do(req)
	checkError(err)
	defer resp.Body.Close()
	checkError(err)
	form := make(url.Values)
	form.Set("show", "1")
	form.Set("origurl", opts.record[3])
	form.Set("type", opts.record[4])
	form.Set("origtype", "")
	form.Set("capcode", "")
	form.Set("type2", "")
	form.Set("origtype2", "")
	form.Set("type3", "")
	form.Set("origtype3", "")
	form.Set("downloadbatch", "")
	links := strings.Split(strings.TrimSpace(opts.record[5]), "\n")
	for _, link := range links {
		form.Add("downloadurl", link)
	}
	form.Set("title", opts.record[1])
	form.Set("poster", opts.record[6])
	form.Set("contents", opts.record[7])
	form.Set("tags", opts.record[2])
	form.Set("origlength", "")
	form.Set("newlength", "")
	form.Set("starts", opts.record[10])
	form.Set("ends", opts.record[11])
	form.Set("act", "確認送出")
	form.Set("preview", opts.record[8])
	form.Set("ID", "")
	form.Set("MM_update", "form1")
	req, err = http.NewRequest("POST", "https://aembed.com/adt8/video_edit.php?order=&page=1&rows_pre_page=50", strings.NewReader(form.Encode()))
	checkError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	reqAddAllCookie(req, opts.cookie)
	resp, err = http.DefaultClient.Do(req)
	checkError(err)
	defer resp.Body.Close()
}

func reqAddAllCookie(req *http.Request, cookie []*http.Cookie) *http.Request {
	for _, cookie := range cookie {
		req.AddCookie(cookie)
	}
	return req
}

func openCsvGetRecords() [][]string {
	file, err := os.Open("AV_1.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	return records
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func loginGetCookie() []*http.Cookie {
	resp := geRespAndSaveImage()
	form := make(url.Values)
	form.Set("acct", "admin")
	form.Set("passwd", "sP61e2K8pTcOSI5J")
	form.Set("chk", inputNumber())
	req, err := http.NewRequest("POST", "https://aembed.com/adt8/index.php", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp.Cookies()
}

func inputNumber() string {
	var num string
	fmt.Println("請打開圖片並輸入check num")
	_, err := fmt.Scan(&num)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(num)
	return num
}

func geRespAndSaveImage() *http.Response {
	resp, err := http.Get(CheckImageURL)
	if err != nil {
		log.Fatalln(err)
	}
	//拿到image and cookie
	imgBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.OpenFile("chk_num.gif", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	writer := bufio.NewWriter(file)
	if _, err = writer.Write(imgBytes); err != nil {
		log.Fatalln(err)
	}
	if err = writer.Flush(); err != nil {
		log.Fatalln(err)
	}
	return resp
}

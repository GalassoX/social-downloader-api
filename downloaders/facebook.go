package downloaders

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type FbVideo struct {
	Url       string `json:"url"`
	SD        string `json:"sd"`
	HD        string `json:"hd"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
}

func Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("authority", "www.facebook.com")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("accept-language", "en-GB,en;q=0.9,tr-TR;q=0.8,tr;q=0.7,en-US;q=0.6")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"89\", \"Chromium\";v=\"89\", \";Not A Brand\";v=\"99\"")
	req.Header.Add("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	req.Header.Add("accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("cookie",
		"sb=Rn8BYQvCEb2fpMQZjsd6L382; datr=Rn8BYbyhXgw9RlOvmsosmVNT; c_user=100003164630629; _fbp=fb.1.1629876126997.444699739; wd=1920x939; spin=r.1004812505_b.trunk_t.1638730393_s.1_v.2_; xs=28%3A8ROnP0aeVF8XcQ%3A2%3A1627488145%3A-1%3A4916%3A%3AAcWIuSjPy2mlTPuZAeA2wWzHzEDuumXI89jH8a_QIV8; fr=0jQw7hcrFdas2ZeyT.AWVpRNl_4noCEs_hb8kaZahs-jA.BhrQqa.3E.AAA.0.0.BhrQqa.AWUu879ZtCw")

	trt := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: trt}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetFBVideo(videoURL string) (*FbVideo, []byte) {
	byteBody, err := Get(videoURL)
	if err != nil {
		fmt.Println("Error when do request: ", err)
		return nil, byteBody
	}

	r := regexp.MustCompile(`"playable_url":"(.*?)"`)
	result := r.FindStringSubmatch(string(byteBody))
	if len(result) < 2 {
		return nil, nil
	}

	sdSplit := strings.Split(result[1], "\\")
	sd := strings.Join(sdSplit, "")

	r = regexp.MustCompile(`"playable_url_quality_hd:"(.*?)"`)
	result = r.FindStringSubmatch(string(byteBody))

	hd := ""
	if len(result) >= 2 {
		hdSplit := strings.Split(result[1], "\\")
		hd = strings.Join(hdSplit, "")
	}

	r = regexp.MustCompile(`"preferred_thumbnail":{"image":{"uri":"(.*?)"`)
	result = r.FindStringSubmatch(string(byteBody))

	thumb := ""
	if len(result) >= 2 {
		thumbSplit := strings.Split(result[1], "\\")
		thumb = strings.Join(thumbSplit, "")
	}

	return &FbVideo{
		SD:        sd,
		HD:        hd,
		Url:       videoURL,
		Thumbnail: thumb,
	}, byteBody
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	url := "https://movie.douban.com/subject/30193669/"

	timeout := time.Duration(5 * time.Second)

	client := &http.Client{
		Timeout:timeout,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http read error")
		return
	}

	bodyString := string(body)

	re, _ := regexp.Compile("background-image: url(.*?)\"")
	divStyles := re.FindAllString(bodyString, -1)

	for i, divStyle := range divStyles {
		re, _ := regexp.Compile("https://(.*?).jpg")
		imageURL := re.FindString(divStyle)

		name := "actor" + strconv.Itoa(i) + ".jpg"
		fmt.Printf("begin download:%s\n",imageURL)

		os.Mkdir("image", os.ModePerm)
		resp, _ := http.Get(imageURL)
		body, _ := ioutil.ReadAll(resp.Body)
		out, _ := os.Create("image/" + name)
		io.Copy(out, bytes.NewReader(body))
	}
}

package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	var (
		inputUrl   string
		outputFile string
	)
	flag.StringVar(&inputUrl, "url", "", "url path to m3u8")
	flag.StringVar(&outputFile, "file", "", "path to output file")

	flag.Parse()
	if inputUrl == "" {
		log.Fatal("url is not required")
	}

	if outputFile == "" {
		log.Fatal("file is required")
	}

	// received file from server
	client := http.Client{}
	req, err := http.NewRequest("GET", inputUrl, nil)
	
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "cookie_name",Value:"value",Expires:expiration, Path: "/"}
	req.AddCookie(&cookie)
	cookie2 := http.Cookie{Name: "cookie_name2",Value:"value2",Expires:expiration, Path: "/"}
	req.AddCookie(&cookie2)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Download error: ", err)
	}
	defer resp.Body.Close()
	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	//log.Printf(string(bodyBytes))

	// create output file
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal("Download error: ", err)
	}
	defer f.Close()

	// read server response line by line
	scanner := bufio.NewScanner(resp.Body)
	i := 0
	for scanner.Scan() {
		l := scanner.Text()

		log.Printf(l)
		// if line contains url address
		if strings.HasPrefix(l, "some prefix") {
			// download file part
			part, err := downloadFilePart("base url" + l)
			if err != nil {
				log.Fatal("Download part error: ", err)
			}

			// write part to output file
			if _, err = f.Write(part); err != nil {
				log.Fatal("Write part to output file: ", err)
			}
			log.Printf("Download part %d\n", i)
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

//downloadFilePart download file part from server
func downloadFilePart(url string) ([]byte, error) {
	log.Printf(string(url))
	result := make([]byte, 0)

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "cookie_name",Value:"value",Expires:expiration, Path: "/"}
	req.AddCookie(&cookie)

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if result, err = ioutil.ReadAll(resp.Body); err != nil {
		return result, err
	}

	return result, err
}

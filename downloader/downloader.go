package downloader

import (
	"fmt"
	"github.com/rjsamson/rss"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const feedUrl = "https://elixirsips.dpdcart.com/feed"
const episodeDirectory = "episodes"

var wg sync.WaitGroup

var config struct {
	Username string
	Password string
	Episodes int
}

func Download(username string, password string, episodes int) {
	config.Username = username
	config.Password = password
	config.Episodes = episodes

	os.Mkdir(episodeDirectory, os.ModeDir|0777)

	client := &http.Client{}
	req, err := http.NewRequest("GET", feedUrl, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
		return
	}

	xmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	feed, err := rss.Parse(xmlData)

	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, item := range feed.Items[0:episodes] {
		wg.Add(1)
		go downloadFile(item.Enclosure.Url)
	}

	wg.Wait()
}

func downloadFile(url string) {
	defer wg.Done()
	fmt.Println(url)
	filename := filenameFromUrl(url)
	fmt.Println("Downloading", filename)

	out, err := os.Create(filename)
	defer out.Close()

	if err != nil {
		log.Fatalln(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(config.Username, config.Password)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
		return
	}

	_, err = io.Copy(out, resp.Body)

	fmt.Println("Finished", filename)
}

func filenameFromUrl(url string) string {
	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1]
	return fmt.Sprintf("%s/%s", episodeDirectory, filename)
}

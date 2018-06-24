package main

import (
	"log"
	"regexp"
	"time"
	"os"
	"os/exec"
	"os/user"
	"path"
	"math/rand"
	"io/ioutil"
	"net/http"
	"encoding/xml"
)

func main() {
	url := "https://www.flickr.com/services/feeds/photos_public.gne?id=130608600@N05&lang=en-us&format=atom"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Feed struct {
		XMLName xml.Name `xml:"feed"`
		Entries []string `xml:"entry>content"`
	}

	feed := Feed{}

	err = xml.Unmarshal([]byte(body), &feed)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().Unix())
	randomized := feed.Entries[rand.Int() % len(feed.Entries)]

	re := regexp.MustCompile("img src=\"https:\\/\\/farm\\d+\\.staticflickr\\.com\\/\\d+\\/(\\d+)_\\w+_m\\.jpg\"")
	img_url := re.FindStringSubmatch(randomized)
	if img_url == nil {
		log.Fatal("Error parsing content: ", randomized)
	}

	full_url := "https://www.flickr.com/photos/spacex/" + img_url[1] + "/sizes/o/"

	resp, err = http.Get(full_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	re2 := regexp.MustCompile("<div id=\"allsizes\\-photo\">\\s+<img src=\"(https:\\/\\/.+staticflickr.com/(\\d+/)*(.+_o.jpg))\">\\s+</div>")
	orig_url := re2.FindSubmatch(body)

	file_url := string(orig_url[1])
	file_name := "spacex_" + string(orig_url[len(orig_url) - 1])

	usr, _ := user.Current()
	dirpath := path.Join(usr.HomeDir, ".spacex_wallpapers")
	fullpath := path.Join(dirpath, file_name)
	os.MkdirAll(dirpath, os.FileMode(0755))

	err = exec.Command("wget", "-O", fullpath, file_url).Run()
	if err != nil {
		log.Fatal(err)
	}



	err = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://" + fullpath).Run()
	if err != nil {
		log.Fatal(err)
	}

}

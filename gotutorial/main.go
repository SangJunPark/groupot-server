package main

import (
	"fmt"
	"gotutorial/cli"
	"gotutorial/db"
	"time"
)

// import (
// 	"gotutorial/blockchain"
// 	"gotutorial/explorer"
// )

// func main() {
// 	blockchain := blockchain.GetBlockchain()
// 	blockchain.AddBlock("first")
// 	blockchain.AddBlock("second")
// 	explorer.Start()
// }

func main() {
	defer db.Close()
	cli.Start()
	//wallet.Wallet()
}

type hiChan chan int

//channel that receive only, read only
func receive(c <-chan int) {
	for {
		time.Sleep(1 * time.Second)

		a, ok := <-c
		fmt.Println("receive : ", a)
		if !ok {
			break
		}
	}
}

//channel that write only
func count(mainChain chan<- int, length int) {
	for i, _ := range [10]int{} {
		fmt.Println("sending : ", i)
		mainChain <- i
		fmt.Println("sent : ", i)

	}

	close(mainChain)
}

//send, receive operation block

// func string2Byte(str string) []byte {
// 	ret := [1]int{1}
// 	for _, b := range str{
// 		ret = append(ret, b)
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// )

// var baseUrl = "https://kr.indeed.com/jobs?q=python&limit=50"

// type extractedJob struct {
// 	id       string
// 	title    string
// 	location string
// }

// func main() {
// 	pages := getPages()

// 	//mainChannel := make(chan map[string]extractedJob)
// 	fmt.Println(pages)

// 	for i := 0; i < pages; i++ {
// 		url := getPage(i)
// 		fmt.Println(url)
// 	}
// }

// func extractData(c chan<- map[string]extractedJob) {
// 	t := extractedJob{
// 		id: "1",
// 		title: "1",
// 		location: "2",
// 	}

// }

// func cleanString(str string) string {
// 	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
// }

// func getPage(page int) string {
// 	pageUrl := baseUrl + "&start=" + strconv.Itoa(page*50)

// 	res, err := http.Get(pageUrl)
// 	checkError(err)
// 	checkResponse(res)

// 	defer res.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(res.Body)

// 	checkError(err)
// 	doc.Find(".tapItem").Each(func(i int, s *goquery.Selection) {
// 		link, _ := s.Attr("href")
// 		title := cleanString(s.Find(".jobTitle>span").Text())
// 		location := cleanString(s.Find(".companyLocation").Text())
// 		summary := cleanString(s.Find(".job-snippet").Text())
// 		fmt.Println(link, title, location, summary)

// 	})
// 	//fmt.Println(doc.Find(".job-seen-beacon").Html())
// 	// doc.Find(".mosaic-provider-jobcards").Each(func(i int, s *goquery.Selection) {
// 	// 	// (s.Find("a").Each(func(i int, ss *goquery.Selection) {
// 	// 	// 	fmt.Println(ss.Html())
// 	// 	// 	pages
// 	// 	// }))
// 	// 	fmt.Print("dfdf")
// 	// 	fmt.Println(s.Html())
// 	// })

// 	return pageUrl
// }

// func getPages() int {
// 	pages := 0
// 	res, err := http.Get(baseUrl)
// 	checkError(err)
// 	checkResponse(res)

// 	defer res.Body.Close()
// 	doc, err := goquery.NewDocumentFromReader(res.Body)

// 	checkError(err)
// 	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
// 		// (s.Find("a").Each(func(i int, ss *goquery.Selection) {
// 		// 	fmt.Println(ss.Html())
// 		// 	pages
// 		// }))
// 		pages = s.Find("a").Length()
// 	})

// 	return pages
// }

// func checkError(err error) {
// 	if err != nil {
// 		log.Fatal("err")
// 	}
// }

// func checkResponse(res *http.Response) {
// 	if res.StatusCode >= 400 {
// 		log.Fatal("status code err")
// 	}
// }

// func (a extractedJob)checkResponse(res *http.Response){

// }

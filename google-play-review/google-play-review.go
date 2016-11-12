package main

import (
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"strings"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	
	getReviewUrl := "https://play.google.com/store/getreviews"

	token := "iJRl6jtoMtdyH525lsTu2yBOOOc:1478943963364"

	file, _ := os.Create("final.csv")
	defer file.Close()
	
	reviewFound := 0
	for pageNum := 0; pageNum < 1500; pageNum++ {

	    requestBody := []byte("reviewType=0&pageNum=1&id=blibli.mobile.commerce&reviewSortOrder=4&xhr=1&hl=in&token=" + token)
	 	getReviewUrlResponse, err := http.Post(getReviewUrl, "application/x-www-form-urlencoded;charset=UTF-8", bytes.NewReader(requestBody))
	 	if err != nil {
	        fmt.Println(err)
	    }
	    defer getReviewUrlResponse.Body.Close()
	    
	    jsonBodyByte, err := ioutil.ReadAll(getReviewUrlResponse.Body)
	    if err != nil {
	        fmt.Println(err)
	    }
	 	jsonBody := string(jsonBodyByte)
	 	jsonBody = jsonBody[17:len(jsonBody)]
	 	jsonBody = jsonBody[0:len(jsonBody)-6]
	 	htmlBody, _ := strconv.Unquote(`"` + jsonBody + `"`)
	 	

	 	//read from output.txt
	 	// buf, _ := ioutil.ReadFile("output.txt")
	 	// htmlBody := string(buf) 
	 	//parsing html

	 	htmlBodyReader := bytes.NewReader([]byte(htmlBody))
	 	doc, _ := goquery.NewDocumentFromReader(htmlBodyReader) 	
	 	bodyItemSelection := doc.Find("body").Children()
	 	size := bodyItemSelection.Size()
	 	if(size > 0){
		 	var bodyItems = make([]*goquery.Selection, size, size)
		 	bodyItems[0] = bodyItemSelection.First()
		 	for i := 1; i < size; i++ {
		 		bodyItemSelection = bodyItemSelection.Next()
		 		bodyItems[i] = bodyItemSelection.First()

		 	}

		 	//write to file
		 	
		 	nextClassName := ""
		 	for i, bodyItem := range bodyItems {
		 		className, _ := bodyItem.Attr("class")
		 		nextIndex := i + 1
		 		if i < size - 1 {
		 			nextClassName, _ = bodyItems[nextIndex].Attr("class")
		 		} else {
		 			nextClassName = ""
		 		}
		 		if className == "single-review" && nextClassName == "developer-reply" {
		 			reviewFound++
		 			line := parseReview(bodyItem) + parseDeveloperReply(bodyItems[nextIndex])
		 			file.WriteString(line + "\n")
		 			file.Sync()
			 		fmt.Println(strconv.Itoa(reviewFound) + ". " + line)
		 		} else if className == "single-review" && nextClassName != "developer-reply"{
		 			reviewFound++
		 			line := parseReview(bodyItem) + "||"
		 			file.WriteString(line + "\n")
		 			file.Sync()
			 		fmt.Println(strconv.Itoa(reviewFound) + ". " + line)
		 		} 
		 	}
	 	}
 	}
}

func parseReview(review *goquery.Selection) string {
	authorName := strings.TrimSpace(review.Find(".author-name").Text())
	line := authorName
	reviewDate := strings.TrimSpace(review.Find(".review-info .review-date").Text())
	line = line + "|" + reviewDate
	reviewTitle := stripNonLatinCharacter(strings.TrimSpace(review.Find(".review-title").Text()))
	line = line + "|" + reviewTitle
	reviewBody := stripNonLatinCharacter(strings.TrimSpace(review.Find(".review-body").Text()))
	reviewBody = strings.TrimSpace(reviewBody[len(reviewTitle): len(reviewBody) - 14])
	line = line + "|" + reviewBody
	return line
}

func stripNonLatinCharacter(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}

func parseDeveloperReply(developerReply *goquery.Selection) string {
 		developerReplyDate := strings.TrimSpace(developerReply.Find(".review-date").Text())
 		line := "|" + developerReplyDate
 		developerReplyContent := stripNonLatinCharacter(strings.TrimSpace(developerReply.Text()))
 		developerReplyContent = strings.TrimSpace(developerReplyContent[len(developerReplyDate): len(developerReplyContent)])
 		line = line + "|" + developerReplyContent
 		return line
}
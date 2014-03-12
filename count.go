package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// http://api.stackoverflow.com/1.1/search?tagged=neo4j&page=1&pagesize=100
type SearchRequest struct {
	Tagged   string
	Page     int
	PageSize int
}

type SearchResponse struct {
	Total     int        `json="total"`
	Page      int        `json="page"`
	PageSize  int        `json="pagesize"`
	Questions []Question `json="questions"`
}

type Question struct {
	Tags         []string `json="tags"`
	AnswerCount  int      `json="answer_count"`
	QuestionId   int      `json="question_id"`
	CreationDate int      `json="creation_date"`
	Title        string   `json="title"`
	Body         string   `json="body"`
}

// http://api.stackoverflow.com/1.1/questions/22323946?body=true
type QuestionRequest struct {
	Ids  string
	Body bool
}

func (qr QuestionRequest) MakeRequest() (*SearchResponse, error) {
	url := fmt.Sprintf("http://api.stackoverflow.com/1.1/questions?ids=%s&body=%v", qr.Ids, qr.Body)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	qresp := SearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qresp)
	if err != nil {
		return nil, err
	}
	return &qresp, nil
}

func (sr SearchRequest) MakeRequest() (*SearchResponse, error) {
	url := fmt.Sprintf("http://api.stackoverflow.com/1.1/search?tagged=%s&page=%d&pagesize=%d", sr.Tagged, sr.Page, sr.PageSize)
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	sresp := SearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&sresp)
	if err != nil {
		return nil, err
	}
	return &sresp, nil
}

func main() {
	sr := SearchRequest{
		Tagged:   "neo4j",
		Page:     1,
		PageSize: 100,
	}
	for {
		resp, err := sr.MakeRequest()
		if err != nil {
			fmt.Println(err)
		}

		if len(resp.Questions) == 0 {
			break
		}
		if sr.Page != resp.Page {
			break
		}
		ids := ""
		for _, q := range resp.Questions {
			ids = ids + strconv.Itoa(q.QuestionId) + ";"
		}
		ids = ids[:len(ids)-1]
		qr := QuestionRequest{Ids: ids, Body: true}
		qresp, err := qr.MakeRequest()
		if err != nil {
			fmt.Println(err)
			break
		}

		for _, q := range qresp.Questions {
			buf, err := json.Marshal(q.Tags)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Sprintf("{question_id:%d, tags:%s, body_length: %d}\n", q.QuestionId, string(buf), len(q.Body))
		}
		sr.Page++
	}
}

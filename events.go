package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"comail.io/go/colog"
)

//Website is a shorthand for the map[string]interface{}
type Website map[string]interface{}

type Author struct {
	Id                  string `json:"id"`
	LastLoginOn         string `json:"lastLoginOn"`
	LastActiveOn        string `json:"lastActiveOn"`
	IsDeactivated       bool   `json:"isDeactivated"`
	Deleted             bool   `json:"deleted"`
	DisplayName         string `json:"displayName"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	EmailVerified       bool   `json:"emailVerified"`
	Bio                 string `json:"bio"`
	RevalidateTimestamp string `json:"revalidateTimestamp"`
	SystemGenerated     bool   `json:"systemGenerated"`
}

type StructuredContent struct {
	Type      string `json:"_type"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type Items struct {
}

type Event struct {
	Id                string            `json:"id"`
	CollectionId      string            `json:"collectionId"`
	RecordType        string            `json:"recordType"`
	AddedOn           string            `json:"addedOn"`
	UpdatedOn         string            `json:"updatedOn"`
	PublishOn         string            `json:"publishOn"`
	AuthorId          string            `json:"authorId"`
	UrlId             string            `json:"urlId"`
	Title             string            `json:"title"`
	SourceUrl         string            `json:"sourceUrl"`
	Body              string            `json:"body"`
	Author            Author            `json:"author"`
	FullUrl           string            `json:"fullUrl"`
	AssetUrl          string            `json:"assetUrl"`
	ContentType       string            `json:"contentType"`
	StructuredContent StructuredContent `json:"structuredContent"`
	StartDate         string            `json:"startDate"`
	EndDate           string            `json:"endDate"`
	Items             []Items           `json:"items"`
}

func main() {
	colog.Register()
	colog.SetFlags(log.Lshortfile)

	fileName := "./testevents.txt"

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("err: could not read file (%v)", err)
		return
	}

	// var dat map[string]interface{}
	w := Event{}
	if err := json.Unmarshal(b, &w); err != nil {
		log.Printf("err: could not unmarshal file (%v)", err)
		return
	}

	log.Printf("ok: unmarshalled file")

	log.Printf("debug: %v", w)
	// log.Printf("debug: %v", w["website"])
	// log.Printf("debug: collection info: %v", w["collection"])
	// log.Printf("debug: upcoming: %v", w["upcoming"])
	//
	// upcoming := Event{}
	// if err := json.Unmarshal(w["upcoming"].([]byte), &upcoming); err != nil {
	// 	log.Printf("error: could not unmarshal upcoming events")
	// }
}

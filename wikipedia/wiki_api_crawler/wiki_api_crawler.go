// https://jdanger.com/build-a-web-crawler-in-go.html
package main

import (
	"GIG-SDK/models"
	"GIG-SDK/request_handlers"
	"GIG-Scripts/wikipedia/wiki_api_crawler/decoders"
	"GIG-Scripts/wikipedia/wiki_api_crawler/requests"
	"flag"
	"log"
	"os"
	"sync"
)

var visited = make(map[string]bool)

func main() {
	flag.Parse()
	args := flag.Args()
	log.Println(args)
	if len(args) < 1 {
		log.Println("starting title not specified")
		os.Exit(1)
	}
	queue := make(chan string)
	go func() { queue <- args[0] }()

	for title := range queue {
		if title != "" {
			entity := enqueue(title, queue)
			if !entity.IsNil() {
				_, err := request_handlers.CreateEntity(entity)
				if err != nil {
					log.Println(err.Error(), title)
				}
			}
		}
	}
	log.Println("end")
}

func enqueue(title string, queue chan string) models.Entity {
	log.Println("fetching", title)
	visited[title] = true
	entity := models.Entity{}

	var requestWorkGroup sync.WaitGroup
	for _, propType := range requests.PropTypes() {

		requestWorkGroup.Add(1)
		go func(prop string) {
			defer requestWorkGroup.Done()
			result, err := requests.GetContent(prop, title)
			if err != nil {
				log.Println(err)
			} else {
				decoders.Decode(result, &entity)
			}
		}(propType)
	}
	requestWorkGroup.Wait()

	if !entity.IsNil() {

		var (
			linkEntities []models.Entity
			err          error
		)

		for _, link := range entity.Links {
			if link.GetTitle() != "" {
				if !visited[link.GetTitle()] {
					//log.Println("	passed link ->", link.Title)
					go func(title string) {
						queue <- title
						//log.Println("	queued link ->", link.Title)
					}(link.GetTitle())
				}
				//add link as an entity
				linkEntities = append(linkEntities, models.Entity{Title: link.GetTitle()})
			}
		}

		entity, err = request_handlers.AddEntitiesAsLinks(entity, linkEntities)

		if err != nil {
			log.Println("error creating links:", err)
		}
	}
	return entity
}

package helpers

import (
	"GIG-Scripts/wikipedia/wiki_web_crawler/constants"
	"bufio"
	"errors"
	"log"
	"os"
)

func LoadQueueFromLog(queue chan string) error {
	files, err := getAllLogs(constants.QueueLogDir)
	//if no log files exist
	if err != nil {
		return err
	}
	if len(files) == 1 {
		return errors.New("no log files found")
	}

	lastLog := files[len(files)-1]

	//open log file
	lastLogFile, err := os.Open(lastLog)

	if err != nil {
		return err
	}
	logScanner := bufio.NewScanner(lastLogFile)
	logScanner.Split(bufio.ScanLines)

	for logScanner.Scan() {
		go func(url string) { queue <- url }(logScanner.Text())
	}
	err = lastLogFile.Close()
	if err != nil {
		return err
	}
	log.Println("queue initialized from log")
	return nil
}

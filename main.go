package main

import (
	"bufio"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

type Track struct {
	ID        int
	Performer string
	Title     string
	StartTime string
}

type SetList struct {
	FileName   string
	TrackCount int
	Tracks     []Track
}

func ConvertHMSToSeconds(time string) int {
	parts := strings.Split(time, ":")
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatal(err)
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Fatal(err)
	}
	return hours*3600 + minutes*60 + seconds
}

func ExportSetList(setList SetList, outputFile *os.File) {
	for _, track := range setList.Tracks {
		timeStamp := ConvertHMSToSeconds(track.StartTime)
		tsFloat := float64(timeStamp)
		fmt.Fprintf(outputFile, "%.6f\t%.6f\t%s\n", tsFloat, tsFloat, track.Performer+" - "+track.Title)
	}
}

func ConvertFile(i *multipart.File, o *os.File) {
	// read the file line by line using scanner
	scanner := bufio.NewScanner(*i)
	SetList := new(SetList)
	newTrack := new(Track)

	for scanner.Scan() {
		// do something with a line
		// log.Printf("line: %s\n", scanner.Text())
		line := scanner.Text()
		if strings.HasPrefix(line, "FILE ") && strings.HasSuffix(line, " WAVE") {
			// get the file name
			fileName := strings.TrimSpace(line[5 : len(line)-5])
			// log.Println("file name:", fileName)
			SetList.FileName = fileName
		}

		if strings.HasPrefix(line, "\tTRACK") && strings.HasSuffix(line, " AUDIO") {
			// log.Println("Track found")
			// split the line into an array of strings
			// using the space as the separator
			parts := strings.Split(line, " ")
			// get the number of tracks
			trackID := parts[1]
			// log.Println("TrackID: ", trackID)
			iTrackID, err := strconv.Atoi(trackID)
			if err != nil {
				log.Fatal(err)
			}
			SetList.TrackCount = iTrackID
		}
		newTrack.ID = SetList.TrackCount

		if strings.HasPrefix(line, "\t\tPERFORMER") {
			// split the line into an array of strings
			// using the space as the separator
			parts := strings.Split(line, "\"")
			// get the number of tracks
			performer := parts[1]
			// log.Println("Performer: ", performer)
			newTrack.Performer = performer
		}
		if strings.HasPrefix(line, "\t\tTITLE") {
			// split the line into an array of strings
			// using the space as the separator
			parts := strings.Split(line, "\"")
			// get the number of tracks
			title := parts[1]
			// log.Println("Title: ", title)
			newTrack.Title = title
		}
		if strings.HasPrefix(line, "\t\tINDEX 01") {
			// split the line into an array of strings
			// using the space as the separator
			parts := strings.Split(line, " ")
			// get the number of tracks
			startTime := parts[2]
			// log.Println("Start_Time: ", startTime)
			newTrack.StartTime = startTime
			SetList.Tracks = append(SetList.Tracks, *newTrack)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ExportSetList(*SetList, o)
	log.Printf("Exported %d tracks\n", SetList.TrackCount)

}

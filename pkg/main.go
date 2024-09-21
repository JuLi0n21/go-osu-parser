package main

import (
	"fmt"
	osuParser "go-osu-parser/parser"
	"log"
	"os"
	"time"
)

func main() {
	filename := "G://Anwendungen/osu!/osu!.db"
	collname := "G://Anwendungen/osu!/collection.db"
	scoresname := "G://Anwendungen/osu!/scores.db"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("osu!.db file does not exist at path: %s", filename)
	}

	var start = time.Now()
	db, err := osuParser.ParseOsuDB(filename)
	if err != nil {
		log.Fatalf("Failed to parse osu!.db: %v", err)
	}
	fmt.Println("Parsed in: ", time.Since(start))

	start = time.Now()
	collection, err := osuParser.ParseCollectionsDB(collname)
	if err != nil {
		log.Fatalf("Failed to parse collections!.db: %v", err)
	}
	fmt.Println("Parsed in: ", time.Since(start))

	start = time.Now()
	scores, err := osuParser.ParseScoresDB(scoresname)
	if err != nil {
		log.Fatalf("Failed to parse scores!.db: %v", err)
	}

	fmt.Println("Parsed in: ", time.Since(start))

	fmt.Println("Collections", collection.NumberOfCollections)
	fmt.Println("Scores", scores.NumberOfScores)

	fmt.Printf("Osu! Version: %d\n", db.Version)
	fmt.Printf("Player Name: %s\n", db.PlayerName)
	fmt.Printf("Number of Beatmaps: %d\n", db.NumberOfBeatmaps)
	fmt.Printf("User Permissions: %d\n", db.UserPermissions)

}

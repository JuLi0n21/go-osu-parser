package main

import (
	"fmt"
	"log"
	"os"
	"time"

	osuParser "github.com/juli0n21/go-osu-parser/parser"
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

	fmt.Print("Parsing all .osu files. this will take sometime...\n\n")

	var SotarksCount int
	var TotalSotarksCircels int

	start = time.Now()
	for i, beatmap := range db.Beatmaps {

		b, err := osuParser.ParseOsuFile(fmt.Sprintf("G://Anwendungen/osu!/Songs/%s/%s", beatmap.FolderName, beatmap.FileName))
		if err != nil {
			log.Printf("Failed to parse osuFile: %v", err)
			continue
		}

		if b.Creator == "Sotarks" {
			SotarksCount++
			TotalSotarksCircels += len(b.HitObjects)
		}
		fmt.Printf("\033[F\r")
		fmt.Printf("\033[K")
		fmt.Printf("%d/%d\n", i, db.NumberOfBeatmaps)
	}

	fmt.Println("All .osu files parsed in: ", time.Since(start))
	fmt.Printf("Found %d Sotarks Diffs. With a total of %d circles/sliders", SotarksCount, TotalSotarksCircels)

}

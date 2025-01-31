package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/juli0n21/go-osu-parser/parser"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	f, err := os.Create("heap.prof")
	if err != nil {
		fmt.Println("Error creating heap profile:", err)
		return
	}
	defer f.Close()

	filename := "G://Anwendungen/osu!/osu!.db"
	collname := "G://Anwendungen/osu!/collection.db"
	scoresname := "G://Anwendungen/osu!/scores.db"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("osu!.db file does not exist at path: %s", filename)
	}

	var start = time.Now()
	db, err := parser.ParseOsuDB(filename)
	if err != nil {
		log.Fatalf("Failed to parse osu!.db: %v", err)
	}
	fmt.Println("Parsed in: ", time.Since(start))

	start = time.Now()
	collection, err := parser.ParseCollectionsDB(collname)
	if err != nil {
		log.Fatalf("Failed to parse collections!.db: %v", err)
	}
	fmt.Println("Parsed in: ", time.Since(start))

	start = time.Now()
	scores, err := parser.ParseScoresDB(scoresname)
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

	var wg sync.WaitGroup
	var mu sync.Mutex

	start = time.Now()
	for _, beatmap := range db.Beatmaps {
		wg.Add(1)
		go func(beatmap *parser.Beatmap) {
			defer wg.Done()

			b, err := parser.ParseOsuFile(fmt.Sprintf("G://Anwendungen/osu!/Songs/%s/%s", beatmap.FolderName, beatmap.FileName))
			if err != nil {
				log.Printf("Failed to parse osuFile: %v", err)
				return
			}

			if b.Creator == "Sotarks" {
				mu.Lock()
				SotarksCount++
				TotalSotarksCircels += len(b.HitObjects)
				mu.Unlock()
			}
		}(beatmap)
	}

	fmt.Printf("Parsed: %s beatmaps", len(db.Beatmaps))
	wg.Wait()

	fmt.Println("All .osu files parsed in: ", time.Since(start))
	fmt.Printf("Found %d Sotarks Diffs. With a total of %d circles/sliders", SotarksCount, TotalSotarksCircels)

	pprof.WriteHeapProfile(f)
}

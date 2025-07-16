package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/juli0n21/go-osu-parser/parser"
)

func TestParseOsuDb(t *testing.T) {
	filename := "./testdata/osu!.db"

	db, err := parser.ParseOsuDB(filename)
	if err != nil {
		t.Errorf("ParsingOsuDb failed: %v", err)
	}

	if db.PlayerName != "JuLi0n_" {
		t.Errorf("Playername mismatch: %s, expexted: JuLi0n_", db.PlayerName)
	}

	if len(db.Beatmaps) != 1 {
		t.Errorf("Beatmapcount mismatch: %d, expected: 1", len(db.Beatmaps))
	}

	if db.Beatmaps[0].BeatmapID != 1853515 {
		t.Errorf("Beatmapid mismatch: %d, expected: 1853515", db.Beatmaps[0].BeatmapID)
	}

	starRating := 9.92
	mods := parser.HardRock | parser.Hidden | parser.DoubleTime
	if db.Beatmaps[0].StarRatingsStandard[int(mods)]-float32(starRating) > 0.01 && db.Beatmaps[0].StarRatingsStandard[int(mods)] < -0.01 {
		t.Errorf("Starratign doesnt not match: %f, expected: 9.92", db.Beatmaps[0].StarRatingsStandard[int(mods)])
	}

}

func TestParseCollectionDb(t *testing.T) {
	filename := "./testdata/collection.db"

	collections, err := parser.ParseCollectionsDB(filename)
	if err != nil {
		t.Errorf("Parsing Collection failed: %v", err)
	}

	if collections.NumberOfCollections != 1 {
		t.Errorf("Expected 1 collection, but got %d", collections.NumberOfCollections)
	}

	if collections.Collections[0].Name != "cool collection" {
		t.Errorf("Expected collection name 'cool collection', but got '%s'", collections.Collections[0].Name)
	}

	if collections.Collections[0].NumberOfBeatmaps != 1 {
		t.Errorf("Expected 1 beatmap in collection 1, but got %d", collections.Collections[0].NumberOfBeatmaps)
	}

	beatmapID := collections.Collections[0].Beatmaps[0]
	if *beatmapID != "b698e35323632c828030a7c366c5c933" {
		t.Errorf("Expected beatmap ID 'b698e35323632c828030a7c366c5c933', but got '%s'", *beatmapID)
	}

}

func TestParseScoresDb(t *testing.T) {
	filename := "./testdata/scores.db"

	scores, err := parser.ParseScoresDB(filename)
	if err != nil {
		t.Errorf("Parsing Presence failed: %v", err)
	}

	if scores.NumberOfScores != 1 {
		t.Errorf("Expected 1 score, but got %d", scores.NumberOfScores)
	}

	if len(scores.Beatmaps) == 0 {
		t.Errorf("Expected at least one BeatmapScores, but found none")
	}

	expectedMD5Hash := "b698e35323632c828030a7c366c5c933"
	if scores.Beatmaps[0].BeatmapMD5Hash != expectedMD5Hash {
		t.Errorf("Expected BeatmapMD5Hash '%s', but got '%s'", expectedMD5Hash, scores.Beatmaps[0].BeatmapMD5Hash)
	}

	if scores.Beatmaps[0].NumberOfScores != 1 {
		t.Errorf("Expected 1 score for the beatmap, but got %d", scores.Beatmaps[0].NumberOfScores)
	}

	score := scores.Beatmaps[0].Scores[0]

	expectedMaxCombo := uint16(468)
	if score.MaxCombo != expectedMaxCombo {
		t.Errorf("Expected MaxCombo %d, but got %d", expectedMaxCombo, score.MaxCombo)
	}

	expectedMods := int32(0)
	if score.Mods != expectedMods {
		t.Errorf("Expected Mods %d, but got %d", expectedMods, score.Mods)
	}

	if score.PerfectCombo {
		t.Error("Expected PerfectCombo to be false, but got true")
	}

	expectedPlayerName := "JuLi0n_"
	if score.PlayerName != expectedPlayerName {
		t.Errorf("Expected PlayerName '%s', but got '%s'", expectedPlayerName, score.PlayerName)
	}

}

func TestFileParsing(t *testing.T) {
	//scan file tree of ./testdata for .osu files and parse them all
	testDir := "./testdata"

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".osu" {
			t.Logf("Parsing file: %s", path)
			_, err := parser.ParseOsuFile(path)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", path, err)
			}

		}
		return nil
	})

	if err != nil {
		t.Fatalf("Error walking through testdata directory: %v", err)
	}
}

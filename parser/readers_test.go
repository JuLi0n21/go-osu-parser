package parser_test

import (
	"testing"

	"github.com/juli0n21/go-osu-parser/parser"
)

func TestParseOsuFile(t *testing.T) {
	filename := "./testdata/1853515 MUZZ - Endgame/MUZZ - Endgame (Time Freeze) [When the Evil Descends].osu"

	osuFile, err := parser.ParseOsuFile(filename)
	if err != nil {
		t.Errorf("Parsing osu file failed: %v", err)
	}

	if osuFile.General.AudioFilename != "audio.mp3" {
		t.Errorf("Expected AudioFilename 'audio.mp3', but got '%s'", osuFile.General.AudioFilename)
	}

	if osuFile.General.Mode != 0 {
		t.Errorf("Expected Mode 0, but got %d", osuFile.General.Mode)
	}

	if osuFile.General.StackLeniency != 0.7 {
		t.Errorf("Expected StackLeniency 0.7, but got %.2f", osuFile.General.StackLeniency)
	}

	if osuFile.Metadata.Title != "Endgame" {
		t.Errorf("Expected Title 'Endgame', but got '%s'", osuFile.Metadata.Title)
	}

	if osuFile.Metadata.Artist != "MUZZ" {
		t.Errorf("Expected Artist 'MUZZ', but got '%s'", osuFile.Metadata.Artist)
	}

	if osuFile.Metadata.Creator != "Time Freeze" {
		t.Errorf("Expected Creator 'Time Freeze', but got '%s'", osuFile.Metadata.Creator)
	}

	if osuFile.Difficulty.HPDrainRate != 5.0 {
		t.Errorf("Expected HPDrainRate 5.0, but got %.2f", osuFile.Difficulty.HPDrainRate)
	}

	if osuFile.Difficulty.CircleSize != 3.80 {
		t.Errorf("Expected CircleSize 3.80, but got %.2f", osuFile.Difficulty.CircleSize)
	}

	if len(osuFile.TimingPointsFile) > 0 {
		firstTimingPoint := osuFile.TimingPointsFile[0]
		if 342.86-firstTimingPoint.BeatLength > 0.01 {
			t.Errorf("Expected BeatLength 342.86, but got %f", firstTimingPoint.BeatLength)
		}

		if firstTimingPoint.Volume != 40 {
			t.Errorf("Expected Volume 40, but got %d", firstTimingPoint.Volume)
		}
	} else {
		t.Error("Expected at least one timing point, but found none")
	}

	if len(osuFile.HitObjects) > 0 {
		firstHitObject := osuFile.HitObjects[0]
		if firstHitObject.X != 32.0 {
			t.Errorf("Expected X 32.0, but got %.2f", firstHitObject.X)
		}

		if firstHitObject.Y != 336.0 {
			t.Errorf("Expected Y 336.0, but got %.2f", firstHitObject.Y)
		}

		if firstHitObject.Time != 37035.0 {
			t.Errorf("Expected Time 37035.0, but got %.2f", firstHitObject.Time)
		}

		if firstHitObject.Type != 5 {
			t.Errorf("Expected Type 5, but got %d", firstHitObject.Type)
		}
	} else {
		t.Error("Expected at least one hit object, but found none")
	}
}

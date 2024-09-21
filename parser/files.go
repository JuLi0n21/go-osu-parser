package osuParser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type General struct {
	AudioFilename            string
	AudioLeadIn              int
	AudioHash                string
	PreviewTime              int
	Countdown                int
	SampleSet                string
	StackLeniency            float64
	Mode                     int
	LetterboxInBreaks        int
	StoryFireInFront         int
	UseSkinSprites           int
	AlwaysShowPlayfield      int
	OverlayPosition          string
	SkinPreference           string
	EpilepsyWarning          int
	CountdownOffset          int
	SpecialStyle             int
	WidescreenStoryboard     int
	SamplesMatchPlaybackRate int
}

type Editor struct {
	Bookmarks       []int
	DistanceSpacing float64
	BeatDivisor     int
	GridSize        int
	TimelineZoom    float64
}

type Metadata struct {
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Creator       string
	Version       string
	Source        string
	Tags          []string
	BeatmapID     int
	BeatmapSetID  int
}

type Difficulty struct {
	HPDrainRate       float64
	CircleSize        float64
	OverallDifficulty float64
	ApproachRate      float64
	SliderMultiplier  float64
	SliderTickRate    float64
}

type Event struct {
	EventType   string
	StartTime   int
	EventParams []string
}

type TimingPointFile struct {
	Time        int
	BeatLength  float64
	Meter       int
	SampleSet   int
	SampleIndex int
	Volume      int
	Uninherited int
	Effects     int
}

type Colour struct {
	Combo               []int
	SliderTrackOverride []int
	SliderBorder        []int
}

type HitObject struct {
	X            float64
	Y            float64
	Time         float64
	Type         int
	HitSound     int
	ObjectParams string
	HitSample    string
}

type OsuFile struct {
	Version int
	General
	Editor
	Metadata
	Difficulty
	Events           []Event
	TimingPointsFile []TimingPointFile
	Colours          []Colour
	HitObjects       []HitObject
}

// not yet Implemented
type ReplayFile struct {
	Gamemode                 byte
	Version                  int32
	beatmapMD5Hash           string
	playername               string
	replayMD5Hash            string
	Count300s                int16
	Count100s                int16
	Count50s                 int16
	Gekis                    int16
	Katus                    int16
	CountMiss                int16
	Combo                    int16
	Score                    int32
	PerfectCombo             byte
	Mods                     int32
	HealthGraph              []*Health
	Timestamp                int64
	LengthInBytes            int32
	Replay                   []*ReplayData
	OnlineScoreId            int64
	AdditionalModInformation float64
	LZMA                     []*byte
}

type Health struct {
	u int32
	v float32
}

type ReplayData struct {
}

func ParseOsuFile(filename string) (*OsuFile, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Unexpected error occored: %v", r)
		}
	}()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s in file: %s", err, filename)
	}

	OsuFile, err := parseOsuFile(filename)
	if err != nil {
		return nil, err
	}

	return OsuFile, nil

}

func parseOsuFile(filename string) (*OsuFile, error) {
	var err error
	var lineNumber int
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v in line %d", r, lineNumber)
		}
	}()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 128*1024)

	byteData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split([]byte(byteData), []byte{'\n'})
	osuFile := &OsuFile{}
	currentSection := ""

	for i, lineStr := range lines {
		lineNumber = i

		line := strings.TrimSpace(string(lineStr))

		if len(line) == 0 || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(line[1 : len(line)-1])
			continue
		}

		switch currentSection {
		case "general":
			parseGeneral(line, &osuFile.General)
		case "editor":
			parseEditor(line, &osuFile.Editor)
		case "metadata":
			parseMetadata(line, &osuFile.Metadata)
		case "difficulty":
			parseDifficulty(line, &osuFile.Difficulty)
		case "events":
			parseEvents(line, &osuFile.Events)
		case "timingpoints":
			parseTimingPoints(line, &osuFile.TimingPointsFile)
		case "colours":
			parseColours(line, &osuFile.Colours)
		case "hitobjects":
			parseHitObjects(line, &osuFile.HitObjects)
		}
	}

	return osuFile, nil
}

func parseGeneral(line string, general *General) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "AudioFilename":
		general.AudioFilename = value
	case "AudioLeadIn":
		general.AudioLeadIn, _ = strconv.Atoi(value)
	case "PreviewTime":
		general.PreviewTime, _ = strconv.Atoi(value)
	case "Countdown":
		general.Countdown, _ = strconv.Atoi(value)
	case "SampleSet":
		general.SampleSet = value
	case "StackLeniency":
		general.StackLeniency, _ = strconv.ParseFloat(value, 64)
	case "Mode":
		general.Mode, _ = strconv.Atoi(value)
	case "LetterboxInBreaks":
		general.LetterboxInBreaks, _ = strconv.Atoi(value)
	case "WidescreenStoryboard":
		general.WidescreenStoryboard, _ = strconv.Atoi(value)
	}
}

func parseEditor(line string, editor *Editor) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Bookmarks":
		for _, v := range strings.Split(value, ",") {
			bookmark, _ := strconv.Atoi(v)
			editor.Bookmarks = append(editor.Bookmarks, bookmark)
		}
	case "DistanceSpacing":
		editor.DistanceSpacing, _ = strconv.ParseFloat(value, 64)
	case "BeatDivisor":
		editor.BeatDivisor, _ = strconv.Atoi(value)
	case "GridSize":
		editor.GridSize, _ = strconv.Atoi(value)
	case "TimelineZoom":
		editor.TimelineZoom, _ = strconv.ParseFloat(value, 64)
	}
}

func parseMetadata(line string, metadata *Metadata) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Title":
		metadata.Title = value
	case "TitleUnicode":
		metadata.TitleUnicode = value
	case "Artist":
		metadata.Artist = value
	case "ArtistUnicode":
		metadata.ArtistUnicode = value
	case "Creator":
		metadata.Creator = value
	case "Version":
		metadata.Version = value
	case "Source":
		metadata.Source = value
	case "Tags":
		metadata.Tags = strings.Split(value, " ")
	case "BeatmapID":
		metadata.BeatmapID, _ = strconv.Atoi(value)
	case "BeatmapSetID":
		metadata.BeatmapSetID, _ = strconv.Atoi(value)
	}
}

func parseDifficulty(line string, difficulty *Difficulty) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "HPDrainRate":
		difficulty.HPDrainRate, _ = strconv.ParseFloat(value, 64)
	case "CircleSize":
		difficulty.CircleSize, _ = strconv.ParseFloat(value, 64)
	case "OverallDifficulty":
		difficulty.OverallDifficulty, _ = strconv.ParseFloat(value, 64)
	case "ApproachRate":
		difficulty.ApproachRate, _ = strconv.ParseFloat(value, 64)
	case "SliderMultiplier":
		difficulty.SliderMultiplier, _ = strconv.ParseFloat(value, 64)
	case "SliderTickRate":
		difficulty.SliderTickRate, _ = strconv.ParseFloat(value, 64)
	}
}

func parseEvents(line string, events *[]Event) {
	parts := strings.Split(line, ",")

	if len(parts) < 1 {
		return
	}

	eventType := parts[0]
	startTime := 0
	var eventParams []string

	if len(parts) > 1 {
		startTime, _ = strconv.Atoi(parts[1])
	}

	if len(parts) > 2 {
		eventParams = parts[2:]
	}

	*events = append(*events, Event{
		EventType:   eventType,
		StartTime:   startTime,
		EventParams: eventParams,
	})
}

func parseTimingPoints(line string, timingPoints *[]TimingPointFile) {
	parts := strings.Split(line, ",")
	if len(parts) < 8 {
		return
	}

	time, _ := strconv.Atoi(parts[0])
	beatLength, _ := strconv.ParseFloat(parts[1], 64)
	meter, _ := strconv.Atoi(parts[2])
	sampleSet, _ := strconv.Atoi(parts[3])
	sampleIndex, _ := strconv.Atoi(parts[4])
	volume, _ := strconv.Atoi(parts[5])
	uninherited, _ := strconv.Atoi(parts[6])
	effects, _ := strconv.Atoi(parts[7])

	*timingPoints = append(*timingPoints, TimingPointFile{
		Time:        time,
		BeatLength:  beatLength,
		Meter:       meter,
		SampleSet:   sampleSet,
		SampleIndex: sampleIndex,
		Volume:      volume,
		Uninherited: uninherited,
		Effects:     effects,
	})
}

func parseColours(line string, colours *[]Colour) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	rgb := strings.Split(value, ",")

	if len(rgb) != 3 {
		return
	}

	r, _ := strconv.Atoi(rgb[0])
	g, _ := strconv.Atoi(rgb[1])
	b, _ := strconv.Atoi(rgb[2])

	switch {
	case strings.HasPrefix(key, "Combo"):
		*colours = append(*colours, Colour{Combo: []int{r, g, b}})
	case key == "SliderTrackOverride":
		*colours = append(*colours, Colour{SliderTrackOverride: []int{r, g, b}})
	case key == "SliderBorder":
		*colours = append(*colours, Colour{SliderBorder: []int{r, g, b}})
	}
}

func parseHitObjects(line string, hitObjects *[]HitObject) {
	parts := strings.Split(line, ",")
	if len(parts) < 5 {
		return
	}

	x, _ := strconv.ParseFloat(parts[0], 64)
	y, _ := strconv.ParseFloat(parts[1], 64)
	time, _ := strconv.ParseFloat(parts[2], 64)
	objectType, _ := strconv.Atoi(parts[3])
	hitSound, _ := strconv.Atoi(parts[4])

	objectParams := ""
	if len(parts) > 5 {
		objectParams = parts[5]
	}
	hitSample := ""
	if len(parts) > 6 {
		hitSample = parts[6]
	}

	*hitObjects = append(*hitObjects, HitObject{
		X:            x,
		Y:            y,
		Time:         time,
		Type:         objectType,
		HitSound:     hitSound,
		ObjectParams: objectParams,
		HitSample:    hitSample,
	})
}

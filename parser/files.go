package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type General struct {
	StackLeniency            float64
	AudioLeadIn              int
	EpilepsyWarning          int
	CountdownOffset          int
	SpecialStyle             int
	WidescreenStoryboard     int
	SamplesMatchPlaybackRate int
	PreviewTime              int
	Countdown                int
	Mode                     int
	LetterboxInBreaks        int
	StoryFireInFront         int
	UseSkinSprites           int
	AlwaysShowPlayfield      int
	AudioHash                string
	SampleSet                string
	OverlayPosition          string
	SkinPreference           string
	AudioFilename            string
}

type Editor struct {
	DistanceSpacing float64
	TimelineZoom    float64
	BeatDivisor     int
	GridSize        int
	Bookmarks       []int
}

type Metadata struct {
	BeatmapID     int
	BeatmapSetID  int
	Tags          []string
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Creator       string
	Version       string
	Source        string
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
	StartTime   int
	EventParams []string
	EventType   string
}

type TimingPointFile struct {
	BeatLength  float64
	Time        int
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
	AdditionalModInformation float64
	Timestamp                int64
	OnlineScoreId            int64
	LengthInBytes            int32
	Mods                     int32
	Score                    int32
	Version                  int32
	Count300s                int16
	Count100s                int16
	Count50s                 int16
	Gekis                    int16
	Katus                    int16
	CountMiss                int16
	Combo                    int16
	Gamemode                 byte
	PerfectCombo             byte
	HealthGraph              []*Health
	Replay                   []*ReplayData
	LZMA                     []*byte
	beatmapMD5Hash           string
	playername               string
	replayMD5Hash            string
}

type Health struct {
	u int32
	v float32
}

type ReplayData struct {
}

func ParseOsuFile(filename string) (osufile *OsuFile, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unexpected error occurred: %v", r)
		}
	}()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s, error: %w", filename, err)
	}

	OsuFile, err := parseOsuFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse osufile: %s, error: %w", filename, err)
	}
	if OsuFile == nil {
		return nil, fmt.Errorf("parsed OsuFile is nil: %s", filename)
	}

	return OsuFile, nil
}

func parseOsuFile(filename string) (*OsuFile, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to parse osufile: %s, panic: %v", filename, r)
		}
	}()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s, error: %w", filename, err)
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 128*1024)

	byteData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytedata: %w", err)
	}

	byteData = bytes.TrimPrefix(byteData, []byte{0xEF, 0xBB, 0xBF})
	lines := bytes.Split(byteData, []byte{'\n'})
	osuFile := &OsuFile{}
	currentSection := ""
	var parseErrs []error

	for i, lineStr := range lines {
		if i == 0 {
			line := strings.TrimSpace(string(lineStr))
			prefix := "osu file format v"
			if strings.HasPrefix(line, prefix) {
				versionStr := strings.TrimSpace(line[len(prefix):])
				osuFile.Version, err = strconv.Atoi(versionStr)
				if err != nil {
					parseErrs = append(parseErrs, fmt.Errorf("failed to parse version: %w", err))
				}
			} else {
				parseErrs = append(parseErrs, fmt.Errorf("unexpected first line format: %s", line))
			}
		}

		line := strings.TrimSpace(string(lineStr))
		if len(line) == 0 || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "[") {
			currentSection = strings.ToLower(line[1 : len(line)-1])
			continue
		}

		switch currentSection {
		case "general":
			if e := parseGeneral(line, &osuFile.General); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("general: %w", e))
			}
		case "editor":
			if e := parseEditor(line, &osuFile.Editor); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("editor: %w", e))
			}
		case "metadata":
			if e := parseMetadata(line, &osuFile.Metadata); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("metadata: %w", e))
			}
		case "difficulty":
			if e := parseDifficulty(line, &osuFile.Difficulty); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("difficulty: %w", e))
			}
		case "events":
			if e := parseEvents(line, &osuFile.Events); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("events: %w", e))
			}
		case "timingpoints":
			if e := parseTimingPoints(line, &osuFile.TimingPointsFile); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("timingpoints: %w", e))
			}
		case "colours":
			if e := parseColours(line, &osuFile.Colours); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("colours: %w", e))
			}
		case "hitobjects":
			if e := parseHitObjects(line, &osuFile.HitObjects); e != nil {
				parseErrs = append(parseErrs, fmt.Errorf("hitobjects: %w", e))
			}
		}
	}

	if len(parseErrs) > 0 {
		return osuFile, errors.Join(parseErrs...)
	}

	return osuFile, nil
}

func parseGeneral(line string, general *General) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid line format: %q", line)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	var err error
	switch key {
	case "AudioFilename":
		general.AudioFilename = value

	case "AudioLeadIn":
		var v int
		v, err = strconv.Atoi(value)
		general.AudioLeadIn = v

	case "PreviewTime":
		var v int
		v, err = strconv.Atoi(value)
		general.PreviewTime = v

	case "Countdown":
		var v int
		v, err = strconv.Atoi(value)
		general.Countdown = v

	case "SampleSet":
		general.SampleSet = value

	case "StackLeniency":
		var v float64
		v, err = strconv.ParseFloat(value, 64)
		general.StackLeniency = v

	case "Mode":
		var v int
		v, err = strconv.Atoi(value)
		general.Mode = v

	case "LetterboxInBreaks":
		var v int
		v, err = strconv.Atoi(value)
		general.LetterboxInBreaks = v

	case "WidescreenStoryboard":
		var v int
		v, err = strconv.Atoi(value)
		general.WidescreenStoryboard = v
	}

	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", key, err)
	}
	return nil
}

func parseEditor(line string, editor *Editor) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid line format: %q", line)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	var err error
	switch key {
	case "Bookmarks":
		for _, v := range strings.Split(value, ",") {
			if v == "" {
				continue
			}
			var bookmark int
			bookmark, err = strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				return fmt.Errorf("failed to parse bookmark %q: %w", v, err)
			}
			editor.Bookmarks = append(editor.Bookmarks, bookmark)
		}

	case "DistanceSpacing":
		editor.DistanceSpacing, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse DistanceSpacing %q: %w", value, err)
		}

	case "BeatDivisor":
		editor.BeatDivisor, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to parse BeatDivisor %q: %w", value, err)
		}

	case "GridSize":
		editor.GridSize, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to parse GridSize %q: %w", value, err)
		}

	case "TimelineZoom":
		editor.TimelineZoom, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse TimelineZoom %q: %w", value, err)
		}
	}

	return nil
}

func parseMetadata(line string, metadata *Metadata) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid line format: %q", line)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	var err error
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
		metadata.BeatmapID, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to parse BeatmapID %q: %w", value, err)
		}
	case "BeatmapSetID":
		metadata.BeatmapSetID, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to parse BeatmapSetID %q: %w", value, err)
		}
	}

	return nil
}
func parseDifficulty(line string, difficulty *Difficulty) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid line format: %q", line)
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	var err error
	switch key {
	case "HPDrainRate":
		difficulty.HPDrainRate, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse HPDrainRate %q: %w", value, err)
		}
	case "CircleSize":
		difficulty.CircleSize, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse CircleSize %q: %w", value, err)
		}
	case "OverallDifficulty":
		difficulty.OverallDifficulty, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse OverallDifficulty %q: %w", value, err)
		}
	case "ApproachRate":
		difficulty.ApproachRate, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse ApproachRate %q: %w", value, err)
		}
	case "SliderMultiplier":
		difficulty.SliderMultiplier, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse SliderMultiplier %q: %w", value, err)
		}
	case "SliderTickRate":
		difficulty.SliderTickRate, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse SliderTickRate %q: %w", value, err)
		}
	}

	return nil
}

var allowedEventTypes map[string]bool = map[string]bool{
	"0":     true, // Background
	"1":     true, // Video (numeric form)
	"2":     true, // Break
	"Video": true, // Video (string form)
	"Break": true, // Break (string form)
}

func parseEvents(line string, events *[]Event) error {
	parts := strings.Split(line, ",")
	if len(parts) < 1 {
		return fmt.Errorf("invalid event line: %q", line)
	}

	eventType := parts[0]

	if !allowedEventTypes[eventType] {
		return nil
	}

	startTime := 0
	var err error

	if eventType == "2" || eventType == "Break" {
		if len(parts) < 3 {
			return fmt.Errorf("invalid break event line: %q", line)
		}
		startTime, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("failed to parse startTime from %q: %w", parts[1], err)
		}
		eventParams := parts[2:]
		*events = append(*events, Event{
			EventType:   eventType,
			StartTime:   startTime,
			EventParams: eventParams,
		})
		return nil
	}

	if len(parts) > 1 {
		startTime, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("failed to parse startTime from %q: %w", parts[1], err)
		}
	}

	eventParams := []string{}
	if len(parts) > 2 {
		eventParams = parts[2:]
	}

	*events = append(*events, Event{
		EventType:   eventType,
		StartTime:   startTime,
		EventParams: eventParams,
	})

	return nil
}

func parseTimingPoints(line string, timingPoints *[]TimingPointFile) error {
	parts := strings.Split(line, ",")

	timeFloat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("failed to parse Time from %q: %w", parts[0], err)
	}
	time := int(timeFloat)

	beatLength, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("failed to parse BeatLength from %q: %w", parts[1], err)
	}

	meter, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("failed to parse Meter from %q: %w", parts[2], err)
	}

	sampleSet, err := strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("failed to parse SampleSet from %q: %w", parts[3], err)
	}

	sampleIndex, err := strconv.Atoi(parts[4])
	if err != nil {
		return fmt.Errorf("failed to parse SampleIndex from %q: %w", parts[4], err)
	}

	volume, err := strconv.Atoi(parts[5])
	if err != nil {
		return fmt.Errorf("failed to parse Volume from %q: %w", parts[5], err)
	}

	uninherited, err := strconv.Atoi(parts[6])
	if err != nil {
		return fmt.Errorf("failed to parse Uninherited from %q: %w", parts[6], err)
	}

	//fall back for older versions...
	effects := 0
	if len(parts) == 8 {
		effects, err = strconv.Atoi(parts[7])
		if err != nil {
			return fmt.Errorf("failed to parse Effects from %q: %w", parts[7], err)
		}
	}

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

	return nil
}

func parseColours(line string, colours *[]Colour) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid colour line (missing colon): %q", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	rgb := strings.Split(value, ",")

	if len(rgb) != 3 {
		return fmt.Errorf("invalid RGB value (expected 3 components): %q", value)
	}

	r, err := strconv.Atoi(strings.TrimSpace(rgb[0]))
	if err != nil {
		return fmt.Errorf("failed to parse R from %q: %w", rgb[0], err)
	}
	g, err := strconv.Atoi(strings.TrimSpace(rgb[1]))
	if err != nil {
		return fmt.Errorf("failed to parse G from %q: %w", rgb[1], err)
	}
	b, err := strconv.Atoi(strings.TrimSpace(rgb[2]))
	if err != nil {
		return fmt.Errorf("failed to parse B from %q: %w", rgb[2], err)
	}

	switch {
	case strings.HasPrefix(key, "Combo"):
		*colours = append(*colours, Colour{Combo: []int{r, g, b}})
	case key == "SliderTrackOverride":
		*colours = append(*colours, Colour{SliderTrackOverride: []int{r, g, b}})
	case key == "SliderBorder":
		*colours = append(*colours, Colour{SliderBorder: []int{r, g, b}})
	default:
		return fmt.Errorf("unrecognized colour key: %q", key)
	}

	return nil
}

func parseHitObjects(line string, hitObjects *[]HitObject) error {
	parts := strings.Split(line, ",")
	if len(parts) < 5 {
		return fmt.Errorf("invalid hit object line (expected at least 5 parts): %q", line)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("failed to parse X from %q: %w", parts[0], err)
	}

	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("failed to parse Y from %q: %w", parts[1], err)
	}

	time, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return fmt.Errorf("failed to parse Time from %q: %w", parts[2], err)
	}

	objectType, err := strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("failed to parse Type from %q: %w", parts[3], err)
	}

	hitSound, err := strconv.Atoi(parts[4])
	if err != nil {
		return fmt.Errorf("failed to parse HitSound from %q: %w", parts[4], err)
	}

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

	return nil
}

func (o *OsuFile) BackgroundImage() string {
	if o == nil || o.Events == nil {
		return ""
	}

	for _, event := range o.Events {
		if event.EventType == "0" || strings.EqualFold(event.EventType, "Background") {
			if len(event.EventParams) > 0 {
				return strings.Trim(event.EventParams[0], "\"")
			}
		}
	}
	return ""
}

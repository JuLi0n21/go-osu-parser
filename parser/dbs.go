package parser

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

type OsuDB struct {
	UnlockDate       time.Time
	Version          int32
	FolderCount      int32
	UserPermissions  int32
	NumberOfBeatmaps int32
	AccountUnlocked  bool
	Beatmaps         []*Beatmap
	PlayerName       string
}

type Beatmap struct {
	LastModificationTime   int64
	LastPlayed             int64
	LastChecked            int64
	DrainTime              int32
	TotalTime              int32
	AudioPreviewStartTime  int32
	DifficultyID           int32
	BeatmapID              int32
	ThreadID               int32
	LastModificationTime2  int32
	OnlineOffset           int16
	LocalBeatmapOffset     uint16
	NumberOfHitCircles     uint16
	NumberOfSliders        uint16
	NumberOfSpinners       uint16
	RankedStatus           byte
	GradeStandard          byte
	GradeTaiko             byte
	GradeCTB               byte
	GradeMania             byte
	ManiaScrollSpeed       byte
	GameplayMode           byte
	IsUnplayed             bool
	IsOsz2                 bool
	IgnoreBeatmapSound     bool
	IgnoreBeatmapSkin      bool
	DisableStoryboard      bool
	DisableVideo           bool
	VisualOverride         bool
	ApproachRate           float32
	CircleSize             float32
	HPDrain                float32
	OverallDifficulty      float32
	StackLeniency          float32
	SliderVelocity         float64
	Artist                 string
	ArtistUnicode          string
	SongTitle              string
	SongTitleUnicode       string
	Creator                string
	Difficulty             string
	AudioFileName          string
	MD5Hash                string
	FileName               string
	SongSource             string
	SongTags               string
	Font                   string
	FolderName             string
	StarRatingsStandard    map[int]float32
	StarRatingsTaiko       map[int]float32
	StarRatingsCTB         map[int]float32
	StarRatingsMania       map[int]float32
	StarRatingsStandardOld map[int]float64
	StarRatingsTaikoOld    map[int]float64
	StarRatingsCTBOld      map[int]float64
	StarRatingsManiaOld    map[int]float64
	TimingPoints           []TimingPoint
	SizeInBytes            *int32
	UnknownShort           *uint16
}

type TimingPoint struct {
	BPM       float64
	Offset    float64
	Inherited bool
}

type Collections struct {
	Version             int32
	NumberOfCollections int32
	Collections         []*Collection
}

type Collection struct {
	NumberOfBeatmaps int32
	Beatmaps         []*string
	Name             string
}

type Scores struct {
	Version        int32
	NumberOfScores int32
	Beatmaps       []*BeatmapScores
}

type BeatmapScores struct {
	NumberOfScores int32
	BeatmapMD5Hash string
	Scores         []*Score
}

type Score struct {
	Timestamp         time.Time
	AdditionalModInfo float64
	OnlineScoreId     int64
	Version           int32
	Mods              int32
	ReplayScore       int32
	Count300s         uint16
	Count100s         uint16
	Count50           uint16
	Gekis             uint16
	Katus             uint16
	CountMiss         uint16
	MaxCombo          uint16
	Gamemode          byte
	PerfectCombo      bool
	BeatmapMD5Hash    string
	PlayerName        string
	ReplayMD5Hash     string
}

type Mods int

const (
	NoMod       Mods = 0
	Easy        Mods = 1 << iota // 1
	NoFail                       // 2
	HalfTime                     // 4
	HardRock                     // 8
	SuddenDeath                  // 16
	DoubleTime                   // 32
	Relax                        // 64
	Hidden                       // 128
	Flashlight                   // 256
	Autoplay                     // 512
	SpunOut                      // 1024
	Relax2                       // 2048 (Autopilot)
	Perfect                      // 4096
	Key4                         // 8192
	Key5                         // 16384
	Key6                         // 32768
	Key7                         // 65536
	Key8                         // 131072
	FadeIn                       // 262144
	Random                       // 524288
	Cinema                       // 1048576
	Target                       // 2097152
	Key9                         // 4194304
	Key10                        // 8388608
	Key1                         // 16777216
	Key3                         // 33554432
	Key2                         // 67108864
	ScoreV2                      // 134217728
	Mirror                       // 268435456
)

func ParseCollectionsDB(filename string) (*Collections, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 128*1024)

	version, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	collectionCount, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	collections := make([]*Collection, 0, collectionCount)
	for i := 0; i < int(collectionCount); i++ {
		collection, err := readCollection(reader)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	if err != nil {
		return nil, err
	}

	return &Collections{
		Version:             version,
		NumberOfCollections: collectionCount,
		Collections:         collections,
	}, nil

}

func readCollection(r io.Reader) (*Collection, error) {

	name, err := readString(r)
	if err != nil {
		return nil, err
	}

	beatmapCount, err := readInt(r)
	if err != nil {
		return nil, err
	}

	beatmaps := make([]*string, 0, beatmapCount)
	for i := 0; i < int(beatmapCount); i++ {
		beatmap, err := readString(r)
		if err != nil {
			return nil, err
		}
		beatmaps = append(beatmaps, &beatmap)
	}

	return &Collection{
		Name:             name,
		NumberOfBeatmaps: beatmapCount,
		Beatmaps:         beatmaps,
	}, nil

}

func ParseScoresDB(filename string) (*Scores, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 128*1024)

	version, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	scoreCount, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	beatmaps := make([]*BeatmapScores, 0, scoreCount)
	for i := 0; i < int(scoreCount); i++ {
		beatmap, err := readBeatmapScore(reader)
		if err != nil {
			return nil, err
		}
		beatmaps = append(beatmaps, beatmap)
	}

	if err != nil {
		return nil, err
	}

	return &Scores{
		Version:        version,
		NumberOfScores: scoreCount,
		Beatmaps:       beatmaps,
	}, nil

}

func readBeatmapScore(r io.Reader) (*BeatmapScores, error) {
	hash, err := readString(r)
	if err != nil {
		return nil, err
	}

	scoreCount, err := readInt(r)
	if err != nil {
		return nil, err
	}

	scores := make([]*Score, 0, scoreCount)
	for i := 0; i < int(scoreCount); i++ {
		score, err := readScore(r)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return &BeatmapScores{
		BeatmapMD5Hash: hash,
		NumberOfScores: scoreCount,
		Scores:         scores,
	}, nil

}

func readScore(r io.Reader) (*Score, error) {

	gamemode, err := readByte(r)
	if err != nil {
		return nil, err
	}

	version, err := readInt(r)
	if err != nil {
		return nil, err
	}

	beatmapMD5Hash, err := readString(r)
	if err != nil {
		return nil, err
	}

	playername, err := readString(r)
	if err != nil {
		return nil, err
	}

	replayMD5Hash, err := readString(r)
	if err != nil {
		return nil, err
	}

	count300, err := readShort(r)
	if err != nil {
		return nil, err
	}

	count100, err := readShort(r)
	if err != nil {
		return nil, err
	}

	count50, err := readShort(r)
	if err != nil {
		return nil, err
	}

	gekis, err := readShort(r)
	if err != nil {
		return nil, err
	}

	katus, err := readShort(r)
	if err != nil {
		return nil, err
	}

	countMiss, err := readShort(r)
	if err != nil {
		return nil, err
	}

	replayScore, err := readInt(r)
	if err != nil {
		return nil, err
	}
	maxcombo, err := readShort(r)
	if err != nil {
		return nil, err
	}

	perfectCombo, err := readBoolean(r)
	if err != nil {
		return nil, err
	}

	mods, err := readInt(r)
	if err != nil {
		return nil, err
	}

	//EmptyString
	_, _ = readString(r)

	ticks, err := readLong(r)
	if err != nil {
		return nil, err
	}
	timestamp := readDateTime(ticks)

	//-1
	_, err = readInt(r)
	if err != nil {
		return nil, err
	}

	onlineScoreId, err := readLong(r)
	if err != nil {
		return nil, err
	}

	var additionalModInfo float64
	if mods<<23 == 1 {
		additionalModInfo, err = readDouble(r)
		if err != nil {
			return nil, err
		}

	}

	return &Score{
		Gamemode:          gamemode,
		Version:           version,
		BeatmapMD5Hash:    beatmapMD5Hash,
		PlayerName:        playername,
		ReplayMD5Hash:     replayMD5Hash,
		Count300s:         count300,
		Count100s:         count100,
		Count50:           count50,
		Gekis:             gekis,
		Katus:             katus,
		CountMiss:         countMiss,
		ReplayScore:       replayScore,
		MaxCombo:          maxcombo,
		PerfectCombo:      perfectCombo,
		Mods:              mods,
		Timestamp:         timestamp,
		OnlineScoreId:     onlineScoreId,
		AdditionalModInfo: additionalModInfo,
	}, nil
}

func readBeatmap(r io.Reader, version int32) (*Beatmap, error) {
	beatmap := &Beatmap{}

	if version < 20191106 {
		sizeInBytes, err := readInt(r)
		if err != nil {
			return nil, err
		}
		beatmap.SizeInBytes = &sizeInBytes
	}

	artist, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.Artist = artist

	artistUnicode, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.ArtistUnicode = artistUnicode

	songTitle, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.SongTitle = songTitle

	songTitleUnicode, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.SongTitleUnicode = songTitleUnicode

	creator, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.Creator = creator

	difficulty, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.Difficulty = difficulty

	audioFileName, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.AudioFileName = audioFileName

	md5Hash, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.MD5Hash = md5Hash

	osuFileName, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.FileName = osuFileName

	var rankedStatus byte
	if err := binary.Read(r, binary.LittleEndian, &rankedStatus); err != nil {
		return nil, err
	}
	beatmap.RankedStatus = rankedStatus

	numberOfHitCircles, err := readShort(r)
	if err != nil {
		return nil, err
	}
	beatmap.NumberOfHitCircles = numberOfHitCircles

	numberOfSliders, err := readShort(r)
	if err != nil {
		return nil, err
	}
	beatmap.NumberOfSliders = numberOfSliders

	numberOfSpinners, err := readShort(r)
	if err != nil {
		return nil, err
	}
	beatmap.NumberOfSpinners = numberOfSpinners

	lastModificationTicks, err := readLong(r)
	if err != nil {
		return nil, err
	}
	beatmap.LastModificationTime = lastModificationTicks

	if version < 20140609 {
		arByte, err := readShort(r)
		if err != nil {
			return nil, err
		}
		arFloat := float32(arByte)
		beatmap.ApproachRate = arFloat

		csByte, err := readShort(r)
		if err != nil {
			return nil, err
		}
		csFloat := float32(csByte)
		beatmap.CircleSize = csFloat

		hpDrainByte, err := readShort(r)
		if err != nil {
			return nil, err
		}
		hpDrainFloat := float32(hpDrainByte)
		beatmap.HPDrain = hpDrainFloat

		odByte, err := readShort(r)
		if err != nil {
			return nil, err
		}
		odFloat := float32(odByte)
		beatmap.OverallDifficulty = odFloat
	} else {
		ar, err := readSingle(r)
		if err != nil {
			return nil, err
		}
		beatmap.ApproachRate = ar

		cs, err := readSingle(r)
		if err != nil {
			return nil, err
		}
		beatmap.CircleSize = cs

		hpDrain, err := readSingle(r)
		if err != nil {
			return nil, err
		}
		beatmap.HPDrain = hpDrain

		od, err := readSingle(r)
		if err != nil {
			return nil, err
		}
		beatmap.OverallDifficulty = od
	}

	sliderVelocity, err := readDouble(r)
	if err != nil {
		return nil, err
	}
	beatmap.SliderVelocity = sliderVelocity

	if version >= 20140609 && version < 20250107 {
		stdStars, err := readIntDoublePairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsStandardOld = stdStars

		taikoStars, err := readIntDoublePairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsTaikoOld = taikoStars

		ctbStars, err := readIntDoublePairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsCTBOld = ctbStars

		maniaStars, err := readIntDoublePairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsManiaOld = maniaStars

	} else if version >= 20250107 {
		stdStars, err := readIntFloatPairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsStandard = stdStars

		taikoStars, err := readIntFloatPairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsTaiko = taikoStars

		ctbStars, err := readIntFloatPairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsCTB = ctbStars

		maniaStars, err := readIntFloatPairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsMania = maniaStars
	}

	drainTime, err := readInt(r)
	if err != nil {
		return nil, err
	}
	beatmap.DrainTime = drainTime

	totalTime, err := readInt(r)
	if err != nil {
		return nil, err
	}
	beatmap.TotalTime = totalTime

	audioPreviewStartTime, err := readInt(r)
	if err != nil {
		return nil, err
	}
	beatmap.AudioPreviewStartTime = audioPreviewStartTime

	timingPoints, err := readTimingPoints(r)
	if err != nil {
		return nil, err
	}
	beatmap.TimingPoints = timingPoints

	difficultyID, err := readInt(r)
	if err != nil {
		return nil, err
	}
	beatmap.DifficultyID = difficultyID

	beatmapID, err := readInt(r)
	if err != nil {
		return nil, err
	}
	beatmap.BeatmapID = beatmapID

	threadID, err := readInt(r)
	if err != nil {
		return nil, err
	}
	beatmap.ThreadID = threadID

	if err := binary.Read(r, binary.LittleEndian, &beatmap.GradeStandard); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &beatmap.GradeTaiko); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &beatmap.GradeCTB); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &beatmap.GradeMania); err != nil {
		return nil, err
	}

	localOffset, err := readShort(r)
	if err != nil {
		return nil, err
	}
	beatmap.LocalBeatmapOffset = localOffset

	stackLeniency, err := readSingle(r)
	if err != nil {
		return nil, err
	}
	beatmap.StackLeniency = stackLeniency

	if err := binary.Read(r, binary.LittleEndian, &beatmap.GameplayMode); err != nil {
		return nil, err
	}

	songSource, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.SongSource = songSource

	songTags, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.SongTags = songTags

	onlineOffset, err := readShortSigned(r)
	if err != nil {
		return nil, err
	}
	beatmap.OnlineOffset = onlineOffset

	font, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.Font = font

	isUnplayed, err := readBoolean(r)
	if err != nil {
		return nil, err
	}
	beatmap.IsUnplayed = isUnplayed

	lastPlayed, err := readLong(r)
	if err != nil {
		return nil, err
	}
	beatmap.LastPlayed = lastPlayed

	isOsz2, err := readBoolean(r)
	if err != nil {
		return nil, err
	}
	beatmap.IsOsz2 = isOsz2

	folderName, err := readString(r)
	if err != nil {
		return nil, err
	}
	beatmap.FolderName = folderName

	lastChecked, err := readLong(r)
	if err != nil {
		return nil, err
	}
	beatmap.LastChecked = lastChecked

	ignoreBeatmapSound, err := readBoolean(r)
	if err != nil {
		return nil, err
	}
	beatmap.IgnoreBeatmapSound = ignoreBeatmapSound

	ignoreBeatmapSkin, err := readBoolean(r)
	if err != nil {
		return nil, err
	}
	beatmap.IgnoreBeatmapSkin = ignoreBeatmapSkin

	disableStoryboard, err := readBoolean(r)
	if err != nil {
		return nil, err
	}
	beatmap.DisableStoryboard = disableStoryboard

	disableVideo, err := readBoolean(r)
	if err != nil {
		return nil, err
	}
	beatmap.DisableVideo = disableVideo

	visualOverride, err := readBoolean(r)
	if err != nil {
		return nil, err
	}
	beatmap.VisualOverride = visualOverride

	if version < 20140609 {
		unknownShort, err := readShort(r)
		if err != nil {
			return nil, err
		}
		beatmap.UnknownShort = &unknownShort
	}

	lastModTime2, err := readInt(r)
	if err != nil {
		return nil, err
	}
	beatmap.LastModificationTime2 = lastModTime2

	if err := binary.Read(r, binary.LittleEndian, &beatmap.ManiaScrollSpeed); err != nil {
		return nil, err
	}

	return beatmap, nil
}

func readTimingPoints(r io.Reader) ([]TimingPoint, error) {
	count, err := readInt(r)
	if err != nil {
		return nil, err
	}

	timingPoints := make([]TimingPoint, 0, count)
	for i := 0; i < int(count); i++ {
		bpm, err := readDouble(r)
		if err != nil {
			return nil, err
		}

		offset, err := readDouble(r)
		if err != nil {
			return nil, err
		}

		inherited, err := readBoolean(r)
		if err != nil {
			return nil, err
		}

		timingPoints = append(timingPoints, TimingPoint{
			BPM:       bpm,
			Offset:    offset,
			Inherited: inherited,
		})
	}
	return timingPoints, nil
}

func ParseOsuDB(filename string) (*OsuDB, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 128*1024)

	version, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	folderCount, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	accountUnlocked, err := readBoolean(reader)
	if err != nil {
		return nil, err
	}

	var unlockDate time.Time
	ticks, err := readLong(reader)
	if err != nil {
		return nil, err
	}
	unlockDate = readDateTime(ticks)

	playerName, err := readString(reader)
	if err != nil {
		return nil, err
	}

	numberOfBeatmaps, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	beatmaps := make([]*Beatmap, 0, numberOfBeatmaps)
	for i := 0; i < int(numberOfBeatmaps); i++ {
		beatmap, err := readBeatmap(reader, version)
		if err != nil {
			return nil, err
		}
		beatmaps = append(beatmaps, beatmap)
	}

	userPermissions, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &OsuDB{
		Version:          version,
		FolderCount:      folderCount,
		AccountUnlocked:  accountUnlocked,
		UnlockDate:       unlockDate,
		PlayerName:       playerName,
		NumberOfBeatmaps: numberOfBeatmaps,
		Beatmaps:         beatmaps,
		UserPermissions:  userPermissions,
	}, nil
}

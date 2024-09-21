package osuParser

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"time"
)

type OsuDB struct {
	Version          int32
	FolderCount      int32
	AccountUnlocked  bool
	UnlockDate       time.Time
	PlayerName       string
	NumberOfBeatmaps int32
	Beatmaps         []*Beatmap
	UserPermissions  int32
}

type Beatmap struct {
	SizeInBytes           *int32
	Artist                string
	ArtistUnicode         string
	SongTitle             string
	SongTitleUnicode      string
	Creator               string
	Difficulty            string
	AudioFileName         string
	MD5Hash               string
	FileName              string
	RankedStatus          byte
	NumberOfHitCircles    uint16
	NumberOfSliders       uint16
	NumberOfSpinners      uint16
	LastModificationTime  int64
	ApproachRate          float32
	CircleSize            float32
	HPDrain               float32
	OverallDifficulty     float32
	SliderVelocity        float64
	StarRatingsStandard   map[int]int64
	StarRatingsTaiko      map[int]int64
	StarRatingsCTB        map[int]int64
	StarRatingsMania      map[int]int64
	DrainTime             int32
	TotalTime             int32
	AudioPreviewStartTime int32
	TimingPoints          []TimingPoint
	DifficultyID          int32
	BeatmapID             int32
	ThreadID              int32
	GradeStandard         byte
	GradeTaiko            byte
	GradeCTB              byte
	GradeMania            byte
	LocalBeatmapOffset    uint16
	StackLeniency         float32
	GameplayMode          byte
	SongSource            string
	SongTags              string
	OnlineOffset          int16
	Font                  string
	IsUnplayed            bool
	LastPlayed            int64
	IsOsz2                bool
	FolderName            string
	LastChecked           int64
	IgnoreBeatmapSound    bool
	IgnoreBeatmapSkin     bool
	DisableStoryboard     bool
	DisableVideo          bool
	VisualOverride        bool
	UnknownShort          *uint16
	LastModificationTime2 int32
	ManiaScrollSpeed      byte
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
	Name             string
	NumberOfBeatmaps int32
	Beatmaps         []*string
}

type Scores struct {
	Version        int32
	NumberOfScores int32
	Beatmaps       []*BeatmapScores
}

type BeatmapScores struct {
	BeatmapMD5Hash string
	NumberOfScores int32
	Scores         []*Score
}

type Score struct {
	Gamemode          byte
	Version           int32
	BeatmapMD5Hash    string
	PlayerName        string
	ReplayMD5Hash     string
	Count300s         uint16
	Count100s         uint16
	Count50           uint16
	Gekis             uint16
	Katus             uint16
	CountMiss         uint16
	ReplayScore       int32
	MaxCombo          uint16
	PerfectCombo      bool
	Mods              int32
	Timestamp         time.Time
	OnlineScoreId     int64
	AdditionalModInfo float64
}

func ParseCollectionsDB(filename string) (*Collections, error) {
	file, err := os.Open(filename)
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
	file, err := os.Open(filename)
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

	if version >= 20140609 {
		stdStars, err := readIntDoublePairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsStandard = stdStars

		taikoStars, err := readIntDoublePairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsTaiko = taikoStars

		ctbStars, err := readIntDoublePairs(r)
		if err != nil {
			return nil, err
		}
		beatmap.StarRatingsCTB = ctbStars

		maniaStars, err := readIntDoublePairs(r)
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
	file, err := os.Open(filename)
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

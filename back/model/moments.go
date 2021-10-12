package model

import (
    "errors"
    "regexp"
    "strconv"
    "strings"
    "time"
)

var (
    ErrEmptyVideoInputPath  = errors.New("empty VideoInputPath")
    ErrEmptyVideoOutputPath = errors.New("empty VideoOutputPath")
    ErrEmptyTextString      = errors.New("empty TextString")
    ErrFontColor            = errors.New("incorrect FontColor")
    ErrXY                   = errors.New("incorrect dimensions")
)

var validateHexaString = regexp.MustCompile(`^(0[xX])[A-Fa-f0-9]+$`)

type VideoEffects struct {
    // Video
    VideoInputPath  string
    VideoOutputPath string
    Duration        string
    Resolution      string

    // Effects
    TextString string
    XY         string
    Fontsize   string
    FontColor  string
    StartTime  string
    EndTime    string
}

type Video struct {
    VideoInputPath  string
    VideoOutputPath string
    Duration        time.Duration
    X, Y            int
}

func (v *Video) Parse(ve VideoEffects) (err error) {
    if ve.VideoInputPath == "" {
        return ErrEmptyVideoInputPath
    }
    v.VideoInputPath = ve.VideoInputPath

    if ve.VideoOutputPath == "" {
        return ErrEmptyVideoOutputPath
    }
    v.VideoOutputPath = ve.VideoOutputPath

    v.Duration, err = time.ParseDuration(strings.ReplaceAll(ve.Duration, " ", ""))
    if err != nil {
        return err
    }

    xy := strings.Split(strings.ReplaceAll(ve.Resolution, " ", ""), "x")
    if len(xy) == 0 {
        return ErrXY
    }

    v.X, err = strconv.Atoi(xy[0])
    if err != nil {
        return err
    }

    v.Y, err = strconv.Atoi(xy[1])
    if err != nil {
        return err
    }

    return nil
}

type Effects struct {
    TextString string
    X, Y       int
    Fontsize   int
    FontColor  string
    StartTime  time.Duration
    EndTime    time.Duration
}

func (e *Effects) Parse(ve VideoEffects) (err error) {
    if ve.TextString == "" {
        return ErrEmptyTextString
    }
    e.TextString = ve.TextString

    xy := strings.Split(strings.ReplaceAll(ve.XY, " ", ""), ",")
    if len(xy) == 0 {
        return ErrXY
    }

    e.X, err = strconv.Atoi(xy[0])
    if err != nil {
        return err
    }

    e.Y, err = strconv.Atoi(xy[1])
    if err != nil {
        return err
    }

    e.Fontsize, err = strconv.Atoi(ve.Fontsize)
    if err != nil {
        return err
    }

    if !validateHexaString.MatchString(ve.FontColor) {
        return ErrFontColor
    }
    e.FontColor = ve.FontColor

    e.StartTime, err = time.ParseDuration(strings.ReplaceAll(ve.StartTime, " ", ""))
    if err != nil {
        return err
    }

    e.EndTime, err = time.ParseDuration(strings.ReplaceAll(ve.EndTime, " ", ""))
    if err != nil {
        return err
    }

    return nil
}
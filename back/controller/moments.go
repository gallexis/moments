package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"steelseries/back/model"
)

func MomentsController(w http.ResponseWriter, r *http.Request) {
	var ve model.VideoEffects
	var video model.Video
	var effect model.Effects

	err := json.NewDecoder(r.Body).Decode(&ve)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = video.Parse(ve)
	if err != nil {
		fmt.Println("error parsing video: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = effect.Parse(ve)
	if err != nil {
		fmt.Println("error parsing effects: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = check(video, effect)
	if err != nil {
		fmt.Println("error check: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value, err := display(video, effect)
	if err != nil {
		fmt.Println("error display: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}


func check(v model.Video, e model.Effects) error {
	if e.X > v.X {
		return errors.New("error string: Invalid X coordinate")
	}

	if e.Y > v.Y {
		return errors.New("error string: Invalid Y coordinate")
	}

	if e.EndTime > v.Duration {
		return errors.New("error string: Invalid End Time")
	}

	if e.StartTime > e.EndTime {
		return errors.New("error string: Start Time is after End Time")
	}

	return nil
}

func display(v model.Video, e model.Effects) (string, error) {
	w := bytes.NewBufferString("")
	m := map[string]interface{}{
		"Input":      v.VideoInputPath,
		"Output":     v.VideoOutputPath,
		"StartAt":    e.StartTime.Seconds(),
		"EndAt":      e.EndTime.Seconds(),
		"TextString": e.TextString,
		"FontColor":  e.FontColor,
		"FontSize":   e.Fontsize,
		"X":          e.X,
		"Y":          e.Y,
	}
	ffmpegTemplate := `ffmpeg -i {{.Input}} -vf drawtext="enable='between(t,{{.StartAt}},{{.EndAt}})':text='{{.TextString}}':fontcolor={{.FontColor}}:fontsize={{.FontSize}}:x={{.X}}:y={{.Y}}" {{.Output}}`

	t := template.Must(template.New("").Parse(ffmpegTemplate))
	err := t.Execute(w, m)

	return w.String(), err
}
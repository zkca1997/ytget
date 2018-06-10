package main

import (
  "fmt"
  "log"
  "errors"
  "regexp"
  "strings"
  "net/url"
  "net/http"
  "io/ioutil"
)

func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoID(url)
  if err != nil { return err }

  err = y.getVideoInfo()
  if err != nil { return err }

  err = y.parseVideoInfo()
  if err != nil { return err }

  return nil
}

func (y *Youtube) parseVideoInfo() error {

	answer, err := url.ParseQuery(y.videoInfo)
	if err != nil {
		return err
	}

	status, ok := answer["status"]
	if !ok {
		err = errors.New("no response status found in the server's answer")
		return err
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return err
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return err
	}

	// read the streams map
	streamMap, ok := answer["url_encoded_fmt_stream_map"]
	if !ok {
		err = errors.New("no stream map found in the server's answer")
		return err
	}

	// read each stream
	streamsList := strings.Split(streamMap[0], ",")

	var streams []stream
	for streamPos, streamRaw := range streamsList {
		streamQry, err := url.ParseQuery(streamRaw)

		if err != nil {
			log.Printf("An error occured while decoding one of the video's stream's information: stream %d: %s\n", streamPos, err)
			continue
		}

		var sig string
		if _, exist := streamQry["sig"]; exist {
			sig = streamQry["sig"][0]
		}

		streams = append(streams, stream{
			"quality": streamQry["quality"][0],
			"type":    streamQry["type"][0],
			"url":     streamQry["url"][0],
			"sig":     sig,
			"title":   answer["title"][0],
			"author":  answer["author"][0],
		})
	}

	y.StreamList = streams
	return nil
}

func (y *Youtube) getVideoInfo() error {

	url := "http://youtube.com/get_video_info?video_id=" + y.VideoID

	resp, err := http.Get(url)
  defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	y.videoInfo = string(body)
	return nil
}

func (y *Youtube) findVideoID(url string) error {

	videoID := url   // videoID is inside the URL

  // if the url contains "youtu" or "?&/<%=", its a valid YouTube URL
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {

    // define locations of video IDs inside the URL
		reList := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}

    // search for video IDs inside the URL
		for _, re := range reList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}

  // if the ID contains the string "?&/<%=" its not valid
	if strings.ContainsAny(videoID, "?&/<%=") {
		return errors.New("invalid characters in video id")
	}

  // if the ID isn't at least 10 characters, its not valid
	if len(videoID) < 10 {
		return errors.New("the video id must be at least 10 characters long")
	}

  y.VideoID = videoID   // set the video ID inside the target struct

	return nil   // return no errors
}

package youtube

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"galched-bot/modules/settings"
)

const (
	YoutubeIDLength = 11
	youtubeRegexpID = `^.*((youtu.be\/)|(embed\/)|(watch\?))\??v?=?([^#\&\?\s]*).*`
)

var (
	urlRegex      = regexp.MustCompile(youtubeRegexpID)
	durationRegex = regexp.MustCompile(`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
)

type (
	Video struct {
		ID       string
		Title    string
		From     string
		Duration string

		Upvotes   uint64
		Downvotes uint64
		Views     uint64
	}

	Requester struct {
		mu *sync.RWMutex

		srv      *youtube.Service
		requests []Video
	}
)

func New(ctx context.Context, s *settings.Settings) (*Requester, error) {
	srv, err := youtube.NewService(ctx, option.WithAPIKey(s.YoutubeToken))
	if err != nil {
		return nil, err
	}

	return &Requester{
		mu:  new(sync.RWMutex),
		srv: srv,
	}, nil
}

func (r *Requester) AddVideo(query, from string) (string, error) {
	var (
		id  string
		err error
	)

	// try parse video id from the query
	id, err = videoID(query)
	if err != nil {
		// if we can't fo that, then search for the query
		resp, err := r.srv.Search.List("snippet").Type("video").MaxResults(1).Q(query).Do()
		if err != nil || len(resp.Items) == 0 || resp.Items[0].Id == nil {
			return "", fmt.Errorf("cannot parse youtube id: %w", err)
		}

		id = resp.Items[0].Id.VideoId
	}

	// get video info from api
	resp, err := r.srv.Videos.List("snippet,statistics,contentDetails").Id(id).Do()
	if err != nil {
		return "", fmt.Errorf("cannot send request to youtube api: %w", err)
	}

	// check if response have all required fields
	if len(resp.Items) == 0 {
		return "", errors.New("youtube api response does not contain items")
	}
	if resp.Items[0].Snippet == nil {
		return "", errors.New("youtube api response does not contain snippet")
	}
	if resp.Items[0].Statistics == nil {
		return "", errors.New("youtube api response does not contain statistics")
	}
	if resp.Items[0].ContentDetails == nil {
		return "", errors.New("youtube api response does not contain content details")
	}

	// check length of the video not more than 5 minutes
	if parseDuration(resp.Items[0].ContentDetails.Duration) == 0 {
		err = errors.New("видео не должно быть трансляцией")
		return err.Error(), err
	}

	// check video is not live
	if parseDuration(resp.Items[0].ContentDetails.Duration) > time.Minute*5 {
		err = errors.New("видео должно быть короче 5 минут")
		return err.Error(), err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// check if video already in the queue
	for i := range r.requests {
		if r.requests[i].ID == id {
			err = errors.New("видео уже есть в очереди")
			return err.Error(), err
		}
	}

	r.requests = append(r.requests, Video{
		ID:        id,
		From:      from,
		Duration:  strings.ToLower(resp.Items[0].ContentDetails.Duration[2:]),
		Title:     resp.Items[0].Snippet.Title,
		Upvotes:   resp.Items[0].Statistics.LikeCount,
		Views:     resp.Items[0].Statistics.ViewCount,
		Downvotes: resp.Items[0].Statistics.DislikeCount,
	})
	log.Printf("yt: added video < %s > from < %s >\n", resp.Items[0].Snippet.Title, from)

	return resp.Items[0].Snippet.Title, nil
}

func (r *Requester) List() []Video {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Video, len(r.requests))
	copy(result, r.requests)

	return result
}

func (r *Requester) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.requests {
		if r.requests[i].ID == id {
			r.requests = append(r.requests[:i], r.requests[i+1:]...)
			return
		}
	}
}

func videoID(url string) (string, error) {
	result := urlRegex.FindStringSubmatch(url)

	ln := len(result)
	if ln == 0 || len(result[ln-1]) != YoutubeIDLength {
		return "", fmt.Errorf("id haven't matched in \"%s\"", url)
	}

	return result[ln-1], nil
}

func parseDuration(str string) time.Duration {
	matches := durationRegex.FindStringSubmatch(str)

	years := parseInt64(matches[1])
	months := parseInt64(matches[2])
	days := parseInt64(matches[3])
	hours := parseInt64(matches[4])
	minutes := parseInt64(matches[5])
	seconds := parseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

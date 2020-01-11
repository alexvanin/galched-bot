package youtube

import (
	"io/ioutil"
	"testing"
)

func TestPlayground(t *testing.T) {
	youtubeTokenPath := "../../tokens/.youtubetoken"
	youtubetoken, err := ioutil.ReadFile(youtubeTokenPath)
	if err != nil {
		t.Errorf("cannot read youtube token: %v", err)
	}
	_ = youtubetoken
}

func TestFetchID(t *testing.T) {
	positiveCases := [][]string{
		{
			"https://www.youtube.com/watch?v=_v5IzvVTw7A&feature=feedrec_grec_index",
			"_v5IzvVTw7A",
		},
		{
			"https://www.youtube.com/watch?v=_v5IzvVTw7A#t=0m10s",
			"_v5IzvVTw7A",
		},
		{
			"https://www.youtube.com/embed/_v5IzvVTw7A?rel=0",
			"_v5IzvVTw7A",
		},
		{
			"https://www.youtube.com/watch?v=_v5IzvVTw7A    ", // multiple spaces in the end
			"_v5IzvVTw7A",
		},
		{
			"https://youtu.be/_v5IzvVTw7A",
			"_v5IzvVTw7A",
		},
		{
			"https://www.youtube.com/watch?v=PCp2iXA1uLE&list=PLvx4lPhqncyf10ymYz8Ph8EId0cafzhdZ&index=2&t=0s",
			"PCp2iXA1uLE",
		},
	}

	negativeCases := []string{
		"https://youtu.be/_vIzvVTw7A",   // short video
		"https://vimeo.com/_v5IzvVTw7A", // incorrect domain
		"youtube prime video",
	}

	for _, testCase := range positiveCases {
		res, err := videoID(testCase[0])
		if err != nil {
			t.Errorf("expecting error: <nil>, got: <%v>, url: %s\n", err, testCase[0])
		}
		if res != testCase[1] {
			t.Errorf("expecting result: %s, got: %s\n", testCase[0], res)
		}
	}

	for _, testCase := range negativeCases {
		_, err := videoID(testCase)
		if err == nil {
			t.Error("expecting error: <non nil>, got: <nil>")
		}
	}
}

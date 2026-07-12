package mediatags

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AppleMusicResponse struct {
	ResultCount int                `json:"resultCount"`
	Results     []AppleMusicResult `json:"results"`
}

type AppleMusicResult struct {
	ArtistName           *string    `json:"artistName,omitempty"`
	CollectionName       *string    `json:"collectionName,omitempty"`
	TrackName            *string    `json:"trackName,omitempty"`
	ArtworkURL100        *string    `json:"artworkUrl100,omitempty"`
	ReleaseDate          *time.Time `json:"releaseDate,omitempty"`
	PrimaryGenreName     *string    `json:"primaryGenreName,omitempty"`
	CollectionArtistName *string    `json:"collectionArtistName,omitempty"`
	DiscCount            *int       `json:"discCount,omitempty"`
	DiscNumber           *int       `json:"discNumber,omitempty"`
	TrackCount           *int       `json:"trackCount,omitempty"`
	TrackNumber          *int       `json:"trackNumber,omitempty"`
}

type ITunesMediaTagsProvider struct{}

var _ IProvider = ITunesMediaTagsProvider{}

func IntPtrToStringPtr(pi *int) *string {
	if pi == nil {
		return nil
	}

	return new(strconv.Itoa(*pi))
}

func (i ITunesMediaTagsProvider) FetchMediaTags(term string) (*MediaTags, error) {
	res, err := http.Get(fmt.Sprintf("https://itunes.apple.com/lookup?id=%s", term))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var x AppleMusicResponse
	err = json.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}

	if len(x.Results) <= 0 {
		return nil, fmt.Errorf("could not find tags for \"%s\"", term)
	}

	match := x.Results[0]
	mediaTags := MediaTags{
		AlbumArtist: match.ArtistName,
		Artist:      match.ArtistName,
		Album:       match.CollectionName,
		Title:       match.TrackName,
		Genre:       match.PrimaryGenreName,
		Track:       IntPtrToStringPtr(match.TrackNumber),
		TrackCount:  IntPtrToStringPtr(match.TrackCount),
		Disc:        IntPtrToStringPtr(match.DiscNumber),
		DiscCount:   IntPtrToStringPtr(match.DiscCount),
	}

	if match.CollectionArtistName != nil {
		mediaTags.AlbumArtist = match.CollectionArtistName
	}

	if match.ReleaseDate != nil {
		mediaTags.Year = new(strconv.Itoa(match.ReleaseDate.In(time.Local).Year()))
	}

	if match.ArtworkURL100 != nil {
		finalUrl := strings.Replace(*match.ArtworkURL100, "100x100bb", "1200x1200bb", 1)
		res, err := http.Get(finalUrl)
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = res.Body.Close()
		}()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		mediaTags.Cover = &MediaTagsCover{
			Data: data,
		}
	}

	return &mediaTags, nil
}

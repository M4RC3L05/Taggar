package mediatags

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"go.senan.xyz/taglib"
)

var numberAndCountRe = regexp.MustCompile(`^\s*(\d+)\s*(?:/\s*(\d+)\s*)?$`)

type TaglibMediaTagsRepository struct{}

var _ MediaTagsRepository = TaglibMediaTagsRepository{}

func parseNumberAndCount(s string) (num *string, count *string) {
	matches := numberAndCountRe.FindStringSubmatch(s)

	var nRes *string
	var cRes *string
	if len(matches) > 1 {
		n := matches[1]
		t := matches[2]

		if n != "" {
			nRes = &matches[1]
		}

		if t != "" {
			cRes = &matches[2]
		}
	}

	return nRes, cRes
}

func (t TaglibMediaTagsRepository) GetMediaTagsFromPath(path string) (*MediaTags, error) {
	tags, err := taglib.ReadTags(path)
	if err != nil {
		return nil, err
	}

	props, err := taglib.ReadProperties(path)
	if err != nil {
		return nil, err
	}

	mediaTags := MediaTags{}

	coverImgIndex := slices.IndexFunc(props.Images, func(img taglib.ImageDesc) bool {
		return strings.ToLower(img.Type) == "front cover"
	})

	if coverImgIndex == -1 && len(props.Images) > 0 {
		coverImgIndex = 0
	}

	if coverImgIndex != -1 {
		imgBytes, err := taglib.ReadImageOptions(path, coverImgIndex)
		if err != nil {
			return nil, err
		}

		mediaTags.Cover = &MediaTagsCover{
			Data:     imgBytes,
			Mimetype: props.Images[coverImgIndex].MIMEType,
		}
	}

	if val, ok := tags[taglib.AlbumArtist]; ok && len(val) > 0 {
		mediaTags.AlbumArtist = &val[0]
	}

	if val, ok := tags[taglib.Album]; ok && len(val) > 0 {
		mediaTags.Album = &val[0]
	}

	if val, ok := tags[taglib.Title]; ok && len(val) > 0 {
		mediaTags.Title = &val[0]
	}

	if val, ok := tags[taglib.Date]; ok && len(val) > 0 {
		mediaTags.Year = &val[0]
	}

	if val, ok := tags[taglib.Artist]; ok && len(val) > 0 {
		mediaTags.Artist = &val[0]
	}

	if val, ok := tags[taglib.Genre]; ok && len(val) > 0 {
		mediaTags.Genre = &val[0]
	}

	if val, ok := tags[taglib.TrackNumber]; ok && len(val) > 0 {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			n, c := parseNumberAndCount(val[0])

			mediaTags.Track = n
			mediaTags.TrackCount = c
		} else {
			mediaTags.Track = &val[0]
		}
	}

	if val, ok := tags["TRACKTOTAL"]; ok && len(val) > 0 {
		mediaTags.TrackCount = &val[0]
	}

	if val, ok := tags[taglib.DiscNumber]; ok && len(val) > 0 {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			n, c := parseNumberAndCount(val[0])

			mediaTags.Disc = n
			mediaTags.DiscCount = c
		} else {
			mediaTags.Disc = &val[0]
		}
	}

	if val, ok := tags["DISCTOTAL"]; ok && len(val) > 0 {
		mediaTags.DiscCount = &val[0]
	}

	return &mediaTags, nil
}

func (t TaglibMediaTagsRepository) SetMediaTagsFromPath(
	path string,
	tags *MediaTags,
) (*MediaTags, error) {
	if tags == nil {
		fmt.Println("No tags to set")
		return t.GetMediaTagsFromPath(path)
	}

	props, err := taglib.ReadProperties(path)
	if err != nil {
		return nil, err
	}

	tagsToSet := map[string][]string{}

	if tags.AlbumArtist != nil {
		tagsToSet[taglib.AlbumArtist] = []string{*tags.AlbumArtist}
	}

	if tags.Album != nil {
		tagsToSet[taglib.Album] = []string{*tags.Album}
	}

	if tags.Title != nil {
		tagsToSet[taglib.Title] = []string{*tags.Title}
	}

	if tags.Year != nil {
		tagsToSet[taglib.Date] = []string{*tags.Year}
	}

	if tags.Artist != nil {
		tagsToSet[taglib.Artist] = []string{*tags.Artist}
	}

	if tags.Genre != nil {
		tagsToSet[taglib.Genre] = []string{*tags.Genre}
	}

	if tags.Track != nil {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			final := ""
			final += *tags.Track

			if tags.TrackCount != nil {
				final += " / " + *tags.TrackCount
			}
			tagsToSet[taglib.TrackNumber] = []string{final}
		} else {
			tagsToSet[taglib.TrackNumber] = []string{*tags.Track}
		}
	}

	if tags.TrackCount != nil && (props.Format == "flac" || props.Format == "opus" ||
		props.Format == "ogg") {
		tagsToSet["TRACKTOTAL"] = []string{*tags.TrackCount}
	}

	if tags.Disc != nil {
		if props.Format != "flac" && props.Format != "opus" && props.Format != "ogg" {
			final := ""
			final += *tags.Disc

			if tags.DiscCount != nil {
				final += " / " + *tags.DiscCount
			}
			tagsToSet[taglib.DiscNumber] = []string{final}
		} else {
			tagsToSet[taglib.DiscNumber] = []string{*tags.Disc}
		}
	}

	if tags.DiscCount != nil && (props.Format == "flac" || props.Format == "opus" ||
		props.Format == "ogg") {
		tagsToSet["DISCTOTAL"] = []string{*tags.Disc}
	}

	err = taglib.WriteTags(path, tagsToSet, 0)
	if err != nil {
		return nil, err
	}

	if tags.Cover != nil {
		err = taglib.WriteImage(path, tags.Cover.Data)
		if err != nil {
			return nil, err
		}
	}

	return t.GetMediaTagsFromPath(path)
}

package mediatags

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

type CobraCmdMediaTagsProvider struct {
	Cmd *cobra.Command
}

var _ IProvider = ITunesMediaTagsProvider{}

func Set(cmd *cobra.Command, n string, dest *(*string)) error {
	if cmd.Flags().Changed(n) {
		val, err := cmd.Flags().GetString(n)
		if err != nil {
			return err
		}

		*dest = &val
	}

	return nil
}

func (i CobraCmdMediaTagsProvider) FetchMediaTags() (*MediaTags, error) {
	res := MediaTags{}

	if err := Set(i.Cmd, "title", &res.Title); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "artist", &res.Artist); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "album", &res.Album); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "albumArtist", &res.AlbumArtist); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "genre", &res.Genre); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "year", &res.Year); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "track", &res.Track); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "trackCount", &res.TrackCount); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "dist", &res.Disc); err != nil {
		return nil, err
	}

	if err := Set(i.Cmd, "distCount", &res.DiscCount); err != nil {
		return nil, err
	}

	if i.Cmd.Flags().Changed("cover") {
		val, err := i.Cmd.Flags().GetString("cover")
		if err != nil {
			return nil, err
		}

		f, err := os.Open(val)
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = f.Close()
		}()

		data, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		res.Cover = &MediaTagsCover{
			Data: data,
		}
	}

	return &res, nil
}

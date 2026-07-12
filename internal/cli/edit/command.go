package edit

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/m4rc3l05/taggar/internal"
	mediatags "github.com/m4rc3l05/taggar/internal/media_tags"
	"github.com/spf13/cobra"
)

type EditFlags struct {
	Path        *string
	Provider    *string
	Term        *string
	Cover       *string
	Title       *string
	Artist      *string
	Album       *string
	AlbumArtist *string
	Genre       *string
	Year        *string
	Track       *string
	TrackCount  *string
	Disc        *string
	DiscCount   *string
}

func Set(cmd *cobra.Command, n string, dest *(*string)) error {
	if cmd.Flags().Changed(n) {
		path, err := cmd.Flags().GetString(n)
		if err != nil {
			return err
		}

		*dest = &path
	}

	return nil
}

type cmd struct {
	data EditFlags
}

func (c *cmd) Pre(cmd *cobra.Command) error {
	x := EditFlags{}

	if err := Set(cmd, "path", &x.Path); err != nil {
		return err
	}

	if err := Set(cmd, "provider", &x.Provider); err != nil {
		return err
	}
	if err := Set(cmd, "term", &x.Term); err != nil {
		return err
	}

	if err := Set(cmd, "cover", &x.Cover); err != nil {
		return err
	}
	if err := Set(cmd, "title", &x.Title); err != nil {
		return err
	}
	if err := Set(cmd, "artist", &x.Artist); err != nil {
		return err
	}
	if err := Set(cmd, "album", &x.Album); err != nil {
		return err
	}
	if err := Set(cmd, "albumArtist", &x.AlbumArtist); err != nil {
		return err
	}
	if err := Set(cmd, "genre", &x.Genre); err != nil {
		return err
	}
	if err := Set(cmd, "year", &x.Year); err != nil {
		return err
	}
	if err := Set(cmd, "track", &x.Track); err != nil {
		return err
	}
	if err := Set(cmd, "trackCount", &x.TrackCount); err != nil {
		return err
	}
	if err := Set(cmd, "dist", &x.Disc); err != nil {
		return err
	}
	if err := Set(cmd, "distCount", &x.DiscCount); err != nil {
		return err
	}

	c.data = x

	return nil
}

func (c cmd) Run(cmd *cobra.Command) error {
	m := &mediatags.MediaTags{}

	if c.data.Provider != nil {
		switch *c.data.Provider {
		case "itunes":
			{
				if c.data.Term == nil {
					return errors.New("term must be provider for itunes provider")
				}

				fmt.Println("Fetching metadata from itunes")
				res, err := mediatags.ITunesMediaTagsProvider{}.FetchMediaTags(*c.data.Term)
				if err != nil {
					return err
				}

				m.CopyFrom(res)
			}
		default:
			{
				return errors.New("provider not supported")
			}
		}
	}

	if c.data.Cover != nil {
		f, err := os.Open(*c.data.Cover)
		if err != nil {
			return err
		}

		defer func() {
			_ = f.Close()
		}()

		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		m.Cover = &mediatags.MediaTagsCover{
			Data: data,
		}
	}

	m.CopyFrom(&mediatags.MediaTags{
		AlbumArtist: c.data.AlbumArtist,
		Album:       c.data.Album,
		Title:       c.data.Title,
		Year:        c.data.Year,
		Artist:      c.data.Artist,
		Genre:       c.data.Genre,
		Track:       c.data.Track,
		TrackCount:  c.data.TrackCount,
		Disc:        c.data.Disc,
		DiscCount:   c.data.DiscCount,
	})

	fmt.Println("Persisting tags")
	tags, err := mediatags.TaglibMediaTagsRepository{}.SetMediaTagsFromPath(*c.data.Path, m)
	if err != nil {
		return err
	}

	return mediatags.DisplayMediaTags(*tags)
}

func NewCommand() *cobra.Command {
	c := cmd{}
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit audio tags",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.Pre(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(cmd)
		},
	}

	editCmd.Flags().StringP("path", "p", "", "the path to the audio file to view metadata")

	editCmd.Flags().
		StringP("provider", "s", "", "the provider to use to prefill metadata\navailable providers:\n\t> itunes")
	editCmd.Flags().
		StringP("term", "t", "", "the search term the provider will use to fetch metadata")

	editCmd.Flags().String("cover", "", "the cover")
	editCmd.Flags().String("title", "", "the title")
	editCmd.Flags().String("artist", "", "the artist")
	editCmd.Flags().String("album", "", "the album")
	editCmd.Flags().String("albumArtist", "", "the album artist")
	editCmd.Flags().String("genre", "", "the genre")
	editCmd.Flags().String("year", "", "the year")
	editCmd.Flags().String("track", "", "the track")
	editCmd.Flags().String("trackCount", "", "the track count")
	editCmd.Flags().String("disc", "", "the disc")
	editCmd.Flags().String("discCount", "", "the disc count")

	internal.Must(editCmd.MarkFlagRequired("path"))
	internal.Must(editCmd.MarkFlagFilename("path"))
	editCmd.MarkFlagsRequiredTogether("provider", "term")

	return editCmd
}

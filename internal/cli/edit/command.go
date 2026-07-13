package edit

import (
	"errors"
	"fmt"

	"github.com/m4rc3l05/taggar/internal"
	mediatags "github.com/m4rc3l05/taggar/internal/media_tags"
	"github.com/spf13/cobra"
)

func getFlag(cmd *cobra.Command, name string) (string, error) {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", err
	}

	return val, nil
}

func NewCommand() *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit audio tags",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var m *mediatags.MediaTags
			path, err := getFlag(cmd, "path")
			if err != nil {
				return err
			}

			provider, err := getFlag(cmd, "provider")
			if err != nil {
				return err
			}

			if cmd.Flags().Changed("provider") {
				term, err := getFlag(cmd, "term")
				if err != nil {
					return err
				}

				switch provider {
				case "itunes":
					{
						fmt.Println("Fetching metadata from itunes")
						res, err := mediatags.ITunesMediaTagsProvider{Id: term}.FetchMediaTags()
						if err != nil {
							return err
						}

						m = res
					}
				default:
					{
						return errors.New("provider not supported")
					}
				}
			}

			res, err := mediatags.CobraCmdMediaTagsProvider{Cmd: cmd}.FetchMediaTags()

			// If setting tags manually (not via a provider) return the error
			if err != nil && !cmd.Flags().Changed("provider") {
				return err
			}

			if m != nil {
				m.CopyFrom(res)
			} else {
				m = res
			}

			fmt.Println("Persisting tags")
			tags, err := mediatags.TaglibMediaTagsRepository{}.SetMediaTagsFromPath(path, m)
			if err != nil {
				return err
			}

			return mediatags.DisplayMediaTags(*tags)
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

package mediatags

type MediaTagsCover struct {
	Data     []byte
	Mimetype string
}

type MediaTags struct {
	Cover       *MediaTagsCover
	AlbumArtist *string
	Album       *string
	Title       *string
	Year        *string
	Artist      *string
	Genre       *string
	Track       *string
	TrackCount  *string
	Disc        *string
	DiscCount   *string
}

func (m *MediaTags) CopyFrom(mm *MediaTags) {
	if mm == nil {
		return
	}

	if mm.Album != nil {
		m.Album = mm.Album
	}

	if mm.AlbumArtist != nil {
		m.AlbumArtist = mm.AlbumArtist
	}

	if mm.Artist != nil {
		m.Artist = mm.Artist
	}

	if mm.Cover != nil {
		m.Cover = mm.Cover
	}

	if mm.Disc != nil {
		m.Disc = mm.Disc
	}

	if mm.DiscCount != nil {
		m.DiscCount = mm.DiscCount
	}

	if mm.Genre != nil {
		m.Genre = mm.Genre
	}

	if mm.Title != nil {
		m.Title = mm.Title
	}

	if mm.Track != nil {
		m.Track = mm.Track
	}

	if mm.TrackCount != nil {
		m.TrackCount = mm.TrackCount
	}

	if mm.Year != nil {
		m.Year = mm.Year
	}
}

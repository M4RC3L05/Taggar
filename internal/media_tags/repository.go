package mediatags

type MediaTagsRepository interface {
	GetMediaTagsFromPath(path string) (*MediaTags, error)
	SetMediaTagsFromPath(path string, tags *MediaTags) (*MediaTags, error)
}

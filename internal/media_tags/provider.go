package mediatags

type IProvider interface {
	FetchMediaTags() (*MediaTags, error)
}

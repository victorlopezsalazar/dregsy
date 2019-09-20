package tagsupport

func NumberOfTags(tags []string) int {
	if len(tags) > 11 {
		return 11
	}
	return len(tags)
}
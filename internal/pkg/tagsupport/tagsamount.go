package tagsupport

func NumberOfTags(tags []string) int {
	if len(tags) > 12 {
		return 12
	}
	return len(tags)
}
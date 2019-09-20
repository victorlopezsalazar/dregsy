package tagsupport

import "regexp"


func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func FilterReleaseTags(tags []string) []string {
	return filter(tags, func(tag string) bool {
		match, _ := regexp.MatchString("^(master|latest|\\d+)$", tag)
		return match
	})
}
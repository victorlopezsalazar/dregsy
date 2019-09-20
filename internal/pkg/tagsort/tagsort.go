package tagsort

import "sort"

type TagSlice []string

func (p TagSlice) Len() int           { return len(p) }
func (p TagSlice) Less(i, j int) bool {
	relativeSize := len(p[i]) - len(p[j])
	if relativeSize == 0{
		return p[i] < p[j]
	}
	return relativeSize < 0
}
func (p TagSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }


func Sort(data []string){
	sort.Sort(TagSlice(data))
}

package model

type Dir struct {
	Key   string
	Value int
}

type DirList []Dir

func (p DirList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p DirList) Len() int           { return len(p) }
func (p DirList) Less(i, j int) bool { return p[i].Value < p[j].Value }

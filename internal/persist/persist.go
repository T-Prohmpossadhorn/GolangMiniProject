package persist

type Persist interface {
	Getfulllist() []string
	Getfirstfilter(filter string) string
}

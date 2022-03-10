package scraper

type Result struct {
	PageURL            string
	InternalLinksCount uint
	ExternalLinksCount uint
	Success            bool
	Error              error
}

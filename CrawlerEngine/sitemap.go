package CrawlerEngine

import "os"

type SitemapStreamer interface {
	Init()
	End()
	NewLoc(link string)
	AppendChildLink(childLink string)
}

func NewSitemapStreamer(filePath string) *sitemapStreamer {
	return &sitemapStreamer{
		filePath: filePath,
	}
}

type sitemapStreamer struct {
	filePath string
	f        *os.File
}

func (s *sitemapStreamer) Init() {
	var err error
	if s.filePath == "" {
		return
	}

	s.f, err = os.Create(s.filePath)
	if err != nil {
		logger.Errorf("Unable to create on %s", s.filePath)
		panic(err)
	}
}

func (s *sitemapStreamer) End() {
	if s.f != nil {
		s.f.Close()
	}
}

func (s *sitemapStreamer) NewLoc(link string) {
	if s.f != nil {
		line := []byte("\n" + link + "\n")
		s.f.Write(line)
	}

	s.f.Sync()
}

func (s *sitemapStreamer) AppendChildLink(childLink string) {
	if s.f != nil {
		line := []byte("\t- " + childLink + "\n")
		s.f.Write(line)
	}
}

package export

import (
	"encoding/csv"
	"io"

	"github.com/stjudewashere/seonaut/internal/crawler"
)

type Link struct {
	Origin      string
	Destination string
	Text        string
}

type Image struct {
	Origin string
	Image  string
	Alt    string
}

type Script struct {
	Origin string
	Script string
}

type Style struct {
	Origin string
	Style  string
}

type Store interface {
	ExportLinks(*crawler.Crawl) <-chan *Link
	ExportExternalLinks(*crawler.Crawl) <-chan *Link
	ExportImages(crawl *crawler.Crawl) <-chan *Image
	ExportScripts(crawl *crawler.Crawl) <-chan *Script
	ExportStyles(crawl *crawler.Crawl) <-chan *Style
}

type Exporter struct {
	store Store
}

func NewExporter(s Store) *Exporter {
	return &Exporter{
		store: s,
	}
}

// Export internal links as a CSV file
func (e *Exporter) ExportLinks(f io.Writer, crawl *crawler.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Destination",
		"Text",
	})

	lStream := e.store.ExportLinks(crawl)

	for v := range lStream {
		w.Write([]string{
			v.Origin,
			v.Destination,
			v.Text,
		})
	}

	w.Flush()
}

// Export internal links as a CSV file
func (e *Exporter) ExportExternalLinks(f io.Writer, crawl *crawler.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Destination",
		"Text",
	})

	lStream := e.store.ExportExternalLinks(crawl)

	for v := range lStream {
		w.Write([]string{
			v.Origin,
			v.Destination,
			v.Text,
		})
	}

	w.Flush()
}

// Export all images as a CSV file
func (e *Exporter) ExportImages(f io.Writer, crawl *crawler.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Image URL",
		"Alt",
	})

	iStream := e.store.ExportImages(crawl)

	for v := range iStream {
		w.Write([]string{
			v.Origin,
			v.Image,
			v.Alt,
		})
	}

	w.Flush()
}

// Export all scripts as a CSV file
func (e *Exporter) ExportScripts(f io.Writer, crawl *crawler.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Script URL",
	})

	sStream := e.store.ExportScripts(crawl)

	for v := range sStream {
		w.Write([]string{
			v.Origin,
			v.Script,
		})
	}

	w.Flush()
}

// Export all CSS styles as a CSV file
func (e *Exporter) ExportStyles(f io.Writer, crawl *crawler.Crawl) {
	w := csv.NewWriter(f)

	w.Write([]string{
		"Origin",
		"Style URL",
	})

	sStream := e.store.ExportStyles(crawl)

	for v := range sStream {
		w.Write([]string{
			v.Origin,
			v.Style,
		})
	}

	w.Flush()
}
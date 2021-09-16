package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func handler(w http.ResponseWriter, req *http.Request) {
	// Setup new logger
	zapconfig := zap.NewProductionConfig()
	zapconfig.Encoding = "console"
	zapconfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := zapconfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync() //nolint:errcheck

	// Fetch rss feed
	span, _ := tracer.StartSpanFromContext(req.Context(), "http.get")
	span.SetTag("http.url", "https://tapas.io/rss/series/3346")
	resp, err := http.Get("https://tapas.io/rss/series/3346")
	span.Finish(tracer.WithError(err))
	if err != nil {
		logger.Error("An error occurred while fetching rss feed", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// Read response body
	span, _ = tracer.StartSpanFromContext(req.Context(), "ioutil.readall")
	body, err := ioutil.ReadAll(resp.Body)
	span.Finish(tracer.WithError(err))
	if err != nil {
		logger.Error("An error occurred while reading rss feed", zap.Error(err))
		return
	}

	// Convert the body to type string
	span, _ = tracer.StartSpanFromContext(req.Context(), "fp.parsestring")
	sb := string(body)
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(sb)
	span.Finish(tracer.WithError(err))
	if err != nil {
		logger.Error("An error occurred while parsing rss feed", zap.Error(err))
		return
	}

	// Rebuild feed
	ffeed := &feeds.Feed{
		Title:       feed.Title,
		Description: feed.Description,
		Link:        &feeds.Link{Href: "https://tapas.io/rss/series/3346"},
		Updated:     *feed.UpdatedParsed,
	}
	for _, item := range feed.Items {
		span, _ = tracer.StartSpanFromContext(req.Context(), "htmlquery.loadurl")
		span.SetTag("http.url", item.Link)
		doc, err := htmlquery.LoadURL(item.Link)
		span.Finish(tracer.WithError(err))
		if err != nil {
			logger.Error("Error building doc", zap.Error(err))
			continue
		}

		// Some magic here to extract the actual comic png
		// The xpath query matches both the png and the gif
		span, _ = tracer.StartSpanFromContext(req.Context(), "htmlquery.queryall")
		nodes, err := htmlquery.QueryAll(doc, "//img[@data-series-id='3346']/@data-src")
		span.Finish(tracer.WithError(err))
		if err != nil {
			logger.Error("Error building nodes", zap.Error(err))
			continue
		}

		ffeed.Items = append(ffeed.Items, &feeds.Item{
			Title:   fmt.Sprintf("%s.", item.Title),
			Created: *item.PublishedParsed,
			// Updated: *item.UpdatedParsed, // This ends up in a nil pointer dereference
			Link: &feeds.Link{
				Href: item.Link,
			},
			Author: &feeds.Author{
				Name:  item.Author.Name,
				Email: item.Author.Email,
			},
			Content: nodes[0].FirstChild.Data,
			// Description: item.Content, // The item content needs to be sanitized first
		})
	}

	span, _ = tracer.StartSpanFromContext(req.Context(), "ffeed.torss")
	sfeed, err := ffeed.ToRss()
	span.Finish(tracer.WithError(err))
	if err != nil {
		logger.Error("An error occurred while generating rss feed", zap.Error(err))
		return
	}

	// For some reasons, strings in ffeed are escaped in sfeed
	// For instance "World's Sweatiest Comic" becomes "World&#39;s Sweatiest Comic"
	// Ideally we would just unescape them, unescaping the whole xml string is good enough
	span, _ = tracer.StartSpanFromContext(req.Context(), "html.unescapestring")
	ufeed := html.UnescapeString(sfeed)
	span.Finish()

	_, err = w.Write([]byte(ufeed))
	if err != nil {
		logger.Error("An error occurred while writing the result", zap.Error(err))
		return
	}
}

func main() {
	tracer.Start(
		tracer.WithService("mrlovenstein-rss"),
		// tracer.WithEnv("env"),
		tracer.WithUDS("/var/run/datadog/apm.sock"),
		tracer.WithSamplingRules([]tracer.SamplingRule{tracer.RateRule(1)}),
	)

	if err := profiler.Start(
		profiler.WithService("mrlovenstein-rss"),
		// profiler.WithEnv("env"),
		profiler.WithUDS("/var/run/datadog/apm.sock"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,

			// The profiles below are disabled by
			// default to keep overhead low, but
			// can be enabled as needed.
			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
		),
	); err != nil {
		tracer.Stop() // defer won't run on panic, so duplicate the tracer stop here
		log.Fatal(err)
	}
	defer tracer.Stop()
	defer profiler.Stop()

	// Create a traced mux router
	mux := httptrace.NewServeMux()
	// Continue using the router as you normally would
	mux.HandleFunc("/mrlovenstein.xml", handler)
	log.Print(http.ListenAndServe(":8080", mux)) // log would print the exit reason for the HTTP server
}

package medad

import (
	"context"
	"html/template"
	"io"

	"github.com/janstoon/toolbox/bricks"
)

type (
	// LangCode is ISO 639-1 two-letter code
	LangCode string

	Language struct {
		Title string
		Rtl   bool
	}
)

var (
	LangCodeEnglish LangCode = "en"
	LangCodePersian LangCode = "fa"

	Languages = map[LangCode]Language{
		LangCodeEnglish: {
			Title: "English",
			Rtl:   false,
		},
		LangCodePersian: {
			Title: "فارسی",
			Rtl:   true,
		},
	}
)

// Article contains a multilingual content
type Article struct {
	Name    string
	Content map[LangCode]io.ReadCloser
}

// Loader loads articles from sources
type Loader interface {
	Load(ctx context.Context, sources ...string) bricks.Bag[Article]
}

// Compiler converts articles and writes output (files) into destination (directory)
type Compiler interface {
	Compile(ctx context.Context, tpl *template.Template, destination string, sources bricks.Bag[Article]) error
}

type Uploader interface {
	Upload(ctx context.Context, local, remote string) error
}

type Settings struct {
	ArticlesGlob  string
	TemplatesGlob string
	DistDirectory string

	FtpHost         string
	FtpPort         string
	FtpUsername     string
	FtpPassword     string
	RemoteDirectory string
}

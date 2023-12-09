package compilers

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"github.com/russross/blackfriday/v2"
	"github.com/yuin/goldmark"
	"golang.org/x/sync/errgroup"

	"github.com/pattack/medad"
)

type markdownCompiler struct {
	RootTemplate string
	RootFilename string

	ArticleTemplate string
	ArticlesDir     string
}

func MarkdownCompiler(opts ...tricks.Option[markdownCompiler]) medad.Compiler {
	c := markdownCompiler{
		RootTemplate: "root.gohtml",
		RootFilename: "index.html",

		ArticleTemplate: "article.gohtml",
		ArticlesDir:     "articles",
	}
	tricks.ApplyOptions(&c, opts...)

	return c
}

func (c markdownCompiler) Compile(
	ctx context.Context, tpl *template.Template, destination string, sources bricks.Bag[medad.Article],
) error {
	err := os.MkdirAll(destination, 0755)
	if nil != err {
		return err
	}

	aa, err := c.compileArticles(ctx, tpl, destination, sources)
	if nil != err {
		return err
	}

	err = c.compileIndex(ctx, tpl, destination, aa...)
	if nil != err {
		return err
	}

	return nil
}

func (c markdownCompiler) compileArticles(
	ctx context.Context, tpl *template.Template, destination string, sources bricks.Bag[medad.Article],
) ([]article, error) {
	dstArticles := filepath.Join(destination, c.ArticlesDir)
	err := os.MkdirAll(dstArticles, 0755)
	if nil != err {
		return nil, err
	}

	aa := make([]article, 0)
	caa := make(chan article, 10)
	wgCollect := sync.WaitGroup{}
	wgCollect.Add(1)
	go func(ctx context.Context) {
		defer wgCollect.Done()

		for {
			select {
			case a, ok := <-caa:
				if !ok {
					return
				}

				aa = append(aa, a)

			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	wgCompile := errgroup.Group{}
	for src, err := sources.Pull(ctx); err != nil; src, err = sources.Pull(ctx) {
		func(src *medad.Article) {
			for langCode := range src.Content {
				wgCompile.Go(func() error {
					a, err := c.compileArticle(ctx, tpl, dstArticles, langCode, src)
					if err != nil {
						return err
					}

					caa <- tricks.PtrVal(a)

					return nil
				})
			}
		}(src)
	}

	if err := wgCompile.Wait(); err != nil {
		close(caa)

		return nil, err
	}

	close(caa)
	wgCollect.Wait()

	return aa, nil
}

func (c markdownCompiler) compileArticle(
	ctx context.Context, tpl *template.Template, destination string, langCode medad.LangCode, src *medad.Article,
) (*article, error) {
	bb, err := io.ReadAll(src.Content[langCode])
	if err != nil {
		return nil, err
	}

	an := src.Name
	bn := fmt.Sprintf("%s.%s.html", an, langCode)

	doc := blackfriday.New().Parse(bb)
	title := getTitle(doc, an)

	fn := filepath.Join(destination, bn)
	fh, err := os.Create(fn)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fh.Close() }()

	md := goldmark.New()
	err = md.Convert(bb, fh)
	if err != nil {
		return nil, err
	}

	return &article{
		Title:    title,
		LangCode: langCode,
		Link:     filepath.Join(c.ArticlesDir, bn),
	}, nil
}

func (c markdownCompiler) compileIndex(
	ctx context.Context, tpl *template.Template, destination string, articles ...article,
) error {
	fh, err := os.Create(filepath.Join(destination, c.RootFilename))
	if nil != err {
		return err
	}
	defer func() { _ = fh.Close() }()

	err = tpl.ExecuteTemplate(fh, c.RootTemplate, newArgs(articles...))
	if nil != err {
		return err
	}

	return nil
}

func extensionless(basename string) string {
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

func getTitle(doc *blackfriday.Node, fallback string) string {
	heading := doc.FirstChild
	for heading != nil && heading.Level != 1 {
		heading = heading.Next
	}

	if heading == nil || heading.FirstChild == nil {
		return fallback
	}

	return string(heading.FirstChild.Literal)
}

type args struct {
	Header    header
	Defaults  defaults
	Languages map[medad.LangCode]medad.Language
	Articles  []article
}

func newArgs(articles ...article) args {
	scope := args{
		Header: newHeader(),
		Defaults: defaults{
			LangCode: medad.LangCodeEnglish,
		},
		Languages: make(map[medad.LangCode]medad.Language),
		Articles:  articles,
	}

	for _, a := range articles {
		scope.Languages[a.LangCode] = medad.Languages[a.LangCode]
	}

	return scope
}

type header struct {
	Date time.Time

	OS   string
	Arch string

	PackageName string
	Filename    string
	LineNumber  string
}

func newHeader() header {
	return header{
		Date: time.Now(),

		OS:   os.Getenv("GOOS"),
		Arch: os.Getenv("GOARCH"),

		PackageName: os.Getenv("GOPACKAGE"),
		Filename:    os.Getenv("GOFILE"),
		LineNumber:  os.Getenv("GOLINE"),
	}
}

type defaults struct {
	LangCode medad.LangCode
}

type article struct {
	Title    string
	LangCode medad.LangCode
	Link     string
}

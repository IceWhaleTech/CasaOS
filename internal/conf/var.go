package conf

import "regexp"

var (
	BuiltAt    string
	GoVersion  string
	GitAuthor  string
	GitCommit  string
	Version    string = "dev"
	WebVersion string
)

var (
	Conf *Config
)

var SlicesMap = make(map[string][]string)
var FilenameCharMap = make(map[string]string)
var PrivacyReg []*regexp.Regexp

var (
	// StoragesLoaded loaded success if empty
	StoragesLoaded = false
)
var (
	RawIndexHtml string
	ManageHtml   string
	IndexHtml    string
)

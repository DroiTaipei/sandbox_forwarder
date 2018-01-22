package file

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	fileBasenameLength = 32
)

// GetCdnURL - Generate CDN url
// Format: http://{cdn_fqdn}/droi/{appid}/{fid}/{basename}
// e.g. http://newmarket1.oo523.com/droi/nf8umbzhxE0XHniHvyVusfK6LXPwathVJqZgaDUC/796239408339750912/foo.jpg
func GetCdnURL(fqdn, appid, fid, filename string) string {
	return fmt.Sprintf("http://%s%s/%s/%s/%s",
		fqdn,
		CDNRootPath,
		appid,
		fid,
		GetFileBasename(fid, filename))
}

// ParseCdnURL - Get fqdn, qppid, fid, and filename from CDN url
func ParseCdnURL(cdnURL string) (fqdn, appid, fid, filename string, err error) {
	urlObj, err := url.Parse(cdnURL)
	if err != nil {
		return "", "", "", "", err
	}
	pathTokens := strings.Split(urlObj.Path, "/")
	if len(pathTokens) != 5 { // e.g. /droi/nf8umbzhxE0XHniHvyVusfK6LXPwathVJqZgaDUC/796239408339750912/foo.jpg
		return "", "", "", "", fmt.Errorf("CDN url path is invalid")
	}
	return urlObj.Host, pathTokens[2], pathTokens[3], pathTokens[4], nil
}

// GetFileBasename - Get raw file basename (may be in Chinese), work on Windows, Linux file system and URL
func GetFileBasename(fid, in string) string {
	var fileBasename string
	ext := filepath.Ext(in)
	// not include only `.`
	if len(ext) > 1 {
		fileBasename = fmt.Sprintf("%s%s", fid, ext)
	} else {
		fileBasename = fid
	}
	// Only accept utf8
	if !utf8.Valid([]byte(in)) {
		return fileBasename
	}
	// Handle Windows filepath
	if strings.Contains(in, "\\") {
		tokens := strings.Split(in, "\\")
		in = tokens[len(tokens)-1]
	}
	// get clean basename
	basename := filepath.Clean(filepath.Base(in))
	// truncate space
	basename = strings.Join(strings.Fields(basename), "")
	// truncate multibyte
	scan := len(basename)
	for scan > 0 {
		_, size := utf8.DecodeLastRuneInString(basename[:scan])
		// remove multibyte char
		if size > 1 {
			basename = basename[:scan-size] + basename[scan:]
		}
		scan -= size
	}
	// truncate length of basename
	if utf8.RuneCountInString(basename) > fileBasenameLength {
		numRunes := 0
		for idx := range basename {
			numRunes++
			if numRunes > utf8.RuneCountInString(basename)-fileBasenameLength {
				basename = basename[idx:]
				break
			}
		}
	}
	// encode path
	basename = url.PathEscape(basename)
	// Replace url sensitive chars
	replacer := strings.NewReplacer(
		"#", "-",
		"?", "-",
		"%", "", // remove %
		"+", "-",
		"=", "-",
		"@", "-",
		"$", "-",
		"&", "-",
	)
	basename = replacer.Replace(basename)
	if len(basename) == 0 || strings.HasPrefix(basename, ".") {
		return fileBasename
	}
	return basename
}

// IsGosunURL - if the url hosted by Gosun
func IsGosunURL(url string) bool {
	return strings.Contains(url, VendorGosunDomain)
}

// IsQiniuURL - if the url hosted by Qiniu
func IsQiniuURL(url string) bool {
	return strings.Contains(url, VendorQiniuDomain)
}

package client

import (
	urlpkg "net/url"
	"path"
)

func Headers(url *urlpkg.URL, headers map[string]string) (map[string]string, error) {
	var defaultHeaders = map[string]string{
		"Host":            url.Host,
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
		"Accept":          headerAccept(url),
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding": "gzip, deflate, br",
		"Referer":         "",
		"Connection":      "keep-alive",
		// "Cookie": "",
		"Upgrade-Insecure-Requests": "1",
		"Pragma":                    "no-cache",
		"Cache-Control":             "no-cache",
	}

	if headers != nil {
		for k, v := range headers {
			defaultHeaders[k] = v
		}
	}
	return defaultHeaders, nil
}

func headerAccept(url *urlpkg.URL) string {
	defaultAccept := "*/*"
	htmlAccept := "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	imageAccept := "image/webp,*/*"
	cssAccept := "text/css,*/*;q=0.1"
	jsAccept := defaultAccept
	extAcceptMap := map[string]string{
		".gif":  imageAccept,
		".png":  imageAccept,
		".jpg":  imageAccept,
		".jpeg": imageAccept,
		".webp": imageAccept,
		".css":  cssAccept,
		".js":   jsAccept,
	}
	ext := path.Ext(url.String())
	for k, v := range extAcceptMap {
		if ext == k {
			return v
		}
	}
	return htmlAccept
}

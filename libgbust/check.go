package libgbust

import (
	"net/http"
	"net/url"
)

// CheckDir is used to execute a directory check
func (a *Attacker) CheckDir(word string) *Result {
	end, err := url.Parse(word)
	if err != nil {
		return &Result{
			Msg: "[!] failed to parse word",
			Err: err,
		}
	}
	fullURL := a.config.URL.ResolveReference(end)
	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return &Result{
			Msg: "[!] failed to create new request",
			Err: err,
		}
	}

	for _, cookie := range a.config.Cookies {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return &Result{
			Err: err,
			Msg: "[!] failed to do request",
		}
	}
	defer resp.Body.Close()

	return &Result{
		StatusCode: resp.StatusCode,
		Size:       resp.ContentLength,
		URL:        resp.Request.URL,
	}
}

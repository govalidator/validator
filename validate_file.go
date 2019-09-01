package govalidator

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// validateFiles validate file size, mimes, extension etc
func validateFiles(r *http.Request, field, rule, msg string, errsBag url.Values) {
	_, _, ext, mime, size, fErr := getFileInfo(r, field)
	// check size
	if strings.HasPrefix(rule, "size:") {
		l, err := strconv.ParseInt(strings.TrimPrefix(rule, "size:"), 10, 64)
		if err != nil {
			panic(errStringToInt)
		}
		if size > l {
			if msg != "" {
				errsBag.Add(field, fmt.Sprintf(msg, l))
			} else {
				errsBag.Add(field, fmt.Sprintf("The %s field size is can not be greater than %d bytes", field, l))
			}
		}
		if fErr != nil {
			errsBag.Add(field, fmt.Sprintf("The %s field failed to read file when fetching size", field))
		}
	}

	// check extension
	if strings.HasPrefix(rule, "ext:") {
		exts := strings.Split(strings.TrimPrefix(rule, "ext:"), ",")
		f := false
		for _, e := range exts {
			if e == ext {
				f = true
			}
		}
		if !f {
			if msg != "" {
				errsBag.Add(field, fmt.Sprintf(msg, ext))
			} else {
				errsBag.Add(field, fmt.Sprintf("The %s field file extension %s is invalid", field, ext))
			}
		}
		if fErr != nil {
			errsBag.Add(field, fmt.Sprintf("The %s field failed to read file when fetching extension", field))
		}
	}

	// check mimes
	if strings.HasPrefix(rule, "mime:") {
		mimes := strings.Split(strings.TrimPrefix(rule, "mime:"), ",")
		f := false
		for _, m := range mimes {
			if m == mime {
				f = true
			}
		}
		if !f {
			if msg != "" {
				errsBag.Add(field, fmt.Sprintf(msg, mime))
			} else {
				errsBag.Add(field, fmt.Sprintf("The %s field file mime %s is invalid", field, mime))
			}
		}
		if fErr != nil {
			errsBag.Add(field, fmt.Sprintf("The %s field failed to read file when fetching mime", field))
		}
	}
}

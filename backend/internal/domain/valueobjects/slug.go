package valueobjects
// Package valueobjects holds immutable, self-validating domain values.
package valueobjects

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// MaxSlugLength bounds a slug so it stays readable inside a URL.
const MaxSlugLength = 60

var (
	// ErrSlugEmpty is returned when the input contains no URL-safe characters.
	ErrSlugEmpty = errors.New("slug: name must contain at least one letter or digit")

	slugPattern      = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
	nonSlugCharacter = regexp.MustCompile(`[^a-z0-9]+`)
)

// Slug is the URL-safe identifier of a pantry, e.g. "familia-suarez".
type Slug string

// NewSlug converts a free-text pantry name into a Slug.
//
// Accents are folded ("Suárez" -> "suarez"), every other non-alphanumeric run
// collapses into a single hyphen, and the result is truncated to
// MaxSlugLength without leaving a trailing hyphen.
func NewSlug(raw string) (Slug, error) {
	folded, err := foldAccents(strings.ToLower(strings.TrimSpace(raw)))
	if err != nil {
		return "", err
	}

	candidate := strings.Trim(nonSlugCharacter.ReplaceAllString(folded, "-"), "-")
	if len(candidate) > MaxSlugLength {
		candidate = strings.Trim(candidate[:MaxSlugLength], "-")
	}

	if !slugPattern.MatchString(candidate) {
		return "", ErrSlugEmpty
	}
	return Slug(candidate), nil
}

// ParseSlug validates an already-slugified value, typically taken from a URL.
func ParseSlug(raw string) (Slug, error) {
	if len(raw) > MaxSlugLength || !slugPattern.MatchString(raw) {
		return "", ErrSlugEmpty
	}
	return Slug(raw), nil
}

// String returns the slug as a plain string.
func (s Slug) String() string { return string(s) }

// foldAccents decomposes characters and drops the combining marks, so "ñ"
// becomes "n" instead of being replaced by a hyphen.
func foldAccents(value string) (string, error) {
	chain := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	folded, _, err := transform.String(chain, value)
	if err != nil {
		return "", err
	}
	return folded, nil
}

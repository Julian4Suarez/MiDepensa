package unit
package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"midepensa/internal/domain/valueobjects"
)

func TestNewSlug_WithVariousNames_ProducesUrlSafeValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "lowercases and hyphenates", input: "Familia Suarez", want: "familia-suarez"},
		{name: "folds accents", input: "Familia Suárez", want: "familia-suarez"},
		{name: "folds enye", input: "La Peña", want: "la-pena"},
		{name: "collapses separator runs", input: "casa   ---   playa", want: "casa-playa"},
		{name: "drops unsafe characters", input: "casa/../etc?x=1", want: "casa-etc-x-1"},
		{name: "trims leading and trailing separators", input: "  --Casa--  ", want: "casa"},
		{name: "keeps digits", input: "Piso 2", want: "piso-2"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			slug, err := valueobjects.NewSlug(test.input)

			require.NoError(t, err)
			assert.Equal(t, test.want, slug.String())
		})
	}
}

func TestNewSlug_WithoutAlphanumericCharacters_ReturnsError(t *testing.T) {
	for _, input := range []string{"", "   ", "///", "!!!"} {
		_, err := valueobjects.NewSlug(input)

		assert.ErrorIs(t, err, valueobjects.ErrSlugEmpty, "input %q", input)
	}
}

func TestNewSlug_WithLongName_TruncatesWithoutTrailingHyphen(t *testing.T) {
	long := "aaaaaaaaaa bbbbbbbbbb cccccccccc dddddddddd eeeeeeeeee ffffffffff gggg"

	slug, err := valueobjects.NewSlug(long)

	require.NoError(t, err)
	assert.LessOrEqual(t, len(slug.String()), valueobjects.MaxSlugLength)
	assert.NotContains(t, []byte{slug.String()[len(slug.String())-1]}, byte('-'))
}

func TestParseSlug_RejectsValuesThatAreNotAlreadySlugified(t *testing.T) {
	for _, input := range []string{"Familia", "familia_suarez", "-familia", "familia-", "fa milia"} {
		_, err := valueobjects.ParseSlug(input)

		assert.Error(t, err, "input %q", input)
	}
}

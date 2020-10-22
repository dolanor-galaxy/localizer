package localizer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"

	"github.com/razor-1/localizer"
	"github.com/razor-1/localizer/store"
)

var enTag = language.Make("en")

func TestNewLocale(t *testing.T) {
	ta := assert.New(t)

	l, err := localizer.NewLocale(enTag)
	ta.NoError(err)
	ta.Equal(enTag, l.Tag)
	ta.NotZero(l.Calendar)
	ta.NotZero(l.Number)
	ta.NotZero(l.Plural)

	ta.NotNil(localizer.GetLocale(enTag))
}

func TestNamedParameters(t *testing.T) {
	ta := assert.New(t)

	bob := localizer.FmtParams{"name": "bob"}
	tests := []struct {
		format string
		params localizer.FmtParams
		out    string
	}{
		{
			format: "%(name)s is cool",
			params: bob,
			out:    "bob is cool",
		},
		{
			format: "what about %(name)",
			params: bob,
			out:    "what about bob",
		},
		{
			format: "%(num)d",
			params: localizer.FmtParams{"num": 22},
			out:    "22",
		},
		{
			format: "%(name)s has %(num)d and %(name) needs %(num2)d",
			params: localizer.FmtParams{"name": "bob", "num": 12, "num2": 20},
			out:    "bob has 12 and bob needs 20",
		},
		{
			format: "%(name)s has %(num)d and %(name)s needs %(num)d",
			params: localizer.FmtParams{"name": "bob", "num": 12},
			out:    "bob has 12 and bob needs 12",
		},
	}

	for _, test := range tests {
		t.Run(test.format, func(t *testing.T) {
			out := localizer.NamedParameters(test.format, test.params)
			ta.Equal(test.out, out)
		})
	}
}

type loader struct {
	catalog store.LocaleCatalog
}

func (ld *loader) GetTranslations(_ language.Tag) (store.LocaleCatalog, error) {
	return ld.catalog, nil
}

func TestLocale_Load(t *testing.T) {
	const msgID = "test"
	const msgStr = "testXlate"
	const pluralMinutes = "%d minutes"
	const pluralMinute = "%d minute"

	ta := assert.New(t)

	l, err := localizer.NewLocale(enTag)
	ta.NoError(err)

	translations := make(map[string]*store.Translation)
	translations[msgID] = &store.Translation{
		ID:     msgID,
		String: msgStr,
	}

	pluralTrans := &store.Translation{
		ID:       "plural-minutes",
		PluralID: pluralMinutes,
		String:   pluralMinute,
		Plurals:  map[plural.Form]string{plural.One: pluralMinute, plural.Other: pluralMinutes},
	}
	translations[pluralTrans.PluralID] = pluralTrans

	testStore := &loader{
		catalog: store.LocaleCatalog{
			Tag:          enTag,
			Translations: translations,
		},
	}

	err = l.Load(testStore)
	ta.NoError(err)
	ta.Equal(msgStr, l.Get(msgID))

	prt := l.NewPrinter()
	ta.Equal(msgStr, prt.Sprintf(msgID))

	ta.Equal("1 minute", l.GetPlural(pluralMinutes, 1, 1))
	ta.Equal("2 minutes", l.GetPlural(pluralMinutes, 2, 2))
}
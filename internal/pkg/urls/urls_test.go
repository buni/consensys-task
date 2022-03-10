package urls

import (
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/buni/scraper/internal/pkg/test"
)

func tempFileHelper(t *testing.T, urls []string) string {
	t.Helper()
	dir := t.TempDir()
	fileName := strconv.Itoa(rand.Int())
	tempFile, err := os.CreateTemp(dir, fileName)
	if err != nil {
		t.Error(err)
	}
	defer tempFile.Close()
	body := ""
	for _, v := range urls {
		body += v + "\n"
	}

	_, err = tempFile.Write([]byte(body))
	if err != nil {
		t.Error(err)
	}
	return tempFile.Name()
}

func TestParseURLs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		setup    func(t *testing.T, urls []string) string
		wantUrls []string
		wantErr  bool
	}{
		{
			name: "test valid url list",
			setup: func(t *testing.T, urls []string) string {
				return tempFileHelper(t, urls)
			},
			wantUrls: []string{
				"https://stackoverflow.com",
				"http://stackoverflow.com",
				"https://stackoverflow.com/",
				"https://stackoverflow.com/questions",
				"https://stackoverflow.com/questions/",
			},
			wantErr: false,
		},
		{
			name: "test empty url",
			setup: func(t *testing.T, urls []string) string {
				return tempFileHelper(t, urls)
			},
			wantUrls: []string{
				"",
			},
			wantErr: true,
		},
		{
			name: "test bad url",
			setup: func(t *testing.T, urls []string) string {
				return tempFileHelper(t, urls)
			},
			wantUrls: []string{
				"\t",
			},
			wantErr: true,
		},
		{
			name: "test no url schema",
			setup: func(t *testing.T, urls []string) string {
				return tempFileHelper(t, urls)
			},
			wantUrls: []string{
				"stackoverflow.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := tt.setup(t, tt.wantUrls)
			f, err := os.Open(path)
			if err != nil {
				t.Error(err)
			}
			defer f.Close()
			gotUrls, err := ParseURLs(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURLConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) == tt.wantErr {
				return
			}
			t.Log(gotUrls)
			if !reflect.DeepEqual(gotUrls, test.StrToURL(t, tt.wantUrls)) {
				t.Errorf("ParseURLConfig() = %v, want %v", gotUrls, tt.wantUrls)
			}
		})
	}
}

package main

import (
	"fmt"
	"testing"
)

func TestStrip(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"https://github.com/gorilla/mux.git", "github.com/gorilla/mux"},
		{"http://github.com/gorilla/mux.git", "github.com/gorilla/mux"},
		{"github.com/gorilla/mux.git", "github.com/gorilla/mux"},
		{"https://gitlab.com/anuchito/payment.git", "gitlab.com/anuchito/payment"},
		{"http://gitlab.com/anuchito/payment.git", "gitlab.com/anuchito/payment"},
		{"gitlab.com/anuchito/payment.git", "gitlab.com/anuchito/payment"},
	}

	for _, tc := range testCases {
		result := strip(tc.input)
		if result != tc.expected {
			t.Errorf("strip(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestParts(t *testing.T) {
	tests := []struct {
		url     string
		domain  string
		account string
		repo    string
		wantErr error
	}{
		{
			url:     "https://github.com/gorilla/mux",
			domain:  "github.com",
			account: "gorilla",
			repo:    "mux",
			wantErr: nil,
		},
		{
			url:     "https://gitlab.com/anuchito/backend/payment",
			domain:  "gitlab.com",
			account: "anuchito",
			repo:    "backend/payment",
			wantErr: nil,
		},
		{
			url:     "https://example.com/foo",
			domain:  "",
			account: "",
			repo:    "",
			wantErr: errInvalidURLFormat,
		},
		{
			url:     "",
			domain:  "",
			account: "",
			repo:    "",
			wantErr: errInvalidURLFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			gotDomain, gotAccount, gotRepo, err := parts(strip(tt.url))

			if err != tt.wantErr {
				fmt.Println("err:", err)
				fmt.Println("wantErr:", tt.wantErr)
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if gotDomain != tt.domain {
				t.Errorf("Expected domain '%s', but got '%s'", tt.domain, gotDomain)
			}

			if gotAccount != tt.account {
				t.Errorf("Expected account '%s', but got '%s'", tt.account, gotAccount)
			}

			if gotRepo != tt.repo {
				t.Errorf("Expected repo '%s', but got '%s'", tt.repo, gotRepo)
			}
		})
	}
}

func TestRoot(t *testing.T) {
	t.Run("with GOPATH environment variable set", func(t *testing.T) {
		mock := func() (string, error) {
			return "", nil
		}

		got := rooted("/user/go", GetwdFunc(mock))

		want := "/user/go/src"
		if got != want {
			t.Errorf("Root directory is %s, expected %s", got, want)
		}
	})

	t.Run("without GOPATH environment variable set", func(t *testing.T) {
		mock := func() (string, error) {
			return "/current/path/dir", nil
		}

		got := rooted("", GetwdFunc(mock))
		want := "/current/path/dir"
		if got != want {
			t.Errorf("Root directory is %s, expected %s", got, want)
		}
	})
}

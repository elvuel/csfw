// Copyright 2015 CoreStore Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
	"github.com/spf13/viper"
)

const (
	LeftDelim  = "{{"
	RightDelim = "}}"

	// PathCSBaseURL main CoreStore base URL, used if no configuration on a store level can be found.
	PathCSBaseURL = "web/corestore/base_url"
	CSBaseURL     = "http://localhost:9500/"
)

const (
	URLTypeAbsent URLType = iota
	// UrlTypeWeb defines the ULR type to generate the main base URL.
	URLTypeWeb
	// UrlTypeStatic defines the url to the static assets like css, js or theme images
	URLTypeStatic
	// UrlTypeLink hmmm
	// UrlTypeLink
	// UrlTypeMedia defines the ULR type for generating URLs to product photos
	URLTypeMedia
)

type (
	// UrlType defines the type of the URL. Used in const declaration.
	// @see https://github.com/magento/magento2/blob/0.74.0-beta7/lib/internal/Magento/Framework/UrlInterface.php#L13
	URLType uint8

	Reader interface {
		// GetString returns a string from the manager. Example usage:
		// Default value: GetString(config.Path("general/locale/timezone"))
		// Website value: GetString(config.Path("general/locale/timezone"), config.ScopeWebsite(w))
		// Store   value: GetString(config.Path("general/locale/timezone"), config.ScopeStore(s))
		GetString(...ScopeOption) string

		// GetBool returns bool from the manager. Example usage see GetString.
		GetBool(...ScopeOption) bool
	}

	Writer interface {
		// Write puts a value back into the manager. Example usage:
		// Default Scope: Write(config.Path("currency", "option", "base"), config.Value("USD"))
		// Website Scope: Write(config.Path("currency", "option", "base"), config.Value("EUR"), config.ScopeWebsite(w))
		// Store   Scope: Write(config.Path("currency", "option", "base"), config.ValueReader(resp.Body), config.ScopeStore(s))
		Write(...ScopeOption) error
	}

	// Manager main configuration struct
	Manager struct {
		// why is Viper private? Because it can maybe replaced by something else ...
		v *viper.Viper
	}
)

var (
	_ Reader = (*Manager)(nil)
	_ Writer = (*Manager)(nil)

	// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
	TableCollection csdb.TableStructureSlice
	ErrEmptyKey     = errors.New("Key is empty")
	DefaultManager  *Manager
)

func init() {
	DefaultManager = NewManager()
}

// NewManager creates the main new configuration for all scopes: default, website and store
func NewManager() *Manager {
	s := &Manager{
		v: viper.New(),
	}
	s.v.SetDefault(ScopeKey(Path(PathCSBaseURL)), CSBaseURL)
	return s
}

// ApplyDefaults reads the map and applies the keys and values to the default configuration
func (m *Manager) ApplyDefaults(ss Sectioner) *Manager {
	ctxLog := logger.WithField("Scope", "ApplyDefaults")
	for k, v := range ss.Defaults() {
		ctxLog.Debug(k, v)
		m.v.SetDefault(k, v)
	}
	return m
}

func (m *Manager) ApplyCoreConfigData(dbrSess dbr.SessionRunner) error {
	var ccd TableCoreConfigDataSlice
	if _, err := csdb.LoadSlice(dbrSess, TableCollection, TableIndexCoreConfigData, &ccd); err != nil {
		return errgo.Mask(err)
	}
	// @todo
	return nil
}

// Write puts a value back into the manager. Example usage:
// Default Scope: Write(config.Path("currency", "option", "base"), config.Value("USD"))
// Website Scope: Write(config.Path("currency", "option", "base"), config.Value("EUR"), config.ScopeWebsite(w))
// Store   Scope: Write(config.Path("currency", "option", "base"), config.ValueReader(resp.Body), config.ScopeStore(s))
func (m *Manager) Write(o ...ScopeOption) error {
	k, v := ScopeKeyValue(o...)
	if k == "" {
		return ErrEmptyKey
	}
	m.v.Set(k, v)
	return nil
}

// GetString returns a string from the manager. Example usage:
// Default value: GetString(config.Path("general/locale/timezone"))
// Website value: GetString(config.Path("general/locale/timezone"), config.ScopeWebsite(w))
// Store   value: GetString(config.Path("general/locale/timezone"), config.ScopeStore(s))
func (m *Manager) GetString(o ...ScopeOption) string {
	return m.v.GetString(ScopeKey(o...))
}

// @todo use the backend model of a config value. most/all magento string slices are comma lists.
func (m *Manager) GetStringSlice(o ...ScopeOption) []string {
	return m.v.GetStringSlice(ScopeKey(o...))
}

// GetBool returns bool from the manager. Example usage see GetString.
func (m *Manager) GetBool(o ...ScopeOption) bool {
	return m.v.GetBool(ScopeKey(o...))
}

// GetFloat64 returns a float64 from the manager. Example usage see GetString.
func (m *Manager) GetFloat64(o ...ScopeOption) float64 {
	return m.v.GetFloat64(ScopeKey(o...))
}

// @todo consider adding other Get* from the viper package

// AllKeys return all keys regardless where they are set
func (m *Manager) AllKeys() []string { return m.v.AllKeys() }

var _ Reader = (*mockReader)(nil)

// MockScopeReader used for testing
type mockReader struct {
	s func(path string) string
	b func(path string) bool
}

// NewMockReader used for testing
func NewMockReader(
	s func(path string) string,
	b func(path string) bool,
) *mockReader {

	return &mockReader{
		s: s,
		b: b,
	}
}

func (sr mockReader) GetString(opts ...ScopeOption) string {
	if sr.s == nil {
		return ""
	}
	return sr.s(ScopeKey(opts...))
}

func (sr mockReader) GetBool(opts ...ScopeOption) bool {
	if sr.b == nil {
		return false
	}
	return sr.b(ScopeKey(opts...))
}

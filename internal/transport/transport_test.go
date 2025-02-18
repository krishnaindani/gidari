// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package transport

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestTimeseries(t *testing.T) {
	t.Parallel()
	t.Run("chunks where end date is before last iteration", func(t *testing.T) {
		t.Parallel()

		timeseries := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		testURL, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := testURL.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T00:00:00Z")
		testURL.RawQuery = query.Encode()

		err = timeseries.chunk(*testURL)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 0o5, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 0, 0, 0, 0, time.UTC),
			},
		}

		if !reflect.DeepEqual(expChunks, timeseries.chunks) {
			t.Fatalf("unexpected chunks: %v", timeseries.chunks)
		}
	})

	t.Run("chunks where end date is equal to last iteration", func(t *testing.T) {
		t.Parallel()

		timeseries := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		testURL, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := testURL.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T01:00:00Z")
		testURL.RawQuery = query.Encode()

		err = timeseries.chunk(*testURL)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 0o5, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 1, 0, 0, 0, time.UTC),
			},
		}

		if !reflect.DeepEqual(expChunks, timeseries.chunks) {
			t.Fatalf("unexpected chunks: %v", timeseries.chunks)
		}
	})

	t.Run("chunks where end date is after last iteration", func(t *testing.T) {
		t.Parallel()
		timeseries := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		testURL, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := testURL.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T02:00:00Z")
		testURL.RawQuery = query.Encode()

		err = timeseries.chunk(*testURL)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 0o5, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 1, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 11, 1, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 2, 0, 0, 0, time.UTC),
			},
		}

		if !reflect.DeepEqual(expChunks, timeseries.chunks) {
			t.Fatalf("unexpected chunks: %v", timeseries.chunks)
		}
	})
}

func TestUpsert(t *testing.T) {
	t.Parallel()

	// Iterate over the fixtures/upsert directory and run each configuration file.
	fixtureRoot := "testdata/upsert"

	fixtures, err := os.ReadDir(fixtureRoot)
	if err != nil {
		t.Fatalf("error reading fixtures: %v", err)
	}

	for _, fixture := range fixtures {
		name := fixture.Name()
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(fixtureRoot, name)

			bytes, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("error reading fixture: %v", err)
			}

			cfg, err := NewConfig(bytes)
			if err != nil {
				t.Fatalf("error creating config: %v", err)
			}

			cfg.Logger = logrus.New()

			// Fill in the authentication details for the fixture.
			cfgAuth := cfg.Authentication
			if cfgAuth.APIKey != nil {
				// The "passhprase" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = Authentication{
					APIKey: &APIKey{
						Key:        os.Getenv(cfgAuth.APIKey.Key),
						Secret:     os.Getenv(cfgAuth.APIKey.Secret),
						Passphrase: os.Getenv(cfgAuth.APIKey.Passphrase),
					},
				}
			}

			if cfgAuth.Auth2 != nil {
				// The "bearer" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = Authentication{
					Auth2: &Auth2{
						Bearer: os.Getenv(cfgAuth.Auth2.Bearer),
					},
				}
			}

			// Upsert the fixture.
			if err := Upsert(context.Background(), cfg); err != nil {
				t.Fatalf("error upserting: %v", err)
			}
		})
	}
}

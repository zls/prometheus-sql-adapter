// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"
	"github.com/prometheus/common/model"
)

// Client allows sending batches of Prometheus samples to Postgres.
type Client struct {
	logger log.Logger

	db *sql.DB
}

// NewClient creates a new Client.
func NewClient(logger log.Logger, conn string, idle int, open int) *Client {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	db, err := sql.Open("postgres", conn)
	if err != nil {
		level.Error(logger).Log(err)
	}
	db.SetMaxIdleConns(idle)
	db.SetMaxOpenConns(open)
	return &Client{
		logger:    logger,
		db:        db,
	}
}

// Write sends a batch of samples to Postgres.
func (c *Client) Write(samples model.Samples) error {
	txn, err := c.db.Begin()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, s := range samples {
		k, l := splitKeyAndLabels(s.Metric)
		t := float64(s.Timestamp.UnixNano()) / 1e9
		v := float64(s.Value)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			level.Debug(c.logger).Log("msg", "cannot send value to Postgres, skipping sample", "value", v, "sample", s)
			continue
		}
		fmt.Fprintf(&buf, "%s %f %f (%s)\n", k, v, t, l)
	}

	level.Debug(c.logger).Log("batch", buf.String())

	err = txn.Commit()
	return err
}

// Name identifies the client as a Postgres client.
func (c Client) Name() string {
	return "postgres"
}

func splitKeyAndLabels(m model.Metric) (key string, labels map[string]string) {
	return string(m[model.MetricNameLabel]), make(map[string]string, 0)
}
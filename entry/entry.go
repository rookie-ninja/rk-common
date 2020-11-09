// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_entry

import (
	"github.com/rookie-ninja/rk-query"
	"time"
)

type Entry interface {
	Bootstrap(rk_query.Event)

	Wait(time.Duration)

	Shutdown(rk_query.Event)

	GetName() string

	GetType() string

	String() string
}

// Copyright 2018 Tony Yip. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package vodka

// Logger defines a logging interface.
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type NullLogger struct{}

func (*NullLogger) Debug(v ...interface{})                 {}
func (*NullLogger) Debugf(format string, v ...interface{}) {}

func (*NullLogger) Info(v ...interface{})                 {}
func (*NullLogger) Infof(format string, v ...interface{}) {}

func (*NullLogger) Error(v ...interface{})                 {}
func (*NullLogger) Errorf(format string, v ...interface{}) {}

func (*NullLogger) Fatal(v ...interface{})                 {}
func (*NullLogger) Fatalf(format string, v ...interface{}) {}

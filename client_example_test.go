// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dosa_test

import (
	"context"

	"github.com/pkg/errors"

	"github.com/uber-go/dosa"
	"github.com/uber-go/dosa/connectors/memory"
)

func ExampleInMemoryClient() {
	// initialize registrar
	reg, err := dosa.NewRegistrar("test", "myteam.myservice", cte1)
	if err != nil {
		// registration will fail if the object is tagged incorrectly
		panic("dosaRenamed.NewRegister returned an error")
	}

	// make an in-memory connector
	conn := memory.NewConnector()

	// create the client using the registry and connector
	client := dosa.NewClient(reg, conn)

	err = client.Initialize(context.Background())
	if err != nil {
		errors.Wrap(err, "client.Initialize returned an error")
	}
}
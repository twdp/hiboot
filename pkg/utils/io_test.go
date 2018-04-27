// Copyright 2018 John Deng (hi.devops.io@gmail.com).
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

package utils

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func TestGetWorkingDir(t *testing.T) {
	wd := GetWorkingDir("")
	expected, err := os.Getwd()
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, wd)
}

func TestWriteFile(t *testing.T) {
	in := "hello, world"
	buf := bytes.NewBufferString(in)
	path := os.TempDir()
	log.Println("path: ", path)
	n, err := WriterFile(path, "foo.txt", buf.Bytes())
	assert.Equal(t, nil, err)
	assert.Equal(t, len(in), n)
}

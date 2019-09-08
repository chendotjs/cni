// Copyright 2016 CNI authors
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

package invoke

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/containernetworking/cni/pkg/types"
)

type RawExec struct {
	Stderr io.Writer
}

func (e *RawExec) ExecPlugin(ctx context.Context, pluginPath string, stdinData []byte, environ []string) ([]byte, error) {
	fmt.Printf("ExecPlugin begin\n")
	fmt.Printf("ExecPlugin pluginPath: %v\n", pluginPath)
	fmt.Printf("ExecPlugin stdinData: %v\n", string(stdinData))
	fmt.Printf("ExecPlugin environ: %v\n", environ)
	for _, env := range environ {
		if strings.Contains(env, "CNI_") {
			fmt.Println(env)
		}
	}
	// 一个CNI插件就是一个可执行文件，我们从配置文件中获取network配置信息，从容器管理系统处获取运行时信息，再将前者以标准输入的形式，后者以环境变量的形式传递传递给插件

	stdout := &bytes.Buffer{}
	c := exec.CommandContext(ctx, pluginPath)
	c.Env = environ
	c.Stdin = bytes.NewBuffer(stdinData)
	c.Stdout = stdout
	c.Stderr = e.Stderr
	if err := c.Run(); err != nil {
		return nil, pluginErr(err, stdout.Bytes())
	}

	return stdout.Bytes(), nil
}

func pluginErr(err error, output []byte) error {
	if _, ok := err.(*exec.ExitError); ok {
		emsg := types.Error{}
		if len(output) == 0 {
			emsg.Msg = "netplugin failed with no error message"
		} else if perr := json.Unmarshal(output, &emsg); perr != nil {
			emsg.Msg = fmt.Sprintf("netplugin failed but error parsing its diagnostic message %q: %v", string(output), perr)
		}
		return &emsg
	}

	return err
}

func (e *RawExec) FindInPath(plugin string, paths []string) (string, error) {
	return FindInPath(plugin, paths)
}

// Copyright 2020
//
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

package trace

// func TestSLACKWriter_WriteMsg(t *testing.T) {
// 	sc := `
// {
//   "webhookurl":"",
//   "level":7
// }
// `
// 	l := newSLACKWriter()
// 	err := l.Init(sc)
// 	if err != nil {
// 		Debug(err)
// 	}
//
// 	err = l.WriteMsg(&LogMsg{
// 		Level: 7,
// 		Msg: `{ "abs"`,
// 		When: time.Now(),
// 		FilePath: "main.go",
// 		LineNumber: 100,
// 		enableFullFilePath: true,
// 		enableFuncCallDepth: true,
// 	})
//
// 	if err != nil {
// 		Debug(err)
// 	}
//
// }

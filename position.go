// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

var newline = []byte{itemNewline}

// unidim2Bidim 用于将一维 tokens 中的列 uCol 变换为对应二维中的行 bLn 和列 bCol。
func (t *Tree) unidim2Bidim(tokens items, uCol int) (bLn, bCol int) {
	bLn = 1
	length := len(tokens)
	var token byte
	for i := 0; i < length && i < uCol; i++ {
		token = tokens[i]
		if itemNewline == token {
			bLn++
			bCol = 1
			continue
		}
		bCol++
	}
	bCol -= bLn - 1 // 减去 \n 个数
	return
}

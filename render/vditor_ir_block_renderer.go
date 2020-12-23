// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/88250/lute/html/atom"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// VditorIRBlockRenderer 描述了 Vditor Instant-Rendering Block DOM 渲染器。
type VditorIRBlockRenderer struct {
	*BaseRenderer
}

// NewVditorIRBlockRenderer 创建一个 Vditor Instant-Rendering Block DOM 渲染器。
func NewVditorIRBlockRenderer(tree *parse.Tree) *VditorIRBlockRenderer {
	ret := &VditorIRBlockRenderer{BaseRenderer: NewBaseRenderer(tree)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeText] = ret.renderText
	ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.RendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.RendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.RendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.RendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.RendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.RendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.RendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	ret.RendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.RendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	ret.RendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.RendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.RendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.RendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.RendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.RendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.RendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.RendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.RendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.RendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.RendererFuncs[ast.NodeHeadingID] = ret.renderHeadingID
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	ret.RendererFuncs[ast.NodeBang] = ret.renderBang
	ret.RendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	ret.RendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	ret.RendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	ret.RendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	ret.RendererFuncs[ast.NodeOpenBrace] = ret.renderOpenBrace
	ret.RendererFuncs[ast.NodeCloseBrace] = ret.renderCloseBrace
	ret.RendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	ret.RendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	ret.RendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.RendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.RendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.RendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	ret.RendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.RendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.RendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.RendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.RendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.RendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.RendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	ret.RendererFuncs[ast.NodeFootnotesDefBlock] = ret.renderFootnotesDefBlock
	ret.RendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	ret.RendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	ret.RendererFuncs[ast.NodeYamlFrontMatterOpenMarker] = ret.renderYamlFrontMatterOpenMarker
	ret.RendererFuncs[ast.NodeYamlFrontMatterContent] = ret.renderYamlFrontMatterContent
	ret.RendererFuncs[ast.NodeYamlFrontMatterCloseMarker] = ret.renderYamlFrontMatterCloseMarker
	ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	ret.RendererFuncs[ast.NodeBlockRefID] = ret.renderBlockRefID
	ret.RendererFuncs[ast.NodeBlockRefSpace] = ret.renderBlockRefSpace
	ret.RendererFuncs[ast.NodeBlockRefText] = ret.renderBlockRefText
	ret.RendererFuncs[ast.NodeBlockRefTextTplRenderResult] = ret.renderBlockRefTextTplRenderResult
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeMark1OpenMarker] = ret.renderMark1OpenMarker
	ret.RendererFuncs[ast.NodeMark1CloseMarker] = ret.renderMark1CloseMarker
	ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSupOpenMarker] = ret.renderSupOpenMarker
	ret.RendererFuncs[ast.NodeSupCloseMarker] = ret.renderSupCloseMarker
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeSubOpenMarker] = ret.renderSubOpenMarker
	ret.RendererFuncs[ast.NodeSubCloseMarker] = ret.renderSubCloseMarker
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.RendererFuncs[ast.NodeBlockQueryEmbedScript] = ret.renderBlockQueryEmbedScript
	ret.RendererFuncs[ast.NodeBlockEmbed] = ret.renderBlockEmbed
	ret.RendererFuncs[ast.NodeBlockEmbedID] = ret.renderBlockEmbedID
	ret.RendererFuncs[ast.NodeBlockEmbedSpace] = ret.renderBlockEmbedSpace
	ret.RendererFuncs[ast.NodeBlockEmbedText] = ret.renderBlockEmbedText
	ret.RendererFuncs[ast.NodeBlockEmbedTextTplRenderResult] = ret.renderBlockEmbedTextTplRenderResult
	ret.RendererFuncs[ast.NodeTag] = ret.renderTag
	ret.RendererFuncs[ast.NodeTagOpenMarker] = ret.renderTagOpenMarker
	ret.RendererFuncs[ast.NodeTagCloseMarker] = ret.renderTagCloseMarker
	ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	ret.RendererFuncs[ast.NodeSuperBlock] = ret.renderSuperBlock
	ret.RendererFuncs[ast.NodeSuperBlockOpenMarker] = ret.renderSuperBlockOpenMarker
	ret.RendererFuncs[ast.NodeSuperBlockLayoutMarker] = ret.renderSuperBlockLayoutMarker
	ret.RendererFuncs[ast.NodeSuperBlockCloseMarker] = ret.renderSuperBlockCloseMarker
	return ret
}

func (r *VditorIRBlockRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "super-block-open-marker"}}, false)
		r.Write([]byte("{{{"))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--info"}, {"data-type", "super-block-layout"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "super-block-close-marker"}}, false)
		r.Write([]byte("}}}"))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"link-ref-defs-block\">")
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.FirstChild.ChildByType(ast.NodeLinkDest).Tokens
		destStr := util.BytesToStr(dest)
		r.WriteString("<span><span>[" + util.BytesToStr(node.Tokens) + "]:")
		if util.Caret != destStr {
			r.WriteString(" </span>")
		}
		r.WriteString(destStr + "\n</span>")
	}
	return ast.WalkSkipChildren
}

func (r *VditorIRBlockRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--tag"}}, false)
		r.WriteByte(lex.ItemCrosshatch)
		r.Tag("/span", nil, false)
		attrs := [][]string{{"class", "vditor-ir__tag"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("span", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--tag"}}, false)
		r.WriteByte(lex.ItemCrosshatch)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.Write(node.Tokens)
		r.WriteString("</span></span>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--mark"}}, false)
		r.WriteString("=")
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("mark", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--mark"}}, false)
		r.WriteString("=")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--mark"}}, false)
		r.WriteString("==")
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("mark", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--mark"}}, false)
		r.WriteString("==")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--sup"}}, false)
		r.WriteString("^")
		r.Tag("/span", nil, false)
		r.Tag("sup", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--sup"}}, false)
		r.WriteString("^")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--sub"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
		r.Tag("sub", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--sub"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) render() (output []byte) {
	r.LastOut = lex.ItemNewline
	r.Writer = &bytes.Buffer{}
	r.Writer.Grow(4096)

	ast.Walk(r.Tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		extRender := r.ExtRendererFuncs[n.Type]
		if nil != extRender {
			output, status := extRender(n, entering)
			r.WriteString(output)
			return status
		}

		render := r.RendererFuncs[n.Type]
		if nil == render {
			if nil != r.DefaultRendererFunc {
				return r.DefaultRendererFunc(n, entering)
			} else {
				return r.renderDefault(n, entering)
			}
		}
		return render(n, entering)
	})

	output = r.Writer.Bytes()
	return
}

func (r *VditorIRBlockRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
		r.Write(node.Tokens)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("<div data-render=\"2\" data-type=\"block-render\"></div></div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		id := node.ChildByType(ast.NodeBlockEmbedID)
		r.WriteString("<div data-block-def-id=\"" + string(id.Tokens) + "\" data-render=\"2\" data-type=\"block-render\"></div>")
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockEmbedID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		attrs = append(attrs, []string{"class", "vditor-ir__marker vditor-ir__marker--link"})
		r.Tag("span", attrs, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockEmbedSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemSpace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockEmbedText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		text := html.EscapeHTML(node.Tokens)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDoublequote)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--info"}}, false)
		r.Write(text)
	} else {
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDoublequote)
		r.Tag("/span", nil, false)
		if nil == node.Next || ast.NodeBlockEmbedTextTplRenderResult != node.Next.Type {
			r.Tag("span", [][]string{{"data-type", "ref-text-tpl-render-result"}, {"class", "vditor-ir__blockref"}}, false)
			r.Tag("/span", nil, false)
		}
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockEmbedTextTplRenderResult(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "ref-text-tpl-render-result"}, {"class", "vditor-ir__blockref"}}, false)
		r.Write(node.Tokens)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		attrs = append(attrs, []string{"class", "vditor-ir__marker vditor-ir__marker--link"})
		r.Tag("span", attrs, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemSpace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDoublequote)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--info"}}, false)
	} else {
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDoublequote)
		r.Tag("/span", nil, false)
		if nil == node.Next || ast.NodeBlockRefTextTplRenderResult != node.Next.Type {
			r.Tag("span", [][]string{{"data-type", "ref-text-tpl-render-result"}, {"class", "vditor-ir__blockref"}}, false)
			r.Tag("/span", nil, false)
		}
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockRefTextTplRenderResult(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "ref-text-tpl-render-result"}, {"class", "vditor-ir__blockref"}}, false)
		r.Write(node.Tokens)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "yaml-front-matter-close-marker"}}, false)
		r.Write(parse.YamlFrontMatterMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		node.Tokens = bytes.TrimSpace(node.Tokens)
		codeLen := len(node.Tokens)
		codeIsEmpty := 1 > codeLen || (len(util.Caret) == codeLen && util.Caret == string(node.Tokens))
		r.Tag("pre", [][]string{{"class", "vditor-ir__marker--pre"}}, false)
		r.Tag("code", [][]string{{"data-type", "yaml-front-matter"}, {"class", "language-yaml"}}, false)
		if codeIsEmpty {
			r.WriteString(util.FrontEndCaret + "\n")
		} else {
			r.Write(html.EscapeHTML(node.Tokens))
		}
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "yaml-front-matter-open-marker"}}, false)
		r.Write(parse.YamlFrontMatterMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
		r.Tag("code", [][]string{{"data-newline", "1"}, {"class", "vditor-ir__marker vditor-ir__marker--pre"}, {"data-type", "html-entity"}}, false)
		r.Write(html.EscapeHTML(node.HtmlEntityTokens))
		r.Tag("/code", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		r.Tag("code", nil, false)
		r.Write(node.HtmlEntityTokens)
		r.Tag("/code", nil, false)
		r.Tag("/span", nil, false)
		r.Tag("/span", nil, false)

		nextNodeText := node.NextNodeText()
		nextNodeText = strings.ReplaceAll(nextNodeText, util.Caret, "")
		if "" == nextNodeText {
			r.WriteString(parse.Zwsp)
		}
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--backslash"}}, false)
		r.WriteByte(lex.ItemBackslash)
		r.WriteString("</span>")
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *VditorIRBlockRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"footnotes-block\">")
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if r.RenderingFootnotes {
		return ast.WalkContinue
	}

	if entering {
		r.WriteString("<div data-type=\"footnotes-def\">")
		r.WriteString("<span>[" + string(node.Tokens) + "]: </span>")
		for c := node.FirstChild; nil != c; c = c.Next {
			ast.Walk(c, func(n *ast.Node, entering bool) ast.WalkStatus {
				return r.RendererFuncs[n.Type](n, entering)
			})
		}
		r.WriteString("</div>")
	}
	return ast.WalkSkipChildren
}

func (r *VditorIRBlockRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, util.Caret, "")
		if "" == previousNodeText {
			r.WriteString(parse.Zwsp)
		}
		idx, def := r.Tree.FindFootnotesDef(node.Tokens)
		if nil == def {
			r.Write(node.Tokens)
			return ast.WalkContinue
		}

		idxStr := strconv.Itoa(idx)
		label := def.Text()
		attrs := [][]string{{"data-type", "footnotes-ref"}}
		text := node.Text()
		expand := strings.Contains(text, util.Caret)
		if expand {
			attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand b3-tooltips b3-tooltips__s"})
		} else {
			attrs = append(attrs, []string{"class", "vditor-ir__node b3-tooltips b3-tooltips__s"})
		}
		attrs = append(attrs, []string{"aria-label", SubStr(html.EscapeString(label), 24)})
		attrs = append(attrs, []string{"data-footnotes-label", string(node.FootnotesRefLabel)})
		r.Tag("sup", attrs, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		r.WriteByte(lex.ItemOpenBracket)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--link"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker--hide"}, {"data-render", "1"}}, false)
		r.WriteString(idxStr)
		r.Tag("/span", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		r.WriteByte(lex.ItemCloseBracket)
		r.Tag("/span", nil, false)
		r.WriteString("</sup>" + parse.Zwsp)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "code-block-close-marker"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--info"}, {"data-type", "code-block-info"}}, false)
		r.WriteString(parse.Zwsp)
		r.Write(node.CodeBlockInfo)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "code-block-open-marker"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(util.Caret) == codeLen && util.Caret == string(node.Tokens))
	isFenced := node.Parent.IsFencedCodeBlock
	var language string
	caretInInfo := false
	if isFenced {
		caretInInfo = bytes.Contains(node.Previous.CodeBlockInfo, util.CaretTokens)
		node.Previous.CodeBlockInfo = bytes.ReplaceAll(node.Previous.CodeBlockInfo, util.CaretTokens, nil)
	}
	var attrs [][]string
	if isFenced && 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
		language = string(infoWords[0])
		attrs = append(attrs, []string{"class", "language-" + language})
		if "mindmap" == language {
			dataCode := r.renderMindmap(node.Tokens)
			attrs = append(attrs, []string{"data-code", string(dataCode)})
		}
	}

	class := "vditor-ir__marker--pre vditor-ir__marker"
	r.Tag("pre", [][]string{{"class", class}}, false)
	r.Tag("code", attrs, false)
	if codeIsEmpty {
		if !caretInInfo {
			r.WriteString(util.FrontEndCaret)
		}
		r.WriteByte(lex.ItemNewline)
	} else {
		r.Write(html.EscapeHTML(node.Tokens))
		r.Newline()
	}
	r.WriteString("</code></pre>")

	r.Tag("pre", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
	preDiv := noHighlight(language)
	if preDiv {
		r.Tag("div", attrs, false)
	} else {
		r.Tag("code", attrs, false)
	}
	tokens := node.Tokens
	tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
	r.Write(html.EscapeHTML(tokens))
	if preDiv {
		r.WriteString("</div></pre>")
	} else {
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-render=\"2\">")
		r.Write(node.Tokens)
		r.WriteString("</span>")
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.Write(node.FirstChild.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-render=\"2\">")
		r.Write(node.Tokens)
		r.WriteString("</span>")
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.Write(node.FirstChild.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, util.Caret, "")
		if "" == previousNodeText {
			r.WriteString(parse.Zwsp)
		}
		r.renderSpanNode(node)
	} else {
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDollar)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Tokens)
		r.Write(tokens)
		r.Tag("/code", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		r.Tag("span", [][]string{{"class", "language-math"}}, false)
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		r.Write(tokens)
		r.WriteString("</span></span>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDollar)
		r.Tag("/span", nil, false)
		r.Tag("code", [][]string{{"data-newline", "1"}, {"class", "vditor-ir__marker vditor-ir__marker--pre"}, {"data-type", "math-inline"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		preText := node.PreviousNodeText()
		if "" != preText && !strings.HasSuffix(preText, " ") {
			r.WriteByte(lex.ItemSpace)
		}
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
		nextText := node.NextNodeText()
		if "" != nextText && !strings.HasPrefix(nextText, " ") {
			r.WriteByte(lex.ItemSpace)
		}
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "math-block-close-marker"}}, false)
		r.Write(parse.MathBlockMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	node.Tokens = bytes.TrimSpace(node.Tokens)
	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(util.Caret) == codeLen && util.Caret == string(node.Tokens))

	class := "vditor-ir__marker--pre"
	if r.Options.VditorMathBlockPreview {
		class += " vditor-ir__marker"
	}
	r.Tag("pre", [][]string{{"class", class}}, false)
	r.Tag("code", [][]string{{"data-type", "math-block"}, {"class", "language-math"}}, false)
	if codeIsEmpty {
		r.WriteString(util.FrontEndCaret + "\n")
	} else {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	r.WriteString("</code></pre>")

	if r.Options.VditorMathBlockPreview {
		r.Tag("pre", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		r.Tag("div", [][]string{{"data-type", "math-block"}, {"class", "language-math"}}, false)
		tokens := node.Tokens
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		r.Write(html.EscapeHTML(tokens))
		r.WriteString("</div></pre>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "math-block-open-marker"}}, false)
		r.Write(parse.MathBlockMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	tag := "td"
	if ast.NodeTableHead == node.Parent.Parent.Type {
		tag = "th"
	}
	if entering {
		var attrs [][]string
		switch node.TableCellAlign {
		case 1:
			attrs = append(attrs, []string{"align", "left"})
		case 2:
			attrs = append(attrs, []string{"align", "center"})
		case 3:
			attrs = append(attrs, []string{"align", "right"})
		}
		r.Tag(tag, attrs, false)
		if nil == node.FirstChild {
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		} else if bytes.Equal(node.FirstChild.Tokens, util.CaretTokens) {
			node.FirstChild.Tokens = []byte(util.Caret + " ")
		}
	} else {
		r.Tag("/"+tag, nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
	} else {
		r.Tag("/tr", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("thead", nil, false)
	} else {
		r.Tag("/thead", nil, false)
		if nil != node.Next {
			r.Tag("tbody", nil, false)
		}
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"data-block", "0"}, {"data-type", "table"}, {"data-node-id", r.NodeID(node)}}
		ial := r.NodeAttrs(node)
		if 0 < len(ial) {
			attrs = append(attrs, ial...)
		}
		r.Tag("table", attrs, false)
	} else {
		if nil != node.FirstChild.Next {
			r.Tag("/tbody", nil, false)
		}
		r.Tag("/table", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--s"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("s", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--s"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--s"}}, false)
		r.WriteString("~~")
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("s", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--s"}}, false)
		r.WriteString("~~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--title"}}, false)
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--link"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeImage == node.Parent.Type {
			r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		} else {
			if 3 == node.Parent.LinkType {
				r.Tag("span", nil, false)
			} else {
				r.Tag("span", [][]string{{"class", "vditor-ir__link"}}, false)
			}
		}
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--paren"}}, false)
		r.WriteByte(lex.ItemCloseParen)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--paren"}}, false)
		r.WriteByte(lex.ItemOpenParen)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--brace"}}, false)
		r.WriteByte(lex.ItemCloseBrace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--brace"}}, false)
		r.WriteByte(lex.ItemOpenBrace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		r.WriteByte(lex.ItemCloseBracket)
		r.Tag("/span", nil, false)

		if 3 == node.Parent.LinkType {
			linkText := node.Parent.ChildByType(ast.NodeLinkText)
			if !bytes.EqualFold(node.Parent.LinkRefLabel, linkText.Tokens) {
				r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--link"}}, false)
				r.WriteByte(lex.ItemOpenBracket)
				r.Write(node.Parent.LinkRefLabel)
				r.WriteByte(lex.ItemCloseBracket)
				r.Tag("/span", nil, false)
			}
		}
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		r.WriteByte(lex.ItemOpenBracket)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		if ast.NodeBlockEmbed == node.Parent.Type && nil != node.Parent.Previous && ast.NodeTaskListItemMarker == node.Parent.Previous.Type {
			r.WriteByte(lex.ItemSpace)
		}
		r.WriteByte(lex.ItemBang)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	needResetCaret := nil != node.Next && ast.NodeText == node.Next.Type && bytes.HasPrefix(node.Next.Tokens, util.CaretTokens)
	single := nil == node.Previous && nil == node.Next
	renderFigure := single

	if entering {
		text := r.Text(node)
		class := "vditor-ir__node"
		if strings.Contains(text, util.Caret) || needResetCaret {
			class += " vditor-ir__node--expand"
		}
		if renderFigure {
			class += " img--center"
		}

		attrs := [][]string{{"class", class}, {"data-type", "img"}}
		r.Tag("span", attrs, false)
	} else {
		if needResetCaret {
			r.WriteString(util.Caret)
			node.Next.Tokens = bytes.ReplaceAll(node.Next.Tokens, util.CaretTokens, nil)
		}

		destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
		destTokens = r.Tree.Context.LinkPath(destTokens)
		destTokens = bytes.ReplaceAll(destTokens, util.CaretTokens, nil)
		attrs := [][]string{{"src", string(destTokens)}}
		alt := node.ChildByType(ast.NodeLinkText)
		if nil != alt && 0 < len(alt.Tokens) {
			altTokens := bytes.ReplaceAll(alt.Tokens, util.CaretTokens, nil)
			attrs = append(attrs, []string{"alt", string(altTokens)})
		}

		attrs = append(attrs, r.NodeAttrs(node.Parent)...)
		r.Tag("img", attrs, true)
		// XSS 过滤
		buf := r.Writer.Bytes()
		idx := bytes.LastIndex(buf, []byte("<img src="))
		imgBuf := buf[idx:]
		if r.Options.Sanitize {
			imgBuf = sanitize(imgBuf)
		}
		r.Writer.Truncate(idx)
		r.Writer.Write(imgBuf)

		if renderFigure {
			if title := node.ChildByType(ast.NodeLinkTitle); nil != title {
				titleTokens := title.Tokens
				titleTokens = bytes.ReplaceAll(titleTokens, util.CaretTokens, nil)
				titleTree := parse.Inline("", titleTokens, r.Tree.Context.ParseOption)
				figureTitle := RenderHeadingText(titleTree.Root)
				attrs = [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}
				r.Tag("span", attrs, false)
				r.Writer.WriteString(figureTitle)
				r.Tag("/span", nil, false)
			}
		}

		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	r.renderDivNode(node)

	tokens := bytes.TrimSpace(node.Tokens)
	class := "vditor-ir__marker--pre"
	if r.Options.VditorHTMLBlockPreview {
		class += " vditor-ir__marker"
	}
	r.Tag("pre", [][]string{{"class", class}}, false)
	r.Tag("code", [][]string{{"data-type", "html-block"}}, false)
	r.Write(html.EscapeHTML(tokens))
	r.WriteString("</code></pre>")

	if r.Options.VditorHTMLBlockPreview {
		r.Tag("pre", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		bilibili := []byte("<iframe src=\"//player.bilibili.com/player.html")
		if bytes.HasPrefix(tokens, bilibili) {
			tokens = bytes.Replace(tokens, bilibili, []byte("<iframe class=\"iframe__video\" src=\"https://player.bilibili.com/player.html"), 1)
		}
		r.Write(tokens)
		r.WriteString("</pre>")
	}
	r.WriteString("</div>")
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	openKbd := bytes.Equal(node.Tokens, []byte("<kbd>"))
	closeKbd := bytes.Equal(node.Tokens, []byte("</kbd>"))
	if openKbd || closeKbd {
		if openKbd {
			if r.tagMatchClose("kbd", node) {
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
				r.Tag("kbd", nil, false)
			} else {
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
			}
		} else {
			if r.tagMatchOpen("kbd", node) {
				r.Tag("/kbd", nil, false)
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
			} else {
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
			}
		}
	} else {
		var attrs [][]string
		htmlRoot := &html.Node{Type: html.ElementNode}
		n, _ := html.ParseFragment(strings.NewReader(util.BytesToStr(node.Tokens)), htmlRoot)
		var rendered bool
		content := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		if 0 < len(n) {
			switch n[0].DataAtom {
			case atom.Br:
				r.Tag("span", [][]string{{"data-type", "html-inline"}}, false)
				attrs = [][]string{{"class", "vditor-ir__br"}}
				r.Tag("span", attrs, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/span", nil, false)
				rendered = true
			default:

			}
		}
		if !rendered {
			if !bytes.Equal([]byte("<>"), content) {
				r.Write(content)
			}

			r.renderSpanNode(node)
			attrs = [][]string{{"class", "vditor-ir__marker"}}
			r.Tag("span", attrs, false)
			r.Write(html.EscapeHTML(node.Tokens))
			r.Tag("/span", nil, false)
		}

		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) tagMatchClose(tag string, node *ast.Node) bool {
	for n := node.Next; nil != n; n = n.Next {
		if ast.NodeInlineHTML == n.Type && "</"+tag+">" == n.TokensStr() {
			return true
		}
	}
	return false
}

func (r *VditorIRBlockRenderer) tagMatchOpen(tag string, node *ast.Node) bool {
	for n := node.Previous; nil != n; n = n.Previous {
		if ast.NodeInlineHTML == n.Type && "<"+tag+">" == n.TokensStr() {
			return true
		}
	}
	return false
}

func (r *VditorIRBlockRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.Tight { // List.ListItem.Paragraph
		return ast.WalkContinue
	}

	if entering {
		attrs := [][]string{{"data-block", "0"}, {"data-node-id", r.NodeID(node)}, {"data-type", "p"}}
		ial := r.NodeAttrs(node)
		if 0 < len(ial) {
			attrs = append(attrs, ial...)
		}
		r.Tag("p", attrs, false)
	} else {
		r.Tag("/p", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if r.Options.FixTermTypo {
			tokens = r.FixTermTypo(tokens)
		}
		if r.Options.ChinesePunct {
			tokens = r.ChinesePunct(tokens)
		}

		// 有的场景需要零宽空格撑起，但如果有其他文本内容的话需要把零宽空格删掉
		if !bytes.EqualFold(tokens, []byte(util.Caret+parse.Zwsp)) {
			tokens = bytes.ReplaceAll(tokens, []byte(parse.Zwsp), nil)
		}

		text := string(html.EscapeHTML(tokens))
		_, text = markText(text, r.Tree.Marks, 1024*1024)
		r.WriteString(text)
	}
	return ast.WalkContinue
}

func markText(text string, keywords []string, beforeLen int) (pos int, marked string) {
	marked = encloseIgnoreCase(text, "<mark>", "</mark>", keywords...)

	pos = strings.Index(marked, "<mark>")
	if 0 > pos {
		return
	}

	var before []rune
	var count int
	for i := pos; 0 < i; { // 关键字前面太长的话缩短一些
		r, size := utf8.DecodeLastRuneInString(marked[:i])
		i -= size
		before = append([]rune{r}, before...)
		count++
		if beforeLen < count {
			break
		}
	}
	marked = string(before) + marked[pos:]
	return
}

func encloseIgnoreCase(text, open, close string, searchStrs ...string) string {
	buf := &bytes.Buffer{}
	textLower := strings.ToLower(text)
	for i := 0; i < len(textLower); i++ {
		sub := textLower[i:]
		var found bool
		for j := 0; j < len(searchStrs); j++ {
			idx := strings.Index(sub, strings.ToLower(searchStrs[j]))
			if 0 != idx {
				continue
			}
			buf.WriteString(open)
			buf.WriteString(text[i : i+len(searchStrs[j])])
			buf.WriteString(close)
			i += len(searchStrs[j]) - 1
			found = true
			break
		}
		if !found {
			buf.WriteByte(text[i])
		}
	}
	return buf.String()
}

func (r *VditorIRBlockRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString(strings.Repeat("`", node.Parent.CodeMarkerLen))
		if bytes.HasPrefix(node.Next.Tokens, []byte("`")) {
			r.WriteByte(lex.ItemSpace)
		}
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("code", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/code", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		if bytes.HasSuffix(node.Previous.Tokens, []byte("`")) {
			r.WriteByte(lex.ItemSpace)
		}
		r.WriteString(strings.Repeat("`", node.Parent.CodeMarkerLen))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--em"}}, false)
		r.WriteByte(lex.ItemAsterisk)
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("em", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--em"}}, false)
		r.WriteByte(lex.ItemAsterisk)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--em"}}, false)
		r.WriteByte(lex.ItemUnderscore)
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("em", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--em"}}, false)
		r.WriteByte(lex.ItemUnderscore)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--strong"}}, false)
		r.WriteString("**")
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("strong", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--strong"}}, false)
		r.WriteString("**")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--strong"}}, false)
		r.WriteString("__")
		r.Tag("/span", nil, false)
		attrs := [][]string{{"data-newline", "1"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("strong", attrs, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--strong"}}, false)
		r.WriteString("__")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"data-block", "0"}, {"data-node-id", r.NodeID(node)}, {"data-type", "blockquote"}}
		ial := r.NodeAttrs(node)
		if 0 < len(ial) {
			attrs = append(attrs, ial...)
		}
		r.Tag("blockquote", attrs, false)
	} else {
		r.WriteString("</blockquote>")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		text := r.Text(node)
		headingID := node.ChildByType(ast.NodeHeadingID)
		level := headingLevel[node.HeadingLevel : node.HeadingLevel+1]
		if strings.Contains(text, util.Caret) || (nil != headingID && bytes.Contains(headingID.Tokens, util.CaretTokens)) {
			r.WriteString("<h" + level + " data-block=\"0\" class=\"vditor-ir__node vditor-ir__node--expand\"")
		} else {
			r.WriteString("<h" + level + " data-block=\"0\" class=\"vditor-ir__node\"")
		}

		r.WriteString(" data-node-id=\"" + r.NodeID(node) + "\" " + r.NodeAttrsStr(node) + " data-type=\"h\"")

		var id string
		if nil != headingID {
			id = string(headingID.Tokens)
		}
		if "" == id {
			id = HeadingID(node)
		}
		r.WriteString(" id=\"ir-" + id + "\"")
		if !node.HeadingSetext {
			r.WriteString(" data-marker=\"#\">")
		} else {
			if 1 == node.HeadingLevel {
				r.WriteString(" data-marker=\"=\">")
			} else {
				r.WriteString(" data-marker=\"-\">")
			}
		}

		if r.Options.HeadingAnchor {
			id := HeadingID(node)
			r.Tag("a", [][]string{{"id", "vditorAnchor-" + id}, {"class", "vditor-anchor"}, {"href", "#" + id}}, false)
			r.WriteString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.Tag("/a", nil, false)
		}

		if !node.HeadingSetext {
			r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--heading"}, {"data-type", "heading-marker"}}, false)
			r.WriteString(strings.Repeat("#", node.HeadingLevel) + " ")
			r.Tag("/span", nil, false)
		}
	} else {
		if node.HeadingSetext {
			r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--heading"}, {"data-type", "heading-marker"}, {"data-render", "2"}}, false)
			r.Newline()
			contentLen := r.setextHeadingLen(node)
			if 1 == node.HeadingLevel {
				r.WriteString(strings.Repeat("=", contentLen))
			} else {
				r.WriteString(strings.Repeat("-", contentLen))
			}
			r.Tag("/span", nil, false)
		}
		r.WriteString("</h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + ">")
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "heading-id"}, {"class", "vditor-ir__marker"}}, false)
		r.WriteString(" {" + string(node.Tokens) + "}")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	tag := "ul"
	if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
		tag = "ol"
	}
	if entering {
		var attrs [][]string
		if node.Tight {
			attrs = append(attrs, []string{"data-tight", "true"})
		}
		if 0 == node.BulletChar {
			if 1 != node.Start {
				attrs = append(attrs, []string{"start", strconv.Itoa(node.Start)})
			}
		}
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", string(node.Marker)})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
		case 3:
			if 0 == node.ListData.BulletChar {
				attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
			} else {
				attrs = append(attrs, []string{"data-marker", string(node.Marker)})
			}
		}
		attrs = append(attrs, []string{"data-block", "0"})
		attrs = append(attrs, []string{"data-node-id", r.NodeID(node)})
		ial := r.NodeAttrs(node)
		if 0 < len(ial) {
			attrs = append(attrs, ial...)
		}
		attrs = append(attrs, []string{"data-type", tag})
		r.renderListStyle(node, &attrs)
		r.Tag(tag, attrs, false)
	} else {
		r.Tag("/"+tag, nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) NodeID(node *ast.Node) (ret string) {
	for _, kv := range node.KramdownIAL {
		if "id" == kv[0] {
			return kv[1]
		}
	}
	return ast.NewNodeID()
}

func (r *VditorIRBlockRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", string(node.Marker)})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
		case 3:
			if 0 == node.ListData.BulletChar {
				attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
			} else {
				attrs = append(attrs, []string{"data-marker", string(node.Marker)})
			}
			if nil != node.FirstChild && nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type /* li.p.task */ ||
				ast.NodeTaskListItemMarker == node.FirstChild.Type /* VditorIRBlockDOM2Tree 忽略了 p */ {
				attrs = append(attrs, []string{"class", r.Options.GFMTaskListItemClass})
			}
		}
		attrs = append(attrs, []string{"data-node-id", r.NodeID(node)})
		ial := r.NodeAttrs(node)
		if 0 < len(ial) {
			attrs = append(attrs, ial...)
		}
		r.Tag("li", attrs, false)
	} else {
		r.Tag("/li", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		if node.TaskListItemChecked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"type", "checkbox"})
		r.Tag("input", attrs, true)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("hr", [][]string{{"data-block", "0"}}, true)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("br", nil, true)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *VditorIRBlockRenderer) renderSpanNode(node *ast.Node) {
	text := r.Text(node)
	var attrs [][]string

	switch node.Type {
	case ast.NodeEmphasis:
		attrs = append(attrs, []string{"data-type", "em"})
	case ast.NodeStrong:
		attrs = append(attrs, []string{"data-type", "strong"})
	case ast.NodeStrikethrough:
		attrs = append(attrs, []string{"data-type", "s"})
	case ast.NodeMark:
		attrs = append(attrs, []string{"data-type", "mark"})
	case ast.NodeSup:
		attrs = append(attrs, []string{"data-type", "sup"})
	case ast.NodeSub:
		attrs = append(attrs, []string{"data-type", "sub"})
	case ast.NodeTag:
		attrs = append(attrs, []string{"data-type", "tag"})
	case ast.NodeLink:
		if 3 != node.LinkType {
			attrs = append(attrs, []string{"data-type", "a"})
		} else {
			attrs = append(attrs, []string{"data-type", "link-ref"})
		}
	case ast.NodeBlockRef:
		attrs = append(attrs, []string{"data-type", "block-ref"})
	case ast.NodeImage:
		attrs = append(attrs, []string{"data-type", "img"})
	case ast.NodeCodeSpan:
		attrs = append(attrs, []string{"data-type", "code"})
	case ast.NodeEmoji:
		attrs = append(attrs, []string{"data-type", "emoji"})
	case ast.NodeInlineMath:
		attrs = append(attrs, []string{"data-type", "inline-math"})
	case ast.NodeHTMLEntity:
		attrs = append(attrs, []string{"data-type", "html-entity"})
	case ast.NodeBackslash:
		attrs = append(attrs, []string{"data-type", "backslash"})
	case ast.NodeInlineHTML:
		attrs = append(attrs, []string{"data-type", "html-inline"})
	case ast.NodeKramdownSpanIAL:
		attrs = append(attrs, []string{"data-type", "span-ial"})
	case ast.NodeBlockQueryEmbedScript:
		attrs = append(attrs, []string{"data-type", "block-query-embed-script"})
		attrs = append(attrs, []string{"class", "vditor-ir__marker vditor-ir__marker--script"})
	default:
		attrs = append(attrs, []string{"data-type", "inline-node"})
	}

	if strings.Contains(text, util.Caret) {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		r.Tag("span", attrs, false)
		return
	}

	preText := node.PreviousNodeText()
	if strings.HasSuffix(preText, util.Caret) {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		r.Tag("span", attrs, false)
		return
	}

	nexText := node.NextNodeText()
	if strings.HasPrefix(nexText, util.Caret) {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		r.Tag("span", attrs, false)
		return
	}

	attrs = append(attrs, []string{"class", "vditor-ir__node"})
	r.Tag("span", attrs, false)
	return
}

func (r *VditorIRBlockRenderer) renderDivNode(node *ast.Node) {
	text := r.Text(node)
	attrs := [][]string{{"data-block", "0"}, {"data-node-id", r.NodeID(node)}}
	ial := r.NodeAttrs(node)
	if 0 < len(ial) {
		attrs = append(attrs, ial...)
	}
	var expand bool
	switch node.Type {
	case ast.NodeCodeBlock:
		attrs = append(attrs, []string{"data-type", "code-block"})
	case ast.NodeSuperBlock:
		attrs = append(attrs, []string{"data-type", "super-block"})
	case ast.NodeMathBlock:
		attrs = append(attrs, []string{"data-type", "math-block"})
	case ast.NodeHTMLBlock:
		attrs = append(attrs, []string{"data-type", "html-block"})
	case ast.NodeYamlFrontMatter:
		attrs = append(attrs, []string{"data-type", "yaml-front-matter"})
	case ast.NodeBlockEmbed:
		attrs = append(attrs, []string{"data-type", "block-ref-embed"})
		text := node.ChildByType(ast.NodeBlockEmbedText)
		tokens := bytes.ReplaceAll(text.Tokens, util.CaretTokens, nil)
		if 0 == len(tokens) {
			attrs = append(attrs, []string{"data-text", "0"})
		}
		id := node.ChildByType(ast.NodeBlockEmbedID)
		expand = bytes.Contains(id.Tokens, util.CaretTokens)
	case ast.NodeBlockQueryEmbed:
		attrs = append(attrs, []string{"data-type", "block-query-embed"})
		script := node.ChildByType(ast.NodeBlockQueryEmbedScript)
		expand = bytes.Contains(script.Tokens, util.CaretTokens)
	}

	if ast.NodeSuperBlock != node.Type {
		if strings.Contains(text, util.Caret) || expand {
			attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		} else {
			attrs = append(attrs, []string{"class", "vditor-ir__node"})
		}
	}
	r.Tag("div", attrs, false)
	return
}

func (r *VditorIRBlockRenderer) Text(node *ast.Node) (ret string) {
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeText, ast.NodeLinkText, ast.NodeLinkDest, ast.NodeLinkSpace, ast.NodeLinkTitle, ast.NodeCodeBlockCode,
				ast.NodeCodeSpanContent, ast.NodeInlineMathContent, ast.NodeMathBlockContent, ast.NodeYamlFrontMatterContent,
				ast.NodeHTMLBlock, ast.NodeInlineHTML, ast.NodeEmojiAlias, ast.NodeBlockRefText, ast.NodeBlockRefSpace,
				ast.NodeBlockEmbedText, ast.NodeBlockEmbedSpace, ast.NodeKramdownSpanIAL:
				ret += string(n.Tokens)
			case ast.NodeCodeBlockFenceInfoMarker:
				ret += string(n.CodeBlockInfo)
			case ast.NodeLink:
				if 3 == n.LinkType {
					ret += string(n.LinkRefLabel)
				}
			}
		}
		return ast.WalkContinue
	})
	return
}

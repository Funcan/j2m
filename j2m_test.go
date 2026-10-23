package j2m

import "testing"

func TestJiraToMD(t *testing.T) {
	tests := []struct {
		name string
		jira string
		want string
	}{
		{"h1", "h1. Title", "# Title"},
		{"h2", "h2. Title", "## Title"},
		{"h3", "h3. Title", "### Title"},
		{"h4", "h4. Title", "#### Title"},
		{"h5", "h5. Title", "##### Title"},
		{"h6", "h6. Title", "###### Title"},
		{"header mid-document", "text\nh2. Title\nmore", "text\n## Title\nmore"},
		{"bold", "*bold*", "**bold**"},
		{"bold in sentence", "some *bold* text", "some **bold** text"},
		{"multiple bolds on a line", "*a* and *b*", "**a** and **b**"},
		{"italic", "_italic_", "*italic*"},
		{"italic in sentence", "some _italic_ text", "some *italic* text"},
		{"snake_case is not italic", "a snake_case_name here", "a snake_case_name here"},
		{"bold and italic", "_*both*_", "***both***"},
		{"monospaced", "{{mono}}", "`mono`"},
		{"monospaced content is untouched", "{{*not bold*}}", "`*not bold*`"},
		{"insert", "+inserted+", "<ins>inserted</ins>"},
		{"superscript", "x^2^", "x<sup>2</sup>"},
		{"subscript", "H~2~O", "H<sub>2</sub>O"},
		{"strikethrough", "a -deleted- b", "a ~~deleted~~ b"},
		{"hyphenated words are not strikethrough", "a well-known non-issue", "a well-known non-issue"},
		{"code block", "{code}\nx := 1\n{code}", "```\nx := 1\n```"},
		{"code block with language", "{code:go}\nx := 1\n{code}", "```go\nx := 1\n```"},
		{"code block with title", "{code:go|title=main.go}\nx := 1\n{code}", "```go\nx := 1\n```"},
		{"code block content is untouched", "{code}\nh1. *x* _y_ [z]\n{code}", "```\nh1. *x* _y_ [z]\n```"},
		{"noformat", "{noformat}\nraw\n{noformat}", "```\nraw\n```"},
		{"noformat content is untouched", "{noformat}\n*x* -y- \n{noformat}", "```\n*x* -y- \n```"},
		{"multiple code blocks", "{code}\na\n{code}\n*b*\n{code}\nc\n{code}", "```\na\n```\n**b**\n```\nc\n```"},
		{"un-named link", "[http://example.com]", "<http://example.com>"},
		{"named link", "[Example|http://example.com]", "[Example](http://example.com)"},
		{"smart link", "Raised here [https://github.com/org/repo/pull/428|https://github.com/org/repo/pull/428|smart-link]", "Raised here [https://github.com/org/repo/pull/428](https://github.com/org/repo/pull/428)"},
		{"smart card", "[Example|http://example.com|smart-card]", "[Example](http://example.com)"},
		{"smart embed", "[Example|http://example.com|smart-embed]", "[Example](http://example.com)"},
		{"image", "!http://example.com/i.png!", "![](http://example.com/i.png)"},
		{"image link", "[!http://example.com/i.png!|http://example.com]", "[![](http://example.com/i.png)](http://example.com)"},
		{"blockquote", "bq. quoted", "> quoted"},
		{"color is removed", "{color:red}warning{color}", "warning"},
		{"un-ordered list", "* a\n* b", "* a\n* b"},
		{"nested un-ordered list", "* a\n** b\n*** c", "* a\n  * b\n    * c"},
		{"ordered list", "# a\n# b", "1. a\n1. b"},
		{"nested ordered list", "# a\n## b\n### c", "1. a\n  1. b\n    1. c"},
		{"formatting inside list", "* _i_ and *b*", "* *i* and **b**"},
		{"panel", "{panel:title=Note}\nbody\n{panel}", "| Note |\n| --- |\n| body |"},
		{"panel after text", "text\n{panel:title=Note}\nbody\n{panel}", "text\n\n| Note |\n| --- |\n| body |"},
		{"table", "||H1||H2||\n|a|b|", "|H1|H2|\n| --- | --- |\n|a|b|"},
		{"table after text", "text\n||H1||H2||\n|a|b|", "text\n\n|H1|H2|\n| --- | --- |\n|a|b|"},
		{"leading blank line before table is kept", "\n||H1||\n|a|", "\n|H1|\n| --- |\n|a|"},
		{"leading blank lines before text are kept", "\n\nhello", "\n\nhello"},
		{"table with leading spaces", "  ||H1||\n  |a|", "|H1|\n| --- |\n|a|"},
		{"url with insert chars", "See https://example.com/display/Space/2024-Feb-08+Gmail+account+locked for details", "See https://example.com/display/Space/2024-Feb-08+Gmail+account+locked for details"},
		{"url with strikethrough chars", "a -b https://example.com/x- c", "a -b https://example.com/x- c"},
		{"url with italic chars", "see https://example.com/_private_/page", "see https://example.com/_private_/page"},
		{"url with superscript chars", "see https://example.com/?q=a^b^c", "see https://example.com/?q=a^b^c"},
		{"url with subscript chars", "see https://example.com/~user/~x", "see https://example.com/~user/~x"},
		{"un-named link url with insert chars", "[https://example.com/a+b+c]", "<https://example.com/a+b+c>"},
		{"named link url with insert chars", "[Page|https://example.com/a+b+c]", "[Page](https://example.com/a+b+c)"},
		{"formatting next to url", "*bold* https://example.com/a+b+c _it_", "**bold** https://example.com/a+b+c *it*"},
		{"plain text", "nothing to see here", "nothing to see here"},
		{"literal < is escaped", "* <add> Fix Vault", "* \\<add> Fix Vault"},
		{"literal < in comparison is escaped", "a < b", "a \\< b"},
		{"literal < in monospaced is not escaped", "{{<add>}}", "`<add>`"},
		{"literal < in code block is not escaped", "{code}\n<add>\n{code}", "```\n<add>\n```"},
		{"run of carets is not superscript", "^^^ will be broken into tickets", "^^^ will be broken into tickets"},
		{"carets around spaced text are not superscript", "^^^ see above ^^^", "^^^ see above ^^^"},
		{"run of pluses is not insert", "+++ will be broken into tickets", "+++ will be broken into tickets"},
		{"pluses around spaced text are not insert", "1 + 2 + 3", "1 + 2 + 3"},
		{"run of tildes is not subscript", "~~~ will be broken into tickets", "~~~ will be broken into tickets"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JiraToMD(tt.jira); got != tt.want {
				t.Errorf("JiraToMD(%q) = %q, want %q", tt.jira, got, tt.want)
			}
		})
	}
}

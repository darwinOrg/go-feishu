package dgfeishu

import "testing"

func TestFeishuFileType(t *testing.T) {
	cases := map[string]string{
		"a.pdf":  "pdf",
		"a.DOCX": "doc",
		"a.xlsx": "xls",
		"a.pptx": "ppt",
		"a.mp4":  "mp4",
		"a.opus": "opus",
		"a.zip":  "stream",
	}
	for file, want := range cases {
		if got := feishuFileType(file); got != want {
			t.Errorf("feishuFileType(%q) = %q, want %q", file, got, want)
		}
	}
}

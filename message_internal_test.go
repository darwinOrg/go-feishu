package dgfeishu

import (
	"encoding/json"
	"testing"
)

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

func TestBuildTableCard(t *testing.T) {
	card, err := buildTableCard("标题", []string{"客户", "金额"}, []string{"120px", ""}, [][]any{{"A", "1"}, {"B"}})
	if err != nil {
		t.Fatal(err)
	}

	var got struct {
		Header   map[string]any `json:"header"`
		Elements []struct {
			Tag     string           `json:"tag"`
			Columns []map[string]any `json:"columns"`
			Rows    []map[string]any `json:"rows"`
		} `json:"elements"`
	}
	if err := json.Unmarshal([]byte(card), &got); err != nil {
		t.Fatal(err)
	}
	if got.Header["template"] != "blue" {
		t.Errorf("header template = %v", got.Header["template"])
	}
	table := got.Elements[0]
	if table.Tag != "table" || len(table.Columns) != 2 {
		t.Fatalf("table tag/columns = %s/%d", table.Tag, len(table.Columns))
	}
	if table.Columns[1]["display_name"] != "金额" {
		t.Errorf("column 1 display_name = %v", table.Columns[1]["display_name"])
	}
	if table.Columns[0]["width"] != "120px" || table.Columns[1]["width"] != "auto" {
		t.Errorf("column widths = %v/%v", table.Columns[0]["width"], table.Columns[1]["width"])
	}
	if table.Rows[0]["col_0"] != "A" || table.Rows[0]["col_1"] != "1" {
		t.Errorf("row 0 = %v", table.Rows[0])
	}
	if v, ok := table.Rows[1]["col_1"]; ok && v != "" {
		t.Errorf("short row col_1 = %v, want empty", v)
	}
}

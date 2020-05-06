package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kjk/notionapi"
	"github.com/yiningv/nblog/pub/log"
	"strconv"
)

const (
	siteConfigPageId   = "159280121cbc44ecb7ef074534a8897f"
	sourceConfigPageId = "926a2ae5b91d45d985961743d230c4c2"
)

var (
	ErrNoTableView = errors.New("页面上没有TableView")
)

func main() {
	//GetSiteConfig()
	//GetSourceConfig()
	GetConfigFromNotion(sourceConfigPageId, nil)
}

// 从notion上获取站点配置
func GetSiteConfig() {
	c := notionapi.Client{}
	page, err := c.DownloadPage(siteConfigPageId)
	if err != nil {
		log.Error(fmt.Sprintf("GetConfigFromNotion error: %v", err))
		return
	}
	if len(page.TableViews) == 0 {
		log.Info(fmt.Sprintf("页面上没有TableView"))
		return
	}
	// 每个页面上只会读取第一个默认的table
	tv := page.TableViews[0]
	headInfos := tv.Columns
	rows := tv.Rows
	for i := range rows {
		siteMap := make(map[string]interface{})
		row := rows[i]
		columns := row.Columns
		for j := range columns {
			column := columns[j]
			if len(column) == 0 {
				continue
			}
			headInfo := headInfos[j]
			if headInfo == nil || headInfo.Schema == nil {
				continue
			}
			value := ""
			switch headInfo.Type() {
			case "text":
				for _, columnInfo := range column {
					value += columnInfo.Text
				}
			default:
				columnInfo := column[0]
				if columnInfo.IsPlain() {
					value = columnInfo.Text
				} else {
					value = columnInfo.Attrs[0][1]
				}
			}
			key := headInfo.Name()
			if key == "image" {
				key = "value"
			}
			siteMap[key] = value
		}
	}
}

type KV struct {
	K string `json:"k"`
	V string `json:"v"`
}

type Date struct {
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Type string `json:"type"`
	TimeZone string `json:"time_zone"`
}

type DealTableFn func(hColumns []*notionapi.ColumnInfo, tRow *notionapi.TableRow)

func GetConfigFromNotion(pageId string, fn DealTableFn) (result []map[string]interface{}, err error) {
	c := notionapi.Client{}
	var page *notionapi.Page
	page, err = c.DownloadPage(pageId)
	if err != nil {
		log.Error(fmt.Sprintf("GetConfigFromNotion for pageId(%s) error: %v", pageId, err))
		return
	}
	if len(page.TableViews) == 0 {
		err = ErrNoTableView
		log.Info(fmt.Sprintf("%v", err))
		return
	}
	// 每个页面上只会读取第一个默认的table
	tv := page.TableViews[0]
	hColumns := tv.Columns
	tRows := tv.Rows
	for i := range tRows {
		siteMap := make(map[string]interface{})
		tRow := tRows[i]
		tColumns := tRow.Columns
		for j := range tColumns {
			hColumn := hColumns[j]
			if hColumn == nil || hColumn.Schema == nil {
				continue
			}
			tColumn := tColumns[j]
			if len(tColumn) == 0 {
				continue
			}
			value := ""
			// person、rollup、formula
			// 不支持这几个类型，跳过
			switch hColumn.Type() {
			case "title":
			case "url":
			case "email":
			case "phone_number":
			case "select":
			case "number":
			case "multi_select":
				value = tColumn[0].Text
			case "text":
				//如果有这个字符‣
				if tColumn[0].Text == "‣" {
					if !tColumn[0].IsPlain() {
						value = tColumn[0].Attrs[0][1]
					}
				} else {
					// 如果字体加了修饰，会有多列
					for _, colInfo := range tColumn {
						value += colInfo.Text
					}
				}
			case "checkbox":
				value = strconv.FormatBool(tColumn[0].Text == "Yes")
			case "date":
				dateStr := tColumn[0].Attrs[0][1]
				var date *Date
				err = json.Unmarshal([]byte(dateStr), date)
				if err != nil {
					return
				}
				var marshal []byte
				marshal, err = json.Marshal(date)
				if err != nil {
					return
				}
				value = string(marshal)
			case "file":
				//kv := make(map[string]string)
				list := make([]*KV, len(tColumn))
				for i, colInfo := range tColumn {
					if !colInfo.IsPlain() {
						kv := &KV{
							K: colInfo.Text,
							V: colInfo.Attrs[0][1],
						}
						list[i] = kv
					}
				}
				var marshal []byte
				marshal, err = json.Marshal(list)
				if err != nil {
					return
				}
				value = string(marshal)
			default:
				continue
			}
			key := hColumn.Name()
			if key == "image" {
				key = "value"
			}
			siteMap[key] = value
		}
	}
	return
}

type SourceConfig struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Table string `json:"table"`
	Des   string `json:"des"`
}

// 获取资源配置
func GetSourceConfig() {
	c := notionapi.Client{}
	page, err := c.DownloadPage(sourceConfigPageId)
	errorPanic(err)
	if len(page.TableViews) == 0 {
		panic(errors.New("没有内容搞个屁"))
	}
	// 每个页面上只会读取第一个默认的table
	tv := page.TableViews[0]
	headInfos := tv.Columns
	rows := tv.Rows
	args := make([]*SourceConfig, 0)
	for i := 0; i < len(rows); i++ {
		sourceConfig := new(SourceConfig)
		sorceMap := make(map[string]interface{})
		row := rows[i]
		columns := row.Columns
		for j := 0; j < len(columns); j++ {
			column := columns[j]
			if len(column) == 0 {
				continue
			}
			headInfo := headInfos[j]

			if headInfo.Schema == nil {
				continue
			}
			value := ""
			columnInfo := column[0]
			if columnInfo.IsPlain() {
				value = columnInfo.Text
			} else {
				value = columnInfo.Attrs[0][1]
			}
			sorceMap[headInfo.Name()] = value
		}
		marshal, err := json.Marshal(sorceMap)
		if err != nil {
			continue
		}
		err = json.Unmarshal(marshal, sourceConfig)
		if err != nil {
			continue
		}
		args = append(args, sourceConfig)
	}
	fmt.Println()
}

func errorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

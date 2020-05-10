package notion

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kjk/notionapi"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"strconv"
	"time"
)

// 从notion上获取站点配置
func GetSiteConfig() (sc map[string]*model.SiteConfig, err error) {
	var tableDatas []map[string]interface{}
	tableDatas, err = getConfigFromNotion(conf.Conf.App.SiteConfigPageId)
	if err != nil {
		return
	}
	var data []byte
	data, err = json.Marshal(tableDatas)
	if err != nil {
		return
	}
	var siteConfigs []*model.SiteConfig
	err = json.Unmarshal(data, &siteConfigs)
	if err != nil {
		return
	}
	sc = map[string]*model.SiteConfig{}
	// 去掉name重复的数据
	for i := range siteConfigs {
		c := siteConfigs[i]
		if c.Name == "" {
			continue
		}
		sc[c.Name] = c
	}
	return
}

// 获取资源配置
func GetSourceConfig() (sc map[string]*model.SourceConfig, err error) {
	var tableDatas []map[string]interface{}
	tableDatas, err = getConfigFromNotion(conf.Conf.App.SourceConfigPageId)
	if err != nil {
		return
	}
	var data []byte
	data, err = json.Marshal(tableDatas)
	if err != nil {
		return
	}
	var sourceConfigs []*model.SourceConfig
	err = json.Unmarshal(data, &sourceConfigs)
	if err != nil {
		return
	}
	sc = map[string]*model.SourceConfig{}
	// 去掉name重复的数据
	for i := range sourceConfigs {
		c := sourceConfigs[i]
		if c.Name == "" {
			continue
		}
		c.OrderNum = i
		sc[c.Name] = c
	}
	return
}

func GetPosts() (postsMap map[string]*model.Posts, err error) {
	var tableDatas []map[string]interface{}
	tableDatas, err = getConfigFromNotion(conf.Conf.App.PostsPageId)
	if err != nil {
		return
	}
	var data []byte
	data, err = json.Marshal(tableDatas)
	if err != nil {
		return
	}
	var posts []*model.Posts
	err = json.Unmarshal(data, &posts)
	if err != nil {
		return
	}
	postsMap = make(map[string]*model.Posts)
	for i := range posts {
		p := posts[i]
		if p.Title == "" {
			continue
		}
		var pTime time.Time
		failed := false
		if p.PublishedTime != "" {
			date := &model.Date{}
			err := json.Unmarshal([]byte(p.PublishedTime), date)
			if err != nil {
				failed = true
				break
			}
			location, err := time.LoadLocation(date.TimeZone)
			if err != nil {
				failed = true
				break
			}
			startTime := "00:00"
			if date.StartTime != "" {
				startTime = date.StartTime
			}
			timeStr := fmt.Sprintf("%s %s", date.StartDate, startTime)
			pTime, err = time.ParseInLocation("2006-01-02 15:04", timeStr, location)
			if err != nil {
				failed = true
				break
			}
		}
		if failed {
			p.PTime = time.Date(2015, 10, 6, 22, 57, 0, 0, time.Local)
		} else {
			p.PTime = pTime
		}

		// 暂时没有找到slug很好的方案，先用pageId代替
		p.Slug = p.PageId
		postsMap[p.PageId] = p
	}
	return
}

func GetStatic(pageId string) (st []*model.Static, err error) {
	var tableDatas []map[string]interface{}
	tableDatas, err = getConfigFromNotion(conf.Conf.App.SourceConfigPageId)
	if err != nil {
		return
	}
	var data []byte
	data, err = json.Marshal(tableDatas)
	if err != nil {
		return
	}
	var sourceConfigs []*model.Static
	err = json.Unmarshal(data, &sourceConfigs)
	if err != nil {
		return
	}

	return
}

func GetHTML(pageId string) {

}

func getConfigFromNotion(pageId string) (result []map[string]interface{}, err error) {
	c := notionapi.Client{}
	var page *notionapi.Page
	page, err = c.DownloadPage(pageId)
	if err != nil {
		log.Error(fmt.Sprintf("getConfigFromNotion for pageId(%s) error: %v", pageId, err))
		return
	}
	if len(page.TableViews) == 0 {
		err = errors.New("页面上没有TableView")
		log.Info(fmt.Sprintf("%v", err))
		return
	}
	result = make([]map[string]interface{}, 0)
	// 每个页面上只会读取第一个默认的table
	tv := page.TableViews[0]
	hColumns := tv.Columns
	tRows := tv.Rows
	for i := range tRows {
		rowMap := make(map[string]interface{})
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
			// 仅处理要支持的类型，其他类型跳过
			switch hColumn.Type() {
			case "title":
				fallthrough
			case "url":
				fallthrough
			case "email":
				fallthrough
			case "phone_number":
				fallthrough
			case "select":
				fallthrough
			case "number":
				fallthrough
			case "multi_select":
				value = tColumn[0].Text
			case "text":
				//如果有这个字符‣
				// 内容为Link
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
				date := &model.Date{}
				err = json.Unmarshal([]byte(dateStr), date)
				if err != nil {
					log.Error(fmt.Sprintf("json.Unmarshal failed: %v", err))
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
				list := make([]*model.KV, len(tColumn))
				for i, colInfo := range tColumn {
					if !colInfo.IsPlain() {
						kv := &model.KV{
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
			}
			key := hColumn.Name()
			rowMap[key] = value
		}
		if len(rowMap) != 0 {
			rowMap["last_edited_time"] = tRow.Page.LastEditedTime
			rowMap["page_id"] = notionapi.ToNoDashID(tRow.Page.ID)
			result = append(result, rowMap)
		}
	}
	return
}

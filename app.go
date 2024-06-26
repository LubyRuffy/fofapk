package main

import (
	"context"
	"fmt"
	"github.com/LubyRuffy/fofapk/pkg/models"
	"github.com/LubyRuffy/gofofa"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var (
	fofaKey   string
	fetchSize = 2000
)

type Result struct {
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Error string      `json:"error"`
}

func NewErrorResult(msg string) *Result {
	return &Result{
		Code:  500,
		Error: msg,
	}
}

func NewDataResult(data any) *Result {
	return &Result{
		Code: 200,
		Data: data,
	}
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	//runtime.EventsOn(ctx, "diff", func(optionalData ...interface{}) {
	//})
}

func queryToMap(c *gofofa.Client, query string, taskid string) map[string]*models.Result {
	fields := models.Result{}.Fields()
	res1, err := c.HostSearch(query, fetchSize, fields)
	if err != nil {
		panic(err)
	}
	sMap := make(map[string]*models.Result)
	for _, r := range res1 {
		result := &models.Result{
			UUID: taskid,
		}
		for i, f := range fields {
			switch f {
			case "host":
				result.Host = r[i]
			case "ip":
				result.IP = r[i]
			case "port":
				result.Port = r[i]
			case "as_organization":
				result.ASOrganization = r[i]
			case "protocol":
				result.Protocol = r[i]
			case "domain":
				result.Domain = r[i]
			case "certs_subject_cn":
				result.CertsSubjectCN = r[i]
			case "title":
				result.Title = r[i]
			case "fid":
				result.FID = r[i]
			default:
				panic("unknown field:" + f)
			}
		}
		sMap[result.IP] = result
	}
	return sMap
}

func diffMaps(m1, m2 map[string]*models.Result) (diff, same map[string]*models.Result) {
	diff = make(map[string]*models.Result)
	same = make(map[string]*models.Result)

	// 遍历m1并检查m2中是否存在相同的键值对
	for k1, v1 := range m1 {
		_, ok := m2[k1]
		if !ok {
			// 如果m2中不存在此键或值不同，则认为是m1中移除的项
			diff[k1] = v1
			diff[k1].From = models.FromA
		} else {
			same[k1] = v1
		}
	}

	// 遍历m2并检查m1中是否存在相同的键值对
	for k2, v2 := range m2 {
		_, ok := m1[k2]
		if !ok {
			// 如果m1中不存在此键，则认为是m2中新添加的项
			diff[k2] = v2
			diff[k2].From = models.FromB
		} else {
			same[k2] = v2
		}
	}

	return diff, same
}

func (a *App) emitProgress(progress float64, logs []string, finished ...bool) {
	runtime.EventsEmit(a.ctx, "onProgress", map[string]interface{}{
		"progress": progress,
		"finished": len(finished) > 0,
		"logs":     logs,
	})
}

func (a *App) emitError(err string) {
	runtime.EventsEmit(a.ctx, "onError", map[string]interface{}{
		"error": err,
	})
}

func (a *App) loadData(taskid string) []models.Result {
	var dbData []models.Result
	if err := models.Get().Where(&models.Result{
		UUID: taskid,
	}).Where("done = ? and `from` != ?", false, models.FromBoth).Find(&dbData).Error; err != nil {
		a.emitError(fmt.Sprintf("receive results from db failed: %v", err))
	}

	return dbData
}

func defaultFofaClient() (*gofofa.Client, error) {
	//var opts []gofofa.ClientOption
	//if fofaKey != "" {
	//	opts = append(opts, gofofa.WithURL(fmt.Sprintf("https://fofa.info/api/?key=%s&version=v2&debuglevel=0", fofaKey)))
	//}
	//c, err := gofofa.NewClient(opts...)
	c, err := gofofa.NewClient()
	if fofaKey != "" {
		c.Key = fofaKey
	}

	return c, err
}

// SplitSlice splits any slice into a slice of slices with the given chunk size.
func SplitSlice[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func (a *App) StartTask(q1, q2 string) *Result {
	t := models.NewTask(q1, q2)
	taskid := t.UUID
	go func(t *models.Task) {
		defer func() {
			if err := recover(); err != nil {
				a.emitError(fmt.Sprintf("%v", err))
			} else {
				a.emitProgress(100.00, []string{"finished"}, true)
			}
		}()

		a.emitProgress(10.00, []string{"initial fofa connection", "current task id is:" + taskid})

		c, err := defaultFofaClient()
		if err != nil {
			panic(err)
		}

		a.emitProgress(11.00, []string{"try to fetch data of query1 from fofa"})

		s1Map := queryToMap(c, q1, taskid)
		a.emitProgress(40.00, []string{fmt.Sprintf("fetch data of query1 from fofa finished, %d ips", len(s1Map)),
			"now try to fetch data of query2 from fofa"})

		s2Map := queryToMap(c, q2, taskid)
		a.emitProgress(70.00,
			[]string{fmt.Sprintf("fetch data of query2 from fofa finished, %d ips", len(s2Map)),
				"now try to diff two results"})

		diff, same := diffMaps(s1Map, s2Map)
		a.emitProgress(80.00, []string{"now try to save results"})

		var results []*models.Result
		for _, m := range []map[string]*models.Result{diff, same} {
			for _, result := range m {
				//log.Println(ip)
				results = append(results, result)
			}
		}
		rand.New(rand.NewSource(time.Now().UnixNano())).
			Shuffle(len(results), func(i, j int) {
				results[i], results[j] = results[j], results[i]
			})

		subSlices := SplitSlice(results, 500)
		for _, sl := range subSlices {
			err = models.Get().Transaction(func(tx *gorm.DB) error {
				return tx.Create(sl).Error
			})
			if err != nil {
				a.emitError(fmt.Sprintf("save db failed: %v", err))
			}
		}

		a.emitProgress(85.00, []string{"now try to receive results"})

		runtime.EventsEmit(a.ctx, "onData", map[string]interface{}{
			"data":     a.loadData(taskid),
			"size1":    len(s1Map),
			"size2":    len(s2Map),
			"diffSize": len(diff),
			"logs":     []string{fmt.Sprintf("found %d equals, %d differents", len(same), len(diff))},
		})

	}(t)
	return NewDataResult(t)
}

func (a *App) UpdateScore(taskid string, ips []string, score int) *Result {
	for _, ip := range ips {
		var r models.Result
		if err := models.Get().Where(&models.Result{
			UUID: taskid,
			IP:   ip,
		}).First(&r).Error; err != nil {
			return NewErrorResult(err.Error())
		}
		if err := models.Get().Model(&r).Updates(&models.Result{
			Score: score,
			Done:  true,
		}).Error; err != nil {
			return NewErrorResult(err.Error())
		}
	}

	var defaultInt int
	var scoreA *int
	if err := models.Get().Model(&models.Result{}).Select("sum(score)").Where(&models.Result{
		UUID: taskid,
		Done: true,
		From: models.FromA,
	}).Scan(&scoreA).Error; err != nil {
		return NewErrorResult(err.Error())
	}
	if scoreA == nil {
		scoreA = &defaultInt
	}

	var scoreB *int
	if err := models.Get().Model(&models.Result{}).Select("sum(score)").Where(&models.Result{
		UUID: taskid,
		Done: true,
		From: models.FromB,
	}).Scan(&scoreB).Error; err != nil {
		return NewErrorResult(err.Error())
	}
	if scoreB == nil {
		scoreB = &defaultInt
	}

	var diffSize int
	if err := models.Get().Debug().Model(&models.Result{}).Select("count(*)").Where(&models.Result{
		UUID: taskid,
	}).Where("`done`=? and (`from`=? or `from`=?)", false, models.FromA, models.FromB).Scan(&diffSize).Error; err != nil {
		return NewErrorResult(err.Error())
	}

	return NewDataResult(map[string]interface{}{
		"score1":   *scoreA,
		"score2":   *scoreB,
		"data":     a.loadData(taskid),
		"diffSize": diffSize,
		"logs":     []string{fmt.Sprintf("found %d differents", diffSize)},
	})
}

func (a *App) FofaStat(ip string) *Result {
	defer func() {
		if err := recover(); err != nil {
			a.emitError(fmt.Sprintf("%v", err))
		}
	}()
	c, err := defaultFofaClient()
	if err != nil {
		panic(err)
	}
	res, err := c.Stats("ip="+ip, 5, []string{"domain", "title", "certs_subject_cn", "fid"})
	if err != nil {
		panic(err)
	}

	return NewDataResult(res)
}

func (a *App) UpdateConfig(fofaKeyString string, fetchSizeInt int) *Result {
	fofaKey = fofaKeyString
	fetchSize = fetchSizeInt
	return NewDataResult(true)
}

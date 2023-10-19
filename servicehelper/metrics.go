package servicehelper

/**
 * 统计处理，对go-metrics封装
 */

import (
	"github.com/rcrowley/go-metrics"
	"log"
	"os"
	"sync"
	"time"
)

const (
	METRICS_TYPE_COUNTER    = "Counter"    // 纯计数类型
	METRICS_TYPE_METERS     = "Meters"     // 单位时间发生次数，最近1,5,15分钟滑动平均
	METRICS_TYPE_GAUGES     = "Gauges"     // 以时间为单位瞬时值
	METRICS_TYPE_HISTOGRAMS = "Histograms" // 历史信息，Count，Min, Max, Mean, Median, 75%, 95%, 99%
	METRICS_TYPE_TIMERS     = "Timers"     // 对某个代码模块同时进行统计调用频率以及调用耗时统计
)

var (
	allMetricsKeys = &MetricsAllKeys{
		list: make(map[string]string),
	}
	allMetricsSamples = &MetricsSamples{
		list: make(map[string]metrics.Sample),
	}
)

// 存储所有统计数据的key
type MetricsAllKeys struct {
	mutex sync.Mutex
	list  map[string]string
}

// 统计history的样本存储器
type MetricsSamples struct {
	mutex sync.Mutex
	list  map[string]metrics.Sample
}

// 初始化metrics执行逻辑
func InitMetrics() {
	// 每五分钟输出一次
	go metrics.Log(metrics.DefaultRegistry, 5*time.Minute, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
}

// 所有key的信息获取
func MetricsGetAllKeys() map[string]string {
	ret := make(map[string]string)
	allMetricsKeys.mutex.Lock()
	defer allMetricsKeys.mutex.Unlock()
	for key, value := range allMetricsKeys.list {
		ret[key] = value
	}
	return ret
}

// 新增key信息
func addMetricsKey(name string, typeStr string) {
	allMetricsKeys.mutex.Lock()
	item, ok := allMetricsKeys.list[name]
	if ok {
		if item != typeStr {
			allMetricsKeys.list[name] = typeStr
		}
	} else {
		allMetricsKeys.list[name] = typeStr
	}
	allMetricsKeys.mutex.Unlock()
}

// //////////////////// 计数相关 //////////////////////////
func MetricsInc(key string, val int64) {
	addMetricsKey(key, METRICS_TYPE_COUNTER)
	metrics.GetOrRegisterCounter(key, nil).Inc(val)
}

func MetricsDec(key string, val int64) {
	metrics.GetOrRegisterCounter(key, nil).Dec(val)
}

func MetricsCounterGet(key string) int64 {
	return metrics.GetOrRegisterCounter(key, nil).Count()
}

/////////////////// Gauge 瞬时值信息 /////////////////////

func MetricsGaugeGet(key string) int64 {
	return metrics.GetOrRegisterGauge(key, nil).Value()
}

func MetricsGaugeUpdate(key string, val int64) {
	addMetricsKey(key, METRICS_TYPE_GAUGES)
	metrics.GetOrRegisterGauge(key, nil).Update(val)
}

// ///////////////// meters 速率计算 /////////////////////
// 用于计算一段时间内的计量，通常用于计算接口调用频率，
// 如QPS(每秒的次数)，主要分为rateMean,Rate1/Rate5/Rate15等指标．
func MetricsMeterMark(key string, val int64) {
	addMetricsKey(key, METRICS_TYPE_METERS)
	metrics.GetOrRegisterMeter(key, nil).Mark(val)
}

func MetricsMeterGetCount(key string) int64 {
	return metrics.GetOrRegisterMeter(key, nil).Snapshot().Count()
}

func MetricsMeterGetRate1(key string) float64 {
	return metrics.GetOrRegisterMeter(key, nil).Snapshot().Rate1()
}

func MetricsMeterGetRate5(key string) float64 {
	return metrics.GetOrRegisterMeter(key, nil).Snapshot().Rate5()
}

func MetricsMeterGetRate15(key string) float64 {
	return metrics.GetOrRegisterMeter(key, nil).Snapshot().Rate15()
}

func MetricsMeterGetRateMean(key string) float64 {
	return metrics.GetOrRegisterMeter(key, nil).Snapshot().RateMean()
}

/////////////////// Histograms 速率计算 /////////////////////
// 主要用于对数据集中的值分布情况进行统计，典型的应用场景为接口耗时，
// 接口每次调用都会产生耗时，记录每次调用耗时来对接口耗时情况进行分析显然不现实．
// 因此将接口一段时间内的耗时看做数据集，
// 并采集Count，Min, Max, Mean, Median, 75%, 95%, 99%等指标．
// 以相对较小的资源消耗，来尽可能反应数据集的真实情况．

// 从样本数据中获取信息
func getMetricsHistogramsSample(key string) metrics.Sample {
	allMetricsSamples.mutex.Lock()
	defer allMetricsSamples.mutex.Unlock()
	item, ok := allMetricsSamples.list[key]
	if ok {
		return item
	} else {
		s := metrics.NewExpDecaySample(4096, 0.015)
		allMetricsSamples.list[key] = s
		return s
	}
}

// 更新历史型样本
func MetricsHistogramsUpdate(key string, val int64) {
	addMetricsKey(key, METRICS_TYPE_HISTOGRAMS)
	s := getMetricsHistogramsSample(key)
	metrics.GetOrRegisterHistogram(key, nil, s).Update(val)
}

func MetricsHistogramsGetCount(key string) int64 {
	return metrics.GetOrRegisterHistogram(key, nil, getMetricsHistogramsSample(key)).
		Snapshot().Count()
}

func MetricsHistogramsGetMin(key string) int64 {
	return metrics.GetOrRegisterHistogram(key, nil, getMetricsHistogramsSample(key)).
		Snapshot().Min()
}

func MetricsHistogramsGetMax(key string) int64 {
	return metrics.GetOrRegisterHistogram(key, nil, getMetricsHistogramsSample(key)).
		Snapshot().Max()
}

func MetricsHistogramsGetMean(key string) float64 {
	return metrics.GetOrRegisterHistogram(key, nil, getMetricsHistogramsSample(key)).
		Snapshot().Mean()
}

// 比如 []float64{0.5, 0.75, 0.95, 0.99}, 0->0.5 1->0.75 ....
func MetricsHistogramsGetPercentiles(key string, percentiles []float64) []float64 {
	return metrics.GetOrRegisterHistogram(key, nil, getMetricsHistogramsSample(key)).
		Snapshot().Percentiles(percentiles)
}

/////////////////// Timer 速率计算 /////////////////////
// 对某个代码模块同时进行统计调用频率以及调用耗时统计．
// 指标就是Histograms以及Meters两种统计方式的合集
// 由于timer需要封装func 进入参数，所以暂时不提供此类函数的封装，将来有需要可以增加封装

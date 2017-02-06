package sensorbee_plugin_psutil

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"time"
)

type psutilIntervalSource struct {
	ctx *core.Context
	w   core.Writer

	interval time.Duration
}

func (s *psutilIntervalSource) GenerateStream(ctx *core.Context, w core.Writer) error {
	var cnt int64
	for ; ; cnt++ {
		v, _ := mem.VirtualMemory()

		virtualMemory := data.Map{
			"total":        data.Int(v.Total),
			"available":    data.Int(v.Available),
			"used_percent": data.Float(v.UsedPercent),
			"free":         data.Int(v.Free),
		}

		ps, _ := cpu.Percent(0, false)
		percentOverAll := data.Float(ps[0])
		ps, _ = cpu.Percent(0, true)

		percentsPercpu := []interface{}{}
		for _, p := range ps {
			percentsPercpu = append(percentsPercpu, data.Float(p))
		}
		percentsPercpuDataArray, _ := data.NewArray(percentsPercpu)

		cpuInfo, _ := cpu.Info()
		cpuInfoArray := []interface{}{}
		for _, info := range cpuInfo {
			infoMap := make(map[string]interface{})
			infoMap["cpu"] = info.CPU
			infoMap["vendor_id"] = info.VendorID
			infoMap["family"] = info.Family
			infoMap["model"] = info.Model
			cpuInfoArray = append(cpuInfoArray, infoMap)
		}
		cpuInfoDataArray, _ := data.NewArray(cpuInfoArray)

		counters, _ := disk.IOCounters()
		countersArray := []interface{}{}
		for name, counter := range counters {
			counterMap := make(map[string]interface{})
			counterMap["read_count"] = counter.ReadCount
			counterMap["write_count"] = counter.WriteCount

			diskInfo := make(map[string]interface{})
			diskInfo[name] = counterMap

			counterDataMap, _ := data.NewMap(counterMap)
			countersArray = append(countersArray, counterDataMap)
		}
		countersDataArray, _ := data.NewArray(countersArray)

		mem := data.Map{
			"virtual_memory": virtualMemory,
		}
		cpu := data.Map{
			"percent_overall": percentOverAll,
			"percents_percpu": percentsPercpuDataArray,
			"info":            cpuInfoDataArray,
		}
		disk := data.Map{
			"counters": countersDataArray,
		}

		tuple := core.NewTuple(data.Map{"mem": mem, "cpu": cpu, "disk": disk})
		if err := w.Write(ctx, tuple); err != nil {
			return err
		}
		time.Sleep(s.interval)
	}
	return nil
}

func (s *psutilIntervalSource) Stop(ctx *core.Context) error {
	return nil
}

func NewIntervalSource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Source, error) {
	interval := 1 * time.Second
	if v, ok := params["interval"]; ok {
		i, err := data.ToDuration(v)
		if err != nil {
			return nil, err
		}
		interval = i
	}

	return core.ImplementSourceStop(&psutilIntervalSource{
		interval: interval,
	}), nil
}

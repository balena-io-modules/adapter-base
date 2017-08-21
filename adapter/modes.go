package adapter

type scanMode func(*ScanOptions, *Worker, chan Job, Job)
type updateMode func(*UpdateOptions, *Worker, chan Job, Job)

var scanModes map[string]scanMode
var updateModes map[string]updateMode

func init() {
	scanModes = map[string]scanMode{
		"simulate": simulateScan,
	}

	updateModes = map[string]updateMode{
		"simulate": simulateUpdate,
	}
}

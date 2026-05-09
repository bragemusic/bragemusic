package types

type JobType string

const (
	JobImporterRun                 JobType = "importer:run"
	JobMetaSyncRun                 JobType = "metasync:run"
	JobAuthClientServerStatus      JobType = "authclient:serverstatus"
	JobSyncerDaemon                JobType = "syncer:daemon"
	JobMediaManagerUpdateSyncItems JobType = "mediamanager:update-searchitems"
	JobAnalyserRunTrackAnalysis    JobType = "analyser:run-track-analysis"
)

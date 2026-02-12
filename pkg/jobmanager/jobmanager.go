package jobmanager

import (
	"context"
	"log/slog"
	"time"

	"github.com/bragemusic/core/pkg/importer"
	"github.com/bragemusic/core/pkg/mediamanager"
	"github.com/bragemusic/core/pkg/metasyncer"
	"github.com/bragemusic/core/pkg/types"
)

type jobDefinition struct {
	interval time.Duration
	run      func(context.Context) error
}

type JobManager struct {
	log      *slog.Logger
	mediamgr *mediamanager.MediaManager
	importer *importer.Importer
	metasync *metasyncer.MetaSyncer
	jobs     map[types.JobType]jobDefinition
}

func (j *JobManager) startJob(ctx context.Context, jobType types.JobType, job jobDefinition) {
	ticker := time.NewTicker(job.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := job.run(ctx); err != nil {
				j.log.ErrorContext(ctx, "job failed",
					"job", jobType,
					"error", err,
				)
			}
		}
	}
}

func (j *JobManager) StartScheduler(ctx context.Context) {
	j.log.InfoContext(ctx, "starting scheduler")

	for jobType, job := range j.jobs {
		go j.startJob(ctx, jobType, job)
	}

	<-ctx.Done()
	j.log.InfoContext(ctx, "jobs finished")
}

func New(slogHandler slog.Handler, m *mediamanager.MediaManager, i *importer.Importer, ms *metasyncer.MetaSyncer) JobManager {
	jobs := map[types.JobType]jobDefinition{
		types.JobImporterRun: {
			interval: 10 * time.Second,
			run:      i.Run,
		},
		types.JobMetaSyncRun: {
			interval: 15 * time.Second,
			run:      ms.Sync,
		},
	}

	return JobManager{
		log:      slog.New(slogHandler).With("service", "job-manager"),
		mediamgr: m,
		importer: i,
		metasync: ms,
		jobs:     jobs,
	}
}

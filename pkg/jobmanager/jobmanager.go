package jobmanager

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/types"
)

type JobConfig struct {
	ImporterRunTiming   int
	MetaSyncerRunTiming int
}

type JobDefinition struct {
	Type     types.JobType
	Interval time.Duration
	Run      func(context.Context) error
	c        chan struct{}
}

type JobManager struct {
	log  *slog.Logger
	jobs []JobDefinition
	berr bragerr.BragErrFactory
}

func (j *JobManager) RunJob(ctx context.Context, jobType types.JobType) error {
	jobFound := false
	for _, job := range j.jobs {
		if job.Type == jobType {
			job.c <- struct{}{}
			jobFound = true
		}
	}

	if !jobFound {
		return j.berr.JobTypeMissing(errors.New("could not run job"), jobType)
	}

	return nil
}

func (j *JobManager) startJob(ctx context.Context, job JobDefinition) {
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			j.log.DebugContext(ctx, "running job", "job", job.Type)
			if err := job.Run(ctx); err != nil {
				j.log.ErrorContext(ctx, "job failed",
					"job", job.Type,
					"error", err,
				)
			}
		case <-job.c:
			j.log.DebugContext(ctx, "running job", "job", job.Type)
			if err := job.Run(ctx); err != nil {
				j.log.ErrorContext(ctx, "job failed",
					"job", job.Type,
					"error", err,
				)
			}

		}
	}
}

func (j *JobManager) StartScheduler(ctx context.Context) {
	j.log.InfoContext(ctx, "starting scheduler")

	for _, job := range j.jobs {
		go j.startJob(ctx, job)
	}

	<-ctx.Done()
	j.log.InfoContext(ctx, "jobs finished")
}

func (j *JobManager) RegisterJob(ctx context.Context, job JobDefinition) {
	job.c = make(chan struct{}, 1)
	j.jobs = append(j.jobs, job)
}

func New(slogHandler slog.Handler) JobManager {
	return JobManager{
		log:  slog.New(slogHandler).With("service", "job-manager"),
		berr: bragerr.NewFactory("job-manager"),
	}
}

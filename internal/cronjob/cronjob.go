package cronjob

import "github.com/robfig/cron/v3"

type CronJob struct {
	c *cron.Cron
}

func NewCronJob() *CronJob {
	return &CronJob{
		c: &cron.Cron{},
	}
}

func (cj *CronJob) Add(s string, f func()) error {
	_, err := cj.c.AddFunc(s, f)
	if err != nil {
		return err
	}

	return nil
}

func (cj *CronJob) Start() {
	cj.c.Start()
}

func (cj *CronJob) Stop() {
	cj.c.Stop()
}

package main

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestLogProc(t *testing.T) {

	lines := strings.Split(
		`255 <158>1 2015-04-02T11:52:34.520012+00:00 host heroku router - at=info method=POST path="/users" host=myapp.com request_id=c1806361-2081-42e7-a8aa-92b6808eac8e fwd="24.76.242.18" dyno=web.1 connect=1ms service=37ms status=201 bytes=828
		229 <45>1 2015-04-02T11:48:16.839257+00:00 host heroku web.1 - source=web.1 dyno=heroku.35930502.b9de5fce-44b7-4287-99a7-504519070cba sample#load_avg_1m=0.01 sample#load_avg_5m=0.02 sample#load_avg_15m=0.03
		542 <134>1 2015-04-02T11:47:55+00:00 host app heroku-postgres - source=HEROKU_POSTGRESQL_TEAL addon=foo sample#current_transaction=6709 sample#db_size=18032824bytes sample#tables=16 sample#active-connections=4 sample#waiting-connections=0 sample#index-cache-hit-rate=0.99971 sample#table-cache-hit-rate=0.99892 sample#load-avg-1m=0.315 sample#load-avg-5m=0.22 sample#load-avg-15m=0.225 sample#read-iops=25.996 sample#write-iops=1.629 sample#memory-total=15666128kB sample#memory-free=233092kB sample#memory-cached=14836812kB sample#memory-postgres=170376kB
		542 <134>1 2015-04-02T11:47:55+00:00 host app heroku-redis - source=REDIS addon=foo sample#active-connections=73 sample#load-avg-1m=0 sample#load-avg-5m=0 sample#load-avg-15m=0 sample#read-iops=0 sample#write-iops=0 sample#memory-total=15664328kB sample#memory-free=14828336kB sample#memory-cached=237920kB sample#memory-redis=176289040bytes sample#hit-rate=0.85243 sample#evicted-keys=0
		222 <134>1 2017-05-13T15:35:33.787162+00:00 host app api - Scaled to mailer@3:Performance-L web@5:Standard-2X by user someuser@gmail.com
		222 <134>1 2015-04-07T16:01:43.517062+00:00 host heroku api - this_is="broken
		222 <134>1 2015-04-07T16:01:43.517062+00:00 host app api - Release v138 created by user foo@bar`, "\n")

	app := "test"
	tags := []string{"tag1", "tag2"}
	prefix := "prefix."
	s := loadServerCtx()
	s.in = make(chan *logData, 3)
	defer close(s.in)
	s.out = make(chan *logMetrics, 3)
	defer close(s.out)

	go logProcess(s.in, s.out)

	for i, l := range lines {
		log.WithField("line", l).Debug("Sending")
		s.in <- &logData{&app, &tags, &prefix, &lines[i]}
	}

	res := <-s.out
	if res.typ != routerMsg {
		t.Error("result must be ROUTE")
	}

	res = <-s.out
	if res.typ != dynoSampleMsg {
		t.Error("result must be DYNO SAMPLE")
	}

	res = <-s.out
	if res.typ != pgSampleMsg {
		t.Error("result must be POSTGRES SAMPLE")
	}

	res = <-s.out
	if res.typ != redisSampleMsg {
		t.Error("result must be REDIS SAMPLE")
	}

	res = <-s.out
	if res.typ != scalingMsg {
		t.Error("result must be SCALE")
	}

	res = <-s.out
	if res.typ != releaseMsg {
		t.Error("result must be RELEASE")
	}
}

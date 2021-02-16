package microservices

import (
	"context"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type LogsStorage struct {
	logs []*Event
	mu   sync.Mutex
}

func (storage *LogsStorage) save(e *Event) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	storage.logs = append(storage.logs, e)
}
func (storage *LogsStorage) getStatFromTime(timestamp int64, ignoreStatisticsMethod bool) *Stat {
	storage.mu.Lock()
	eventsFiltered := funk.Filter(storage.logs, func(e *Event) bool {
		return e.Timestamp >= timestamp
	}).([]*Event)
	storage.mu.Unlock()

	byMethod, byConsumer := make(map[string]uint64), make(map[string]uint64)
	for _, e := range eventsFiltered {
		if ignoreStatisticsMethod && e.Method == "/microservices.Admin/Statistics" {
			continue
		}
		if _, ok := byMethod[e.Method]; ok {
			storage.mu.Lock()
			byMethod[e.Method]++
			storage.mu.Unlock()
		} else {
			storage.mu.Lock()
			byMethod[e.Method] = 1
			storage.mu.Unlock()
		}
		if _, ok := byConsumer[e.Consumer]; ok {
			storage.mu.Lock()
			byConsumer[e.Consumer]++
			storage.mu.Unlock()
		} else {
			storage.mu.Lock()
			byConsumer[e.Consumer] = 1
			storage.mu.Unlock()
		}
	}

	return &Stat{
		Timestamp:  time.Now().UnixNano(),
		ByMethod:   byMethod,
		ByConsumer: byConsumer,
	}
}

type adminServiceImpl struct {
	aclRules          ACL
	mu                sync.Mutex
	logStreams        []Admin_LoggingServer
	statisticsStreams []Admin_StatisticsServer
	consumersWatchingLogs,
	consumersWatchingStatistics []string
	logsCh      chan *Event
	logsStorage *LogsStorage
}

func (adminService *adminServiceImpl) isAccessAllowed(ctx context.Context) bool {
	return checkAccessByContext(ctx, adminService.aclRules)
}

func (adminService *adminServiceImpl) SaveLogEvent(e *Event) {
	adminService.logsStorage.save(e)
	adminService.propagateEventStreams(e)
}

func (adminService *adminServiceImpl) propagateEventStreams(e *Event) {
	if len(adminService.logStreams) > 0 {
		adminService.logsCh <- e
	}
}

func (adminService *adminServiceImpl) Logging(data *Nothing, stream Admin_LoggingServer) error {

	_ = data
	if !adminService.isAccessAllowed(stream.Context()) {
		return status.Error(codes.Unauthenticated, "Access denied")
	}

	adminService.logStreams = append(adminService.logStreams, stream)

	consumer, _ := getConsumerFromContext(stream.Context())
	adminService.mu.Lock()
	if funk.IndexOfString(adminService.consumersWatchingLogs, consumer) == -1 {
		adminService.consumersWatchingLogs = append(adminService.consumersWatchingLogs, consumer)
		go adminService.SaveLogEvent(createLogEventByContext(stream.Context()))
	}
	adminService.mu.Unlock()

	for event := range adminService.logsCh {
		for _, stream := range adminService.logStreams {
			if err := stream.Send(event); err != nil {
				stream.Context().Done()
				return status.Error(codes.Internal, "Sending error")
			}
		}
	}

	stream.Context().Done()
	return nil
}

func (adminService *adminServiceImpl) Statistics(interval *StatInterval, stream Admin_StatisticsServer) error {

	if !adminService.isAccessAllowed(stream.Context()) {
		return status.Error(codes.Unauthenticated, "Access denied")
	}

	adminService.mu.Lock()
	adminService.statisticsStreams = append(adminService.statisticsStreams, stream)
	ignoreStatisticsMethod := len(adminService.statisticsStreams) > 1
	adminService.mu.Unlock()

	consumer, _ := getConsumerFromContext(stream.Context())
	adminService.mu.Lock()
	if funk.IndexOfString(adminService.consumersWatchingStatistics, consumer) == -1 {
		adminService.consumersWatchingStatistics = append(adminService.consumersWatchingStatistics, consumer)
		go adminService.SaveLogEvent(createLogEventByContext(stream.Context()))
	}
	adminService.mu.Unlock()

	lastLogsTimestamp := int64(0)
	go func() {
		for {
			<-time.After(time.Duration(interval.IntervalSeconds) * time.Second)
			stat := adminService.logsStorage.getStatFromTime(lastLogsTimestamp, ignoreStatisticsMethod)
			lastLogsTimestamp = time.Now().UnixNano()
			if err := stream.Send(stat); err != nil {
				return
			}
		}
	}()

	<-stream.Context().Done()
	return nil

}

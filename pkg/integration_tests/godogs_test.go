package main

import (
	cfg "calendar/pkg/config"
	logging "calendar/pkg/logger"
	"calendar/pkg/models"
	DB "calendar/pkg/psql"
	"calendar/pkg/rabbit"
	pb "calendar/pkg/services/api/gen"
	"calendar/pkg/services/scheduler"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	pt "github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"time"
)

type apiTest struct {
	config *cfg.APIConfig
	logger *zap.Logger
	DB *DB.PSQLInterface
	apiClient pb.APIClient
	Rabbit *rabbit.Rabbit
	SchedulerInstance *scheduler.Scheduler
	addEventState error
	updateEventState error
	getDayEventState *pb.Events
	getWeekEventState *pb.Events
	getMonthEventState *pb.Events
	baseTime time.Time
	uuids map[string]string
}

func NewApiTest() *apiTest{
	//host := "0.0.0.0"

	// Read api config
	apiConfig, err := cfg.ReadAPIConfig()
	if err != nil {
		log.Fatalf("reading api config error: %v \n", err)
	}
	// Read scheduler config
	schedulerConfig, err := cfg.ReadSchedulerConfig()
	if err != nil {
		log.Fatalf("reading scheduler config error: %v \n", err)
	}
	// Create api logger
	apiLogger, err := logging.CreateLogger(&apiConfig.API.Logger)
	if err != nil {
		log.Fatalf("creating api logger error: %v \n", err)
	}
	// Create scheduler logger
	schedulerLogger, err := logging.CreateLogger(&schedulerConfig.Scheduler.Logger)
	if err != nil {
		log.Fatalf("creating scheduler logger error: %v \n", err)
	}
	// Create GRPC Client
	conn, err := grpc.Dial(fmt.Sprintf("api:%v", apiConfig.API.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect to grpc server did not connect: %s", err)
	}

	grpcClient := pb.NewAPIClient(conn)

	// Connect to dbt
	dsn := fmt.Sprintf("user=%v password=%v dbname=%v host=%v sslmode=disable", apiConfig.API.DB.User, apiConfig.API.DB.Password, apiConfig.API.DB.DBName, apiConfig.API.DB.Host)
	//fmt.Printf("DSN:%v", dsn)
	dbInterface, err := DB.NewPSQLInterface(apiLogger, dsn)
	if err != nil {
		log.Fatalf("error during create PSQL instance: %v", err)
	}

	// Connect to Rabbit
	//apiConfig.API.Rabbit.Host = host
	apiConfig.API.Rabbit.QueueName = "calendar_test"
	rabbitInstance , err := rabbit.NewRabbit(apiLogger, &apiConfig.API.Rabbit)
	if err != nil {
		log.Fatalf("rabbit initialization give error: %v", err)
	}

	sch := &scheduler.Scheduler{
		Config: schedulerConfig,
		Rabbit: rabbitInstance,
		PSQL: dbInterface,
		Logger: schedulerLogger,
	}


	return &apiTest{
		config: apiConfig,
		logger: apiLogger,
		DB: dbInterface,
		apiClient: grpcClient,
		SchedulerInstance: sch,
		Rabbit: rabbitInstance,
		addEventState: nil,
		updateEventState: nil,
		getDayEventState: nil,
		getWeekEventState: nil,
		getMonthEventState: nil,
		baseTime: time.Date(2019, 10, 5, 9, 30, 0, 0, time.UTC), // without timezone
		uuids: map[string]string{
			"Create": "3ebcf64b-1402-47da-ab23-3656b3268135",
			"Update": "4ebcf64b-1402-47da-ab23-3656b3268135",
			"Day": "5ebcf64b-1402-47da-ab23-3656b3268135",
			"Week": "6ebcf64b-1402-47da-ab23-3656b3268135",
			"Month": "7ebcf64b-1402-47da-ab23-3656b3268135",
			"Old": "8ebcf64b-1402-47da-ab23-3656b3268135",
			"Notify": "9ebcf64b-1402-47da-ab23-3656b3268135",
		},
	}
}


func (as *apiTest) iAddEventThroughGprcRequest() error {

	date, _ := pt.TimestampProto(as.baseTime)
	newUUID := as.uuids["Create"]
	event := &pb.Event{
		User: "Vasya",
		UUID: newUUID,
		//UUID: newUUID.String(),
		Summary: "Create event summary",
		StartDate: date,
		EndDate: date,
		Description: "Create event description",
		NotifyTime: date,
	}

	ctx := context.Background()

	fmt.Printf("\nCreate event...\n")
	_, err := as.apiClient.EventCreate(ctx, event)
	as.addEventState = err

	return nil
}

func (as *apiTest) theCreateResponseShouldNotContainsError() error {
	if as.addEventState != nil {
		return fmt.Errorf("Event create response return error: %v", as.addEventState)
	}
	return nil
}


func (as *apiTest) iUpdateEventThroughGprcRequest() error {
	newUUID := as.uuids["Update"]
	date, _ := pt.TimestampProto(as.baseTime)

	event := &pb.Event{
		User: "Vasya",
		UUID: newUUID,
		//UUID: newUUID.String(),
		Summary: "Updated event summary",
		StartDate: date,
		EndDate: date,
		Description: "Updated event description",
		NotifyTime: date,
	}

	ctx := context.Background()

	fmt.Printf("\nUpdate event...\n")
	_, err := as.apiClient.EventUpdate(ctx, &pb.UpdateRequest{ID: as.uuids["Create"], Event:event})
	as.updateEventState = err

	return nil
}

func (as *apiTest) theUpdateResponseShouldNotContainsError() error {
	if as.updateEventState != nil {
		return fmt.Errorf("Event update response return error: %v", as.updateEventState)
	}
	return nil
}


func (as *apiTest) iGetDayEventsThroughGprcRequest() error {

	date, _ := pt.TimestampProto(as.baseTime)
	newUUID := as.uuids["Day"]
	event := &pb.Event{
		User: "Vasya",
		UUID: newUUID,
		//UUID: newUUID.String(),
		Summary: "Day event summary",
		StartDate: date,
		EndDate: date,
		Description: "Day event description",
		NotifyTime: date,
	}
	ctx := context.Background()
	_, err := as.apiClient.EventCreate(ctx, event)
	if err != nil {
		return fmt.Errorf("iGetDayEventsThroughGprcRequest EventCreate give error: %v", err)
	}

	events, err := as.apiClient.EventDayList(ctx, date)
	as.getDayEventState = events
	if err != nil {
		return fmt.Errorf("EventDayList give error: %v", err)
	}

	return nil

}

func (as *apiTest) theResponseShouldContainsDayEvents() error {

	for _, ev := range as.getDayEventState.Events {
		if ev.UUID == as.uuids["Day"]{
			return nil
		}
	}

	return fmt.Errorf("created day event is not contains in day events list")
}


func (as *apiTest) iGetWeekEventsThroughGprcRequest() error {

	date, _ := pt.TimestampProto( time.Date(as.baseTime.Year(), as.baseTime.Month(), as.baseTime.Day() + 2, as.baseTime.Hour(), as.baseTime.Minute(), as.baseTime.Second(), as.baseTime.Nanosecond(), time.UTC))

	newUUID := as.uuids["Week"]
	event := &pb.Event{
		User: "Vasya",
		UUID: newUUID,
		//UUID: newUUID.String(),
		Summary: "Week event summary",
		StartDate: date,
		EndDate: date,
		Description: "Week event description",
		NotifyTime: date,
	}
	ctx := context.Background()
	_, err := as.apiClient.EventCreate(ctx, event)
	if err != nil {
		return fmt.Errorf("iGetWeekEventsThroughGprcRequest EventCreate give error: %v", err)
	}

	events, err := as.apiClient.EventWeekList(ctx, date)
	as.getWeekEventState = events
	if err != nil {
		return fmt.Errorf("EventWeekList give error: %v", err)
	}

	return nil
}

func (as *apiTest) theResponseShouldContainsWeekEvents() error {

	for _, ev := range as.getWeekEventState.Events {
		if ev.UUID == as.uuids["Week"]{
			return nil
		}
	}
	return fmt.Errorf("created week event is not contains in week events list")
}


func (as *apiTest) iGetMonthThroughGprcRequest() error {

	date, _ := pt.TimestampProto( time.Date(as.baseTime.Year(), as.baseTime.Month(), as.baseTime.Day() + 20, as.baseTime.Hour(), as.baseTime.Minute(), as.baseTime.Second(), as.baseTime.Nanosecond(), time.UTC))

	newUUID := as.uuids["Month"]
	event := &pb.Event{
		User: "Vasya",
		UUID: newUUID,
		Summary: "Week event summary",
		StartDate: date,
		EndDate: date,
		Description: "Week event description",
		NotifyTime: date,
	}
	ctx := context.Background()
	_, err := as.apiClient.EventCreate(ctx, event)
	if err != nil {
		return fmt.Errorf("iGetMonthThroughGprcRequest EventCreate give error: %v", err)
	}

	events, err := as.apiClient.EventMonthList(ctx, date)
	as.getMonthEventState = events
	if err != nil {
		return fmt.Errorf("EventMonthList give error: %v", err)
	}

	return nil
}

func (as *apiTest) theResponseShouldContainsMonthEvents() error {

	for _, ev := range as.getMonthEventState.Events {
		if ev.UUID == as.uuids["Month"]{
			return nil
		}
	}
	return fmt.Errorf("created month event is not contains in month events list")
}


func (as *apiTest) schedulerClearOldEvents() error {

	date, _ := pt.TimestampProto( time.Date(as.baseTime.Year() - 2, as.baseTime.Month(), as.baseTime.Day(), as.baseTime.Hour(), as.baseTime.Minute(), as.baseTime.Second(), as.baseTime.Nanosecond(), time.UTC))

	newUUID := as.uuids["Old"]
	event := &pb.Event{
		User: "Vasya",
		UUID: newUUID,
		Summary: "Old event summary",
		StartDate: date,
		EndDate: date,
		Description: "Old event description",
		NotifyTime: date,
	}
	ctx := context.Background()
	_, err := as.apiClient.EventCreate(ctx, event)
	if err != nil {
		return fmt.Errorf("schedulerClearOldEvents EventCreate give error: %v", err)
	}
	err = as.SchedulerInstance.CleanEventsOnce()
	if err != nil {
		return fmt.Errorf("CleanEventsOnce give error: %v", err)
	}

	return nil
}

func (as *apiTest) oldEventsWillBeRemovedFromDB() error {

	findedEvent, err := as.DB.GetEventByUUID(as.uuids["Old"])
	if err != nil{
		return fmt.Errorf("oldEventsWillBeRemovedFromDB GetEventByUUID give error: %v", err)
	}
	if findedEvent != nil{
		return fmt.Errorf("GetEventByUUID return old event which should be already deleted!")
	}

	return nil
}


func (as *apiTest) eventNotifyTimeIsComing() error {

	startDate, _ := pt.TimestampProto(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour() + 1, time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC))
	endDate, _ := pt.TimestampProto(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour() + 3, time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC))
	notifyDate, _ := pt.TimestampProto(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour() - 1, time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC))

	newUUID := as.uuids["Notify"]
	event := &pb.Event{
		User: "Vasya",
		UUID: newUUID,
		Summary: "Notify event summary",
		StartDate: startDate,
		EndDate: endDate,
		Description: "Notify event description",
		NotifyTime: notifyDate,
	}
	ctx := context.Background()
	_, err := as.apiClient.EventCreate(ctx, event)
	if err != nil {
		return fmt.Errorf("eventNotifyTimeIsComing EventCreate give error: %v", err)
	}
	return nil
}

func (as *apiTest) schedulerSendNotificationToNotificatorQueue() error {

	// Consume from queue
	response, err := as.Rabbit.Channel.Consume(
		as.Rabbit.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		as.Rabbit.Logger.Sugar().Infof("Consuming error: %v", err)
		return err
	}

	// Scheduler send event notification to rabbit
	_ = as.SchedulerInstance.SendEventsOnce()

	//as.Rabbit.Channel.Close()
	for resp := range response {
		var event models.DBEvent
		err = json.Unmarshal(resp.Body, &event)
		if err != nil {
			return fmt.Errorf("Unmarshaling event give error: %v", err)
		} else {
			var evUUID string
			err = event.UUID.AssignTo(&evUUID)
			if err != nil {
				return fmt.Errorf("assigning pg.UUID to string give error: %v", err)
			} else {
				if evUUID == as.uuids["Notify"] {
					return nil
				}
			}

		}
	}
	return fmt.Errorf("Notify event notification is not exists in rabbit channel")
}


func (as *apiTest) ClearAll(*gherkin.Feature) {

	as.logger.Info("Clearing DB after all tests...")

	allEvents, err := as.DB.GetAllEvents()
	if err != nil {
		as.logger.Sugar().Errorf("ClearAll | GetAllEvents give error: %v", err)
		return
	}
	deleted := 0
	for _, ev := range allEvents {
		var pbUUID string
		err := ev.UUID.AssignTo(&pbUUID)
		if err != nil {
			as.logger.Sugar().Errorf("ClearAll | Cannot scan UUID from event: %v", err)
			return
		}

		err = as.DB.DeleteEvent(pbUUID)
		if err != nil {
			as.logger.Sugar().Errorf("ClearAll | DeleteEvent give error: %v", err)
			return
		}
		deleted += 1
	}

	as.logger.Sugar().Infof("ClearAll | Deleted %v events", deleted)
	return
}

func FeatureContext(s *godog.Suite) {

	newTestApiInstance := NewApiTest()

	// Clear DB after tests
	s.AfterFeature(newTestApiInstance.ClearAll)

	s.Step(`^I add event through gprc request$`, newTestApiInstance.iAddEventThroughGprcRequest)
	s.Step(`^The create response should not contains error$`, newTestApiInstance.theCreateResponseShouldNotContainsError)
	s.Step(`^I update event through gprc request$`, newTestApiInstance.iUpdateEventThroughGprcRequest)
	s.Step(`^The update response should not contains error$`, newTestApiInstance.theUpdateResponseShouldNotContainsError)
	s.Step(`^I get day events through gprc request$`, newTestApiInstance.iGetDayEventsThroughGprcRequest)
	s.Step(`^The response should contains day events$`, newTestApiInstance.theResponseShouldContainsDayEvents)
	s.Step(`^I get week events through gprc request$`, newTestApiInstance.iGetWeekEventsThroughGprcRequest)
	s.Step(`^The response should contains week events$`, newTestApiInstance.theResponseShouldContainsWeekEvents)
	s.Step(`^I get month events through gprc request$`, newTestApiInstance.iGetMonthThroughGprcRequest)
	s.Step(`^The response should contains month events$`, newTestApiInstance.theResponseShouldContainsMonthEvents)
	s.Step(`^Scheduler clear old events$`, newTestApiInstance.schedulerClearOldEvents)
	s.Step(`^Old events will be removed from DB$`, newTestApiInstance.oldEventsWillBeRemovedFromDB)
	s.Step(`^Event notify time is coming$`, newTestApiInstance.eventNotifyTimeIsComing)
	s.Step(`^Scheduler send notification to notificator queue$`, newTestApiInstance.schedulerSendNotificationToNotificatorQueue)

}

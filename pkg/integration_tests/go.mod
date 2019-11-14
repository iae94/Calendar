module calendar/pkg/integration_tests

go 1.12

require (
	calendar/pkg/config v0.0.0
	calendar/pkg/logger v0.0.0
	calendar/pkg/models v0.0.0
	calendar/pkg/psql v0.0.0
	calendar/pkg/rabbit v0.0.0
	calendar/pkg/services/api/gen v0.0.0
	calendar/pkg/services/scheduler v0.0.0

	github.com/DATA-DOG/godog v0.7.13
	github.com/golang/protobuf v1.3.2
)

replace calendar/pkg/config v0.0.0 => ../../pkg/config

replace calendar/pkg/logger v0.0.0 => ../../pkg/logger

replace calendar/pkg/services/api/gen v0.0.0 => ../../pkg/services/api/gen

replace calendar/pkg/services/api/server v0.0.0 => ../../pkg/services/api/server

replace calendar/pkg/models v0.0.0 => ../../pkg/models

replace calendar/pkg/psql v0.0.0 => ../../pkg/psql

replace calendar/pkg/rabbit v0.0.0 => ../../pkg/rabbit

replace calendar/pkg/services/scheduler v0.0.0 => ../../pkg/services/scheduler

Feature: Api server
  As client of calendar service
  In order to manage my calendar
  I want to add/update/get events

  Scenario: Add event
    When I add event through gprc request
    Then The create response should not contains error

  Scenario: Update event
    When I update event through gprc request
    Then The update response should not contains error

  Scenario: Get day events
    When I get day events through gprc request
    Then The response should contains day events

  Scenario: Get week events
    When I get week events through gprc request
    Then The response should contains week events

  Scenario: Get month events
    When I get month events through gprc request
    Then The response should contains month events

  Scenario: Clean old events
    When Scheduler clear old events
    Then Old events will be removed from DB

  Scenario: Send notifications
    When Event notify time is coming
    Then Scheduler send notification to notificator queue

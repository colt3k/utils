# Usage

- Define a `Rule` and call `Process`

`retry` runs a task with backoff-aware retry behavior and context cancellation.

NOTE: You can vary how your rule is configured, if you don't define fields they will be
    defaulted if you supply the MaxInterval you MUST set MaxIntervalDurationType

- Example

        retry.Rule{
            MaxAttempts:       3,
            MaxInterval:       10,
            MaxIntervalDurationType: time.Second,
            MaxElapsed:        15 * time.Second,
            SleepDurationType: retry.Seconds,
        }

## Fields Explained

- `MaxAttempts`: how many times to retry
- `MaxInterval`: maximum time between retries entered as a whole number
- `MaxIntervalDurationType`: this is the time unit intended for `MaxInterval` using `time.Duration`
- `MaxElapsed`: maximum amount of time for all retries before quitting
- `SleepDurationType`: the sleep format to wait in, for example milliseconds, seconds, or minutes

## Other Examples

- Minutes

          retry.Rule{
              MaxAttempts:             3,
              MaxInterval:             1,
              MaxIntervalDurationType: time.Minute,
              MaxElapsed:              10 * time.Minute,
              SleepDurationType:       retry.Minutes,
          }

- Minutes with delay in seconds

          retry.Rule{
              MaxAttempts:             3,
              MaxInterval:             1,
              MaxIntervalDurationType: time.Minute,
              MaxElapsed:              10 * time.Minute,
              SleepDurationType:       retry.Seconds,
          }

- Milliseconds

          retry.Rule{
              MaxAttempts:             3,
              MaxInterval:             1000,
              MaxIntervalDurationType: time.Millisecond,
              MaxElapsed:              15000 * time.Millisecond,
              SleepDurationType:       retry.Milliseconds,
          }

## Development

- `cd retry && go test ./...`
- `Process` uses `context.Context`, so cancellation behavior should be covered when changing retry logic

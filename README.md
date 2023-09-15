## LimitedExec: A Rate Limited Executor for Go

### Overview

`LimitedExec` is a Go package designed to manage the execution rate of certain functions based on unique keys. This ensures that the associated function for a particular key isn't executed more than once within a specified time duration.

A common use case might be in situations where there's a frequent, repetitive event (e.g., requests from a banned IP address) that could overwhelm logs or system resources. Using `LimitedExec`, one can throttle these events to ensure they're logged or processed at a manageable rate.

### Features

- **Rate Limiting**: Provides an easy-to-use function to rate limit executions.
- **Concurrent Safety**: Uses Go's `sync.Map` and `sync.Mutex` to ensure concurrent safety.
- **Periodic Cleanup**: Periodically cleans up the map to ensure it doesn’t grow infinitely.

### Usage

The main function, `Exec`, is used to attempt the execution of a function based on a key:

```go
Exec(k interface{}, t time.Duration, fun func())
```

**Parameters**:
- `k`: Unique key for identifying the function (e.g., "banned "+ip).
- `t`: Time duration within which the function shouldn’t be executed again.
- `fun`: Function to be executed if allowed.

**Example**:

Consider you're tracking events of banned IPs:

```go
12:12: came banned ip xx
12:13: came banned ip xx
...
```

Using `LimitedExec`, you can reduce the frequency:

```go
limitedexec.Exec("banned "+ip, time.Minute, func() { log.Print("came banned ip xx") })
```

This will produce:

```go
12:12: came banned ip xx
13:12: came banned ip xx
...
```

### Testing

The package also contains tests to validate its functionality. It tests scenarios like multiple executions within the limited duration and checks if the map cleanup mechanism works as expected.
# Crontab Resource

Implements a resource that reports new versions when the current time
matches the crontab expression

## Add to your Concourse deployment

Note: Specifying the resource_types in your manifest will overwrite the [defaults](https://github.com/concourse/concourse/blob/master/jobs/groundcrew/spec#L68-L96). You will need to add the defaults in your manifest as well.

```
groundcrew:
  resource_types:
  - type: krafa
    image: docker:///cftoolsmiths/croncourse-resource
```

## Source Configuration

* `expression`: *Required.* The crontab expression:

    |field       | allowed values |
    |-------------|----------------|
    |minute       | 0-59 |
    |hour         | 0-23 |
    |day of month | 1-31 |
    |month        | 1-12 (or names, see below) |
    |day of week  | 0-7 (0 or 7 is Sun, or use names) |

  e.g.

    `0 23 * * 1-5` # Run at 11:00pm from Monday to Friday

* `location`: *Optional.* Defaults to UTC. Accepts any timezone that
  can be parsed by https://godoc.org/time#LoadLocation

  e.g.

  `America/New_York`
       
  `America/Vancouver`

## Behavior

### `check`: Report the current time.

Returns `time.Now()` as the version only if a minute since we last
fired matches the crontab expression. The first time the script runs
it will fire if a minute in the last hour matches the crontab
expression.

#### Parameters

*None.*

### `in`: Report the given time

If triggered by `check`, returns the original version as the resulting
version.

#### Parameters

1. *Output directory.* The directory where the in script will store
   the requested version

### `out`: Not implemented.

# Crontab Resource

Implements a resource that reports new versions when the current time
matches the crontab expression

---
## Update your pipeline

Update your pipeline to include this new declaration of resource types. See the example pipeline yml snippet below or the Concourse docs for more details [here](https://concourse.ci/configuring-resource-types.html).
```
---
resource_types:
- name: cron-resource
  type: docker-image
  source:
    repository: pivotal-cf-experimental/cron-resource

resources:
  - name: 10-min-trigger
    type: cron-resource
    source:
      expression: "*/10 * * * *"
      location: "America/New_York"
      fire_immediately: true
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

* `fire_immediately`: *Optional.* Defaults to false. Immediately triggers the resource the first time it is checked.

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

---
## Developer Notes

### Building the jobs

To build the resource's go binaries, run the following command from within the cron-resource directory:

```
docker run -v "$PWD":/go/src/github.com/pivotal-cf-experimental/cron-resource/ \
           -it golang:1.7 \
           /bin/bash /go/src/github.com/pivotal-cf-experimental/cron-resource/build_in_docker_container.sh
```

You should see two new binaries named `built-in` and `built-check`.

### Testing the go binaries

Start an interactive session of your docker container to run the binaries:

```
docker run -v "$PWD":/go/src/github.com/pivotal-cf-experimental/cron-resource/ \
           -it golang:1.7 \
           /bin/sh
```

Within the interactive session:

```
cd /go/src/github.com/pivotal-cf-experimental/cron-resource/
./built-check
```

It looks like it hangs, but it's waiting for you to enter some JSON:

```
{"source":{"expression":"* * * * *","location":"America/New_York"} } # Paste this in after running ./built-check
[{"time":"2016-08-19T10:15:27.183011117-04:00"}] # This is the successful output
>>>>>>> Updates README with instructions on how to use Concourse's resource types, Adds developer notes.
```

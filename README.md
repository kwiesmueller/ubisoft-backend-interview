# Ubisoft backend developer role test.
[![Go Report Card](https://goreportcard.com/badge/github.com/kwiesmueller/ubisoft-backend-interview)](https://goreportcard.com/report/github.com/kwiesmueller/ubisoft-backend-interview)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/dd692d4293644b1ca429a1f69001c225)](https://www.codacy.com/app/kwiesmueller/ubisoft-backend-interview?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=kwiesmueller/ubisoft-backend-interview&amp;utm_campaign=Badge_Grade)
[![API Docs](https://img.shields.io/badge/API%20Docs-available-green.svg)](https://ubisoftbackendinterview.docs.apiary.io/#)
[![Build Status](https://travis-ci.org/kwiesmueller/ubisoft-backend-interview.svg?branch=master)](https://travis-ci.org/kwiesmueller/ubisoft-backend-interview)
[![Docker Repository on Quay](https://quay.io/repository/finch/ubisoft-backend-interview/status "Docker Repository on Quay")](https://quay.io/repository/finch/ubisoft-backend-interview)

## About

- **Author:** kwiesmueller
- **Started:** March 17th, 2018 - 8pm GMT+1
- **Purpose:** Just for fun... Get in touch if you like. 
- **License:** Do whatever you want with it...
Would be happy to know if this appears in any games I play, though.

## Description

### Intro

If you wonder, why that much text for a test I wasn't even asked for? Let's call out loud thinking and self reminding. While this text might explain my thoughts and decisions to somebody at Ubi stumbling over my PR, it is mainly thought as self-improvement. I took the chance to build something easy aside from my other projects to go over my own decisions, patterns and best/worst practices. I guess when I am done I might convert this into some kind of endless refactoring experiment of sorts. Who knows.

### Facts

My implementation of this test is written in Golang. This choice was somehow obvious to me, as I 
on one hand really like to work with go, and on the other hand think it is a perfect match on this.
The implementation might be slightly longer than using NodeJS, Python or others but for that, one gets a solid, well tested and naturally multi-threading service with a very small footprint.

Due to the motivation to ensure production grade quality and the testing added to all parts of the application, the entire development took longer than expected, which you might see based on my commits. Other languages, especially those allowing much easier mocking for tests might have been faster there, but I am still happy with the amount of time spent and the resulting quality.
As one point in the original task was to only comment code if needed, only comments specifically required by Golang were added.

While using a quite general persistence design internally, my implementation of it is using PostgreSQL.
The choice for SQL is obvious I think, as the data is structured and it's design is known.
PostgreSQL should be able to easily handle this service both in single or as replicated deployments (maybe using something like spilo for clustering PostgreSQL might be good in a production environment).

I also thought about not doing any common database persistence as the required limit by 15 looked like feedback further in the past wouldn't ever be used, but (as you can see in my service implementation) I wanted to keep the option open for retrieving more data.

Switching to a different persistence technology (like using batcher internally) or even existing NoSQL databases, is fairly easy. All that has to be done is implementing the [Repository Interface](pkg/feedback/repository.go). I might even supply another database option once I'm finished.

The app itself is naturally packet into a docker image, but can also be built and deployed as single binary if necessary. Kubernetes manifests are supplied with the image as an example, too.

API documentation is available on [Apiary](https://ubisoftbackendinterview.docs.apiary.io/#).
The database design can be found in [db.sql](db.sql). 

## Dependencies

- First of all [Golang](https://golang.org/dl/) 1.9 is required (older and newer might work, but is not tested by me for now)
- For building the image and starting PostgreSQL, [Docker](https://www.docker.com/community-edition) is required

All code dependencies and vendor libraries are checked in, so there should be no need to do anything else. If you want to test the code run `make deps` before to get the used Golang testing and coverage tools.

Dependency management is done with [dep](https://github.com/golang/dep).

## Usage

To run the implementation the database server has to be up first. To do so, the [Makefile](Makefile) contains a little helper to start the [postgres container](helpers/make_db).
To start it, execute the following:
```bash
DB_PASSWORD=db make start-db
```
Note that this starts the database attached to your tty, so best do this in a separate terminal.
If you are starting the db for the first time, building the [database structure](db.sql) is required:
```bash
cat db.sql | docker exec -i ubisoft-backend-interview-db psql -U db -d db -
```
After that, the service should be able to reach your local instance of PostgreSQL and work.

To start, either run `make dev` for debug output or `make run` to build and run the binary.

## Logging and Monitoring

Please note, that due to the used logging library configuration (down at the core [uber-go/zap](go.uber.org/zap)) running without debug won't print INFO either. This could be changed easily, but in my own deployments I saw this information is mostly not required and very verbose. If there is the need of debugging through info logs, I prefer real debugging (or cloud debugging using breakpoints etc.).

For monitoring my choice is prometheus, for which a simple middleware is being used to collect usage metrics and statistics.

Some benefits of the chosen logging library are type-safe structured logging, good performance even on high throughput and the added [Sentry](https://sentry.io) integration which will forward all errors logged to the supplied Sentry instance (see the `-sentryDsn` parameter).
Additionally it would be no effort to use the built-in tracing library [Jaeger](github.com/uber/jaeger-lib), which has been spared for now, as it might have been out of scope.

Error handling is a question I thought for some time as well. As you might see on my code, I tend to pass on errors until the end. Database grade errors get caught and overwritten with more user-friendly ones as I did not want to disclose database insight to the user, but cases like the query param verification are still something i have to decide on.
At the moment errors from strconv get (while wrapped) returned to the user (see [API Docs](https://ubisoftbackendinterview.docs.apiary.io/#)). As I am not entirely sure about the kind of deployment or what access the user will have to those error messages, this might be not most beautiful decision.

It would be possible (and pretty easy) to 


## Test

You are to write a new micro-service that will allow users to share feedback on their last game session and allow visibility to a live operations team.

Users can rate their session from 1 to 5 and leave a comment. Session id is provided in the url path and the user id is in the header named `Ubi-UserId`.

**Players within the same gaming session rate via the same session id, but a player can only leave one feedback per session.**

Following RESTful principles:

1 - Write an HTTP endpoint for players to post a new feedback for a session.

2 - Write an HTTP endpoint to get the **last** 15 feedbacks left by players and allow filtering by rating.


## Rules

This test has a time limit of one week. To submit your result use a github or gitlab repository.
You can share answers/design/documentation via markdown in the repository.
Submit all your work at the end of your week, whether it's completed or not.
No specific language is required. You may use the language that you're most comfortable with, or even explore a new one. The same applies for the database you choose for your tool.


Tips to improve your application :
- Document your database design. (schema/index/query)
- Document your api design. (routes, payload)
- Document your code only if needed.
- Document how to run your project or run tests. (easier for us to evaluate)


Good Luck :v:
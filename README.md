# Ubisoft backend developer role test.
[![Go Report Card](https://goreportcard.com/badge/github.com/kwiesmueller/ubisoft-backend-interview)](https://goreportcard.com/report/github.com/kwiesmueller/ubisoft-backend-interview)
[![Build Status](https://travis-ci.org/kwiesmueller/ubisoft-backend-interview.svg?branch=master)](https://travis-ci.org/kwiesmueller/ubisoft-backend-interview)
[![Docker Repository on Quay](https://quay.io/repository/finch/ubisoft-backend-interview/status "Docker Repository on Quay")](https://quay.io/repository/finch/ubisoft-backend-interview)

## About
- **Author:** kwiesmueller
- **Started:** March 17th, 2018 - 8pm GMT+1
- **Purpose:** Just for fun... Get in touch if you like. 
- **License:** Do whatever you want with it...
Would be happy to know if this appears in any games I play, though.

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
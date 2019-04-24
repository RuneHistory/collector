# RuneHistory Collector Microservice

This microservice:
* Assignes buckets to accounts on creation (accounts microservice should emit an event)
* Fetches buckets of users every 5 minutes, get's the latest highscore for each, and emits an event for each record

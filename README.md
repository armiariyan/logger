# Logger
Logger is an interface to produce common log format that already specified in document

## Usage

Open folder `/example` to see the sample usage.

## Key Features
* Logging context using middleware (see example): now you can search all related log with same Thread ID or Journey ID without specifying it in each call.
* Add additional data that available across all log, such as user_id via context [1]
* Multi-writer using Zap: neat code
* Closing io.Writer on rotate logs and Kafka

[1] Please note don't add large data as you will need more memory to pass data via context
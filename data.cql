-- Keyspace
CREATE KEYSPACE IF NOT EXISTS project WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Back-end
CREATE TABLE IF NOT EXISTS project.spans (
span_id         text,
trace_id        text,
session_id      text,
user_id         text,
chapter_id      text,

time_sent       bigint,
trigger_route   text,
status_code     smallint,
data            text,
request_data    text,
time_duration   duration,

PRIMARY KEY(span_id, time_sent)
) WITH CLUSTERING ORDER BY (time_sent DESC);

CREATE INDEX ON project.spans (trace_id);
CREATE INDEX ON project.spans (session_id);
CREATE INDEX ON project.spans (user_id);
CREATE INDEX ON project.spans (chapter_id);
CREATE INDEX ON project.spans (time_sent);
CREATE INDEX ON project.spans (request_data);

-- Front-end
CREATE TABLE IF NOT EXISTS project.events (
session_id      text,
user_id         text,
chapter_id      text,
data            text,

PRIMARY KEY(data)
);
CREATE INDEX ON project.events (session_id);
CREATE INDEX ON project.events (user_id);
CREATE INDEX ON project.events (chapter_id);
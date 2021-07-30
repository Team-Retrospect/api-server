from cassandra.cluster import Cluster
import json
import pytest
import requests
import yaml

with open('./config.yml', 'r') as f:
    CONFIG = yaml.load(f, Loader=yaml.FullLoader)
    HOST = f"http{'s' if CONFIG['use_https'] else ''}://localhost"
    PORT = CONFIG['port']
with open('./tests/data/sample_span.json', 'r') as f :
    SAMPLE_SPAN = json.load(f)
with open('./tests/data/sample_trace.json', 'r') as f :
    SAMPLE_TRACE = json.load(f)
with open('./tests/data/sample_snapshot.json', 'r') as f :
    SAMPLE_SNAPSHOT = json.load(f)

IDS = {
    'span'      : 'some_span_uuid',
    'trace'     : 'some_trace_uuid',
    'session'   : 'some_session_uuid',
    'user'      : 'some_user_uuid',
    'chapter'   : 'some_chapter_uuid',
    'trigger'   : 'get http://some_trigger_route.com/dashboard',
    'nonexist'  : 'this_should_not_exist',
    'data'      : 'InsiU2FtcGxlIjoiVGhpcyBpcyB0ZXN0In0i',
    'time_sent' : 1,
}
HEADERS = {
    "user-id"       : "some_user_uuid",
    "session-id"    : "some_session_uuid",
    "chapter-id"    : "some_chapter_uuid",
}


CLUSTER = Cluster([CONFIG['cluster']], port=9042)
KEYSPACE = 'retrospect'
SESSION = CLUSTER.connect(KEYSPACE, wait_for_all_pools=True)
SESSION.execute(f'USE {KEYSPACE};')

# setup:

INSERT_QUERY = f"INSERT INTO {'{table}'} JSON '{'{payload}'}';"

def setup_insert_sample_span():
    SESSION.execute(INSERT_QUERY.format(
        table='spans',
        payload=json.dumps({
            'span_id': IDS['span'],
            'trace_id': IDS['trace'],
            'session_id': IDS['session'],
            'user_id': IDS['user'],
            'chapter_id': IDS['chapter'],
            'trigger_route': IDS['trigger'],
            'time_sent': 1,
            })
    ))
def setup_insert_sample_trace():
    SESSION.execute(INSERT_QUERY.format(
        table='events',
        payload=json.dumps({
            'session_id': IDS['session'],
            'user_id': IDS['user'],
            'chapter_id': IDS['chapter'],
            'data': IDS['data'],
            })
    ))
def setup_insert_sample_snapshot():
    SESSION.execute(INSERT_QUERY.format(
        table='snapshots',
        payload=json.dumps({
            'session_id': IDS['session'],
            'data': IDS['data'],
            })
    ))

# teardown:

def teardown_delete_sample_span():
    SESSION.execute(f"DELETE FROM spans WHERE span_id='{IDS['span']}';")

def teardown_delete_sample_trace():
    SESSION.execute(f"DELETE FROM events WHERE data='{IDS['data']}';")

def teardown_delete_sample_snapshot():
    SESSION.execute(f"DELETE FROM snapshots WHERE data='{IDS['data']}';")


# Assertions:

def assert_is_json_array(r):
    assert (r.text[0] == '[') and (r.text[-1] == ']'), "body isn't an array"

def assert_is_json_object(r):
    assert (r.text[0] == '{') and (r.text[-1] == '}'), "body isn't an object"

def assert_is_json(r):
    try :
        j = r.json()
        assert True
    except json.decoder.JSONDecodeError :
        pytest.fail("doesn't parse by json")

def assert_is_empty_array(r):
    assert r.text == '[]', "returned nonempty array"

def assert_contents_are_type(r, t):
    try :
        j = r.json()
        assert type(j[0]) == t, f"type is not {type(t)}"
    except json.decoder.JSONDecodeError :
        pytest.fail("doesn't parse by json")

def assert_isnt_empty(r):
    assert len(r.json())

def assert_is_ok(r):
    assert r.ok, "returned error code"

def assert_is_not_ok(r):
    assert not r.ok, "returned success but shouldnt"

# others:

def url(endpoint):
    return HOST + PORT + endpoint

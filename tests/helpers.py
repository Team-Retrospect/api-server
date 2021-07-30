import pytest
import requests
import json

HOST = "http://localhost:8081"
with open('./data/sample_span.json', 'r') as f :
    SAMPLE_SPAN = json.load(f)
with open('./data/sample_trace.json', 'r') as f :
    SAMPLE_TRACE = json.load(f)
with open('./data/sample_snapshot.json', 'r') as f :
    SAMPLE_SNAPSHOT = json.load(f)

IDS = {
    'user'      : 'some_user_uuid',
    'span'      : 'some_span_uuid',
    'trace'     : 'some_trace_uuid',
    'session'   : 'some_session_uuid',
    'chapter'   : 'some_chapter_uuid',
    'trigger'   : 'get http://some_trigger_route.com/dashboard',
    'nonexist'  : 'this_should_not_exist',
}
HEADERS = {
    "user-id"       : "some_user_uuid",
    "session-id"    : "some_session_uuid",
    "chapter-id"    : "some_chapter_uuid",
}

def url(endpoint):
    return HOST + endpoint

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
        assert type(j[0]) == t
    except json.decoder.JSONDecodeError :
        pytest.fail("doesn't parse by json")


def assert_isnt_empty(r):
    assert len(r.json())

def assert_is_ok(r):
    assert r.ok, "returned error code"

def assert_is_not_ok(r):
    assert not r.ok, "returned success but shouldnt"

# setup:

def setup_insert_sample_span():
    # (todo: use a cassandra driver to input this?)
    requests.post(url('/spans'), data=json.dumps(SAMPLE_SPAN), headers=HEADERS)
def setup_insert_sample_trace():
    # (todo: use a cassandra driver to input this?)
    requests.post(url('/events'), data=json.dumps(SAMPLE_TRACE), headers=HEADERS)
def setup_insert_sample_snapshot():
    # (todo: use a cassandra driver to input this?)
    requests.post(url('/events'), data=json.dumps(SAMPLE_SNAPSHOT), headers=HEADERS)

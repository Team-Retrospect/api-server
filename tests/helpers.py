import pytest
import json

HOST = "http://localhost:8081"
with open('sample_span.json', 'r') as f :
    SAMPLE_SPAN = json.load(f)
with open('sample_trace.json', 'r') as f :
    SAMPLE_TRACE = json.load(f)
with open('sample_snapshot.json', 'r') as f :
    SAMPLE_SNAPSHOT = json.load(f)

IDS = {
    'user'      : 'some_user_uuid',
    'span'      : 'some_span_uuid',
    'trace'     : 'some_trace_uuid',
    'session'   : 'some_session_uuid',
    'chapter'   : 'some_chapter_uuid',
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
    except :
        pytest.fail("doesn't parse by json")

def assert_isnt_empty(r):
    j = r.json()
    assert len(j)

def assert_is_ok(r):
    assert r.ok, "returned error code"

def assert_is_not_ok(r):
    assert not r.ok, "returned success but shouldnt"

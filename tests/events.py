import json
import pytest
import requests

from helpers import *

# -------------------------------------------------------------------- #

@pytest.fixture(autouse=True)
def run_around_tests():
    setup_insert_sample_trace()
    yield

# -------------------------------------------------------------------- #

# -------------------------------
# test basic OK and response type
# -------------------------------

def test_get_events():
    r = requests.get(url('/events'))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)

def test_get_events_by_chapter():
    r = requests.get(url(f"/events_by_chapter/{IDS['chapter']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)
    assert_isnt_empty(r)

def test_get_events_by_session():
    r = requests.get(url(f"/events_by_session/{IDS['session']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)
    assert_isnt_empty(r)

def test_post_events_works():
    r = requests.post(url('/events'), data=json.dumps(SAMPLE_TRACE), headers=HEADERS)
    assert_is_ok(r)
    assert r.text == 'Creation was successful'

# ----------------
# test error codes
# ----------------

def test_get_events_by_chapter_returns_array_when_id_does_not_exist():
    r = requests.get(url(f"/events_by_chapter/{IDS['nonexist']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_empty_array(r)

def test_get_events_by_session_returns_array_when_id_does_not_exist():
    r = requests.get(url(f"/events_by_session/{IDS['nonexist']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_empty_array(r)

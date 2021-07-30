import json
import pytest
import requests

from helpers import *

# -------------------------------------------------------------------- #

@pytest.fixture(autouse=True)
def run_around_tests():
    setup_insert_sample_snapshot()
    yield

# -------------------------------------------------------------------- #

# -------------------------------
# test basic OK and response type
# -------------------------------

def test_get_snapshots():
    r = requests.get(url('/events/snapshots'))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)

def test_get_events_by_session():
    r = requests.get(url(f"/events/snapshots_by_session/{IDS['session']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)
    assert_isnt_empty(r)

def test_post_snapshots_works():
    r = requests.post(url('/events/snapshots'), data=json.dumps(SAMPLE_SNAPSHOT))
    assert_is_ok(r)
    assert r.text == 'Creation was successful'

# ----------------
# test error codes
# ----------------

def test_get_snapshots_by_session_returns_array_when_id_does_not_exist():
    r = requests.get(url(f"/events/snapshots_by_session/{IDS['nonexist']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert r.text == '[]', "returned nonempty array"

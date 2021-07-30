import json
import pytest
import requests

from helpers import *

# -------------------------------------------------------------------- #

@pytest.fixture(autouse=True)
def run_around_tests():
    setup_insert_sample_span()
    yield
    teardown_delete_sample_span()

# -------------------------------------------------------------------- #

# -------------------------------
# test basic OK and response type
# -------------------------------

def test_get_events():
    r = requests.get(url('/events'))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)

def test_get_chapters_by_session():
    r = requests.get(url(f"/chapter_ids_by_session/{IDS['session']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)
    assert_isnt_empty(r)
    assert_contents_are_type(r, dict)

def test_post_chapters_by_trigger():
    r = requests.post(url(f"/chapter_ids_by_trigger"), data=IDS['trigger'])
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)
    assert_isnt_empty(r)
    assert_contents_are_type(r, str)

# ----------------
# test error codes
# ----------------

def test_get_chapters_by_session_returns_empty_array_when_id_does_not_exist():
    r = requests.get(url(f"/chapter_ids_by_session/{IDS['nonexist']}"))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_empty_array(r)

def test_get_chapters_by_trigger_returns_array_of_empty_string_when_id_does_not_exist():
    r = requests.post(url(f"/chapter_ids_by_trigger"), data=IDS['nonexist'])
    assert_is_ok(r)
    assert_is_json_array(r)
    assert r.text == '[""]'

def test_get_chapters_by_trigger_returns_message_when_no_body_supplied():
    r = requests.post(url(f"/chapter_ids_by_trigger"), data="")
    assert r.status_code == 400
    assert r.text == 'No content supplied in the request body'

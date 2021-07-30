import json
import pytest
import requests

from helpers import *

# -------------------------------------------------------------------- #

@pytest.fixture(autouse=True)
def run_around_tests():
    setup_insert_sample_trace()
    yield
    teardown_delete_sample_trace()

# -------------------------------------------------------------------- #

# -------------------------------
# test basic OK and response type
# -------------------------------

def test_get_trigger_routes():
    r = requests.get(url('/trigger_routes'))
    assert_is_ok(r)
    assert_is_json_array(r)
    assert_is_json(r)

# ----------------
# test error codes
# ----------------

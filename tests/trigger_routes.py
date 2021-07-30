import json
import pytest
import requests

from helpers import *

# -------------------------------------------------------------------- #

@pytest.fixture(autouse=True)
def run_around_tests():
    # (todo: use a cassandra driver to input this?)
    requests.post(url('/events'), data=json.dumps(SAMPLE_TRACE), headers=HEADERS)
    yield

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

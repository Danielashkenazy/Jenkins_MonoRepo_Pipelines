from main import compute_total

def test_compute_total_basic():
    assert compute_total([10, 20, -5]) == 30

def test_compute_total_empty():
    assert compute_total([]) == 0

def test_compute_total_negative_only():
    assert compute_total([-10, -2]) == 0

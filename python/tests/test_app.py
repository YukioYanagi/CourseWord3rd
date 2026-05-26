from fastapi.testclient import TestClient

from app import app, transform

client = TestClient(app)


def test_health() -> None:
    r = client.get("/health")
    assert r.status_code == 200
    assert r.json().get("status") == "ok"


def test_transform_json_to_xml() -> None:
    body = {
        "source_format": "json",
        "target_format": "xml",
        "payload": '{"x": 1}',
    }
    r = client.post("/transform", json=body)
    assert r.status_code == 200
    data = r.json()
    assert "result" in data
    assert "<" in data["result"]


def test_transform_invalid_soap() -> None:
    body = {
        "source_format": "soap",
        "target_format": "json",
        "payload": "<root/>",
    }
    r = client.post("/transform", json=body)
    assert r.status_code == 400


def test_transform_same_format() -> None:
    s = '{"a": 1}'
    out = transform("json", "json", s)
    assert out == s


def test_transform_rejects_xml_entities() -> None:
    payload = """<?xml version="1.0"?>
<!DOCTYPE root [
<!ENTITY xxe SYSTEM "file:///etc/passwd">
]>
<root>&xxe;</root>"""
    body = {
        "source_format": "xml",
        "target_format": "json",
        "payload": payload,
    }
    r = client.post("/transform", json=body)
    assert r.status_code == 400

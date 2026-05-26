"""Трансформация данных между JSON, XML и SOAP (обёртка SOAP 1.1)."""

from __future__ import annotations

import json
import re
from typing import Any

import xml.etree.ElementTree as ET
from defusedxml.ElementTree import fromstring
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

app = FastAPI(title="Transform service", version="1.0.0")

SOAP_NS = "http://schemas.xmlsoap.org/soap/envelope/"


class TransformIn(BaseModel):
    source_format: str = Field(..., pattern="^(json|xml|soap)$")
    target_format: str = Field(..., pattern="^(json|xml|soap)$")
    payload: str


class TransformOut(BaseModel):
    result: str
    error: str | None = None


def _json_to_obj(payload: str):
    return json.loads(payload)


def _obj_to_json(obj: Any) -> str:
    return json.dumps(obj, ensure_ascii=False, indent=2)


def _elem_to_obj(elem: ET.Element) -> Any:
    children = list(elem)
    if not children:
        t = (elem.text or "").strip()
        return t if t else None
    if all(c.tag == children[0].tag for c in children) and len({c.tag for c in children}) == 1:
        return [_elem_to_obj(c) for c in children]
    out: dict[str, Any] = {}
    for c in children:
        v = _elem_to_obj(c)
        tag = c.tag.split("}")[-1] if "}" in c.tag else c.tag
        if tag in out:
            if not isinstance(out[tag], list):
                out[tag] = [out[tag]]
            out[tag].append(v)
        else:
            out[tag] = v
    return out


def _xml_to_obj(xml_str: str) -> Any:
    root = fromstring(xml_str)
    return {root.tag.split("}")[-1]: _elem_to_obj(root)}


def _obj_to_xml(obj: Any, root_tag: str = "root") -> str:
    def build(parent: ET.Element, data: Any, name: str) -> None:
        if isinstance(data, dict):
            for k, v in data.items():
                tag = re.sub(r"[^\w.\-]", "_", str(k))
                if not tag:
                    tag = "item"
                el = ET.SubElement(parent, tag)
                build(el, v, tag)
        elif isinstance(data, list):
            for item in data:
                el = ET.SubElement(parent, name or "item")
                build(el, item, name)
        elif data is None:
            return
        else:
            parent.text = str(data)

    root = ET.Element(root_tag)
    if isinstance(obj, dict) and len(obj) == 1:
        only_key = next(iter(obj))
        sub = ET.SubElement(root, re.sub(r"[^\w.\-]", "_", only_key) or "item")
        build(sub, obj[only_key], only_key)
    else:
        build(root, obj, root_tag)
    return '<?xml version="1.0" encoding="utf-8"?>\n' + ET.tostring(root, encoding="unicode")


def _unwrap_soap(xml_str: str) -> str:
    tree = fromstring(xml_str)
    tag_lower = tree.tag.lower()
    if "envelope" not in tag_lower:
        raise ValueError("not a SOAP Envelope")
    body = None
    for child in tree:
        ctag = child.tag.split("}")[-1].lower()
        if ctag == "body":
            body = child
            break
    if body is None or not len(body):
        raise ValueError("SOAP Body is empty")
    inner = body[0]
    return ET.tostring(inner, encoding="unicode")


def _wrap_soap(inner_xml: str) -> str:
    inner = inner_xml.strip()
    inner = re.sub(r"^<\?xml[^>]*\?>\s*", "", inner, count=1)
    return (
        '<?xml version="1.0" encoding="utf-8"?>\n'
        f'<s:Envelope xmlns:s="{SOAP_NS}">\n'
        f"  <s:Body>\n{inner}\n  </s:Body>\n"
        "</s:Envelope>"
    )


def transform(source: str, target: str, payload: str) -> str:
    source, target = source.lower(), target.lower()
    if source == target:
        return payload

    if source == "json":
        obj = _json_to_obj(payload)
        if target == "xml":
            return _obj_to_xml(obj)
        if target == "soap":
            return _wrap_soap(_obj_to_xml(obj))
        raise ValueError("unsupported target")

    if source == "xml":
        obj = _xml_to_obj(payload)
        if target == "json":
            return _obj_to_json(obj)
        if target == "soap":
            return _wrap_soap(payload)
        raise ValueError("unsupported target")

    if source == "soap":
        inner = _unwrap_soap(payload)
        if target == "xml":
            return inner
        if target == "json":
            return _obj_to_json(_xml_to_obj(inner))
        raise ValueError("unsupported target")

    raise ValueError("unsupported source")


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok", "service": "python-transform"}


@app.post("/transform", response_model=TransformOut)
def do_transform(body: TransformIn) -> TransformOut:
    try:
        out = transform(body.source_format, body.target_format, body.payload)
        return TransformOut(result=out)
    except Exception as e:
        print("ERROR:", repr(e))
        raise HTTPException(status_code=400, detail=str(e)) from e

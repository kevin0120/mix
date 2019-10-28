# -*- coding: utf-8 -*-

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from odoo.http import Response
import logging
import json

_logger = logging.getLogger(__name__)

MAGIC_ERROR_CODE = 59999  # 永远不会定义

SA_ERROR_CODE = {
    40001: "Tightening Tool Is Not Success Configuration"
}


def sa_http_session(
        retries=5,
        backoff_factor=0.3,
        status_forcelist=(500, 502, 504),
        session=None,
):
    session = session or requests.Session()
    retry = Retry(
        total=retries,
        read=retries,
        connect=retries,
        backoff_factor=backoff_factor,
        status_forcelist=status_forcelist,
    )
    adapter = HTTPAdapter(max_retries=retry)
    session.mount('http://', adapter)
    session.mount('https://', adapter)
    return session


def sa_success_resp(status_code=200, **kwargs):
    if status_code > 400:
        _logger.error("Success Response The Status Code Must Be Less Than 400, But Now Is {0}".format(status_code))
        return Response(200)
    data = {
        "msg": kwargs.get("msg") or kwargs.get("message"),
        "extra": kwargs.get("extra", "") or kwargs.get("extra_info", "")
    }
    body = json.dumps(data)
    headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
    resp = Response(body, status=status_code, headers=headers)
    return resp


def sa_fail_response(error_code=MAGIC_ERROR_CODE, **kwargs):
    msg = SA_ERROR_CODE.get(error_code)
    if not msg:
        _logger.error("Error Code: {0} Is Not Defined".format(error_code))
        msg = kwargs.get("msg") or kwargs.get("message"),
    data = {
        "msg": msg,
        "extra": kwargs.get("extra", "") or kwargs.get("extra_info", "")
    }
    body = json.dumps(data)
    headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
    return Response(body, status=400, headers=headers)

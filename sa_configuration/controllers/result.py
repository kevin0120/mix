# -*- coding: utf-8 -*-

from .base import BaseRestClass
import json

class Result(BaseRestClass):
    def __init__(self):
        self.id = None
        self.name = False

    @classmethod
    def convert_from_quailty_check(cls, quality_check):
        ret = cls()
        ret.name = quality_check['name']
        ret.id = quality_check['id']
        return ret


class ResultList(BaseRestClass):
    def __init__(self):
        self.result_list = []

    def append_result(self, result):
        self.result_list.append(result.__dict__)

    def to_json(self):
        return json.dumps(self.result_list)
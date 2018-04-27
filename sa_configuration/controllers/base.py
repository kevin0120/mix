# -*- coding: utf-8 -*-

import json

class BaseRestClass(object):
    def validation(self):
        return True

    def to_json(self):
        if self.validation():
            return json.dumps(self.__dict__)
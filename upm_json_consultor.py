#!/usr/bin/env python
# -*- coding: utf-8 -*-
import json
import io
import objectpath
from requests import get


def get_etsisi_degrees_info_json():
    with io.open('data.json', 'r', encoding='utf8') as data_file:
        data = json.load(data_file)

    jsons_joint = {}
    for i, degree in enumerate(data["upm_jsons"]):
        jsons_joint[degree] = get(data["upm_jsons"][degree]).json()
        # print jsons_joint[degree]

    with io.open('etsisi_degrees_data.json', 'w', encoding='utf8') as data_file:
        json.dump(jsons_joint, data_file)

    return jsons_joint
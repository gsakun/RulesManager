# -*- coding: utf-8 -*-
from .File import write_yaml_file, get_rules, get_detail, update_yaml_file, delete_yaml_file
import yaml
import os

def create_rules_model(groupname, alertname, expr, _for, labels, annotations):
    try:
        data = {'groups': [{
            "name": groupname,
            "rules": [{
                "alert": alertname,
                "expr": expr,
                "for": _for,
                "labels": labels,
                "annotations": annotations
                }]
            }]
        }
        filename = groupname + "____" + alertname
        rs = write_yaml_file(filename, data)
        return rs
    except Exception as e:
        print(e)
        return 1


def update_rules_model(groupname, alertname, expr, _for, labels, annotations):
    try:
        data = {'groups': [{
            "name": groupname,
            "rules": [{
                "alert": alertname,
                "expr": expr,
                "for": _for,
                "labels": labels,
                "annotations": annotations
                }]
            }]
        }
        filename = groupname + "____" + alertname
        rs = update_yaml_file(filename, data)
        return rs
    except Exception as e:
        print(e)
        return 1


def delete_rules_model(groupname,alertname):
    try:
        filename = groupname + "____" + alertname
        rs = delete_yaml_file(filename)
        return rs
    except Exception as e:
        print(e)
        return 1


def get_rule():
    rs = get_rules()
    return rs


def get_rule_detail(name):
    return get_detail(name)
# -*- coding: utf-8 -*-
from .Config import filepath
import os
import yaml


def write_yaml_file(filename, data):

    if filepath.endswith('/'):
        rule_file = filepath + filename + '.yml'
    else:
        rule_file = filepath + '/' + filename + '.yml'
    exist = os.path.exists(rule_file)
    if exist:
        return 2
    try:
        with open(rule_file, 'w', encoding='utf-8') as f:
            yaml.dump(data,f)
        return 0
    except Exception as e:
        print(e)
        return 1

def update_yaml_file(filename, data):

    if filepath.endswith('/'):
        rule_file = filepath + filename + '.yml'
    else:
        rule_file = filepath + '/' + filename + '.yml'
    exist = os.path.exists(rule_file)
    if not exist:
        print("this rule not exist,create it")
    try:
        with open(rule_file, 'w', encoding='utf-8') as f:
            f.truncate()
            yaml.dump(data,f)
        return 0
    except Exception as e:
        print(e)
        return 1

def delete_yaml_file(_name):

    if filepath.endswith('/'):
        _file = filepath + _name + '.yml'
    else:
        _file = filepath + '/' + _name + '.yml'
    try:
        if(os.path.exists(_file)):
            os.remove(_file)
            return 0
        else:
            return 2
    except Exception as e:
        print(e)
        return 1
        
def get_rules():
    rs = []
    try:
        for root, dirs, files in os.walk(filepath):
            rs.append(files)
    except Exception as e:
        print(e)
    return rs


def get_detail(filename):
    try:
        if filepath.endswith('/'):
            rule_file = filepath + filename + '.yml'
        else:
            rule_file = filepath + '/' + filename + '.yml'

        content = os.popen('head -1 ' + rule_file).read()
        values = eval(content[1:len(content)])
        return values
    except Exception as e:
        print(e)
        return {'name': '', 'alert': '', 'expr': '', '_for': '', 'level': '', 'summary': '', 'description': ''}

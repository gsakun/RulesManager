# -*- coding: utf-8 -*-
from flask import Flask, request, render_template, make_response, Response
import sys
import json
from func.Func import create_rules_model, get_rule, get_rule_detail, update_rules_model, delete_rules_model
import func.Config
import os

app = Flask(__name__)


@app.route('/ruler/add', methods=['POST'])
def add():
    content = json.dumps('create successful')
    data = request.get_data()
    info = json.loads(data)
    try:
        groups = info['groups']
        for group in groups:
            for rule in group['rules']:
                rs = create_rules_model(group['name'],rule['alert'],rule['expr'],rule['for'],rule['labels'],rule['annotations'])
                if rs == 2:
                    return json.dumps('error! this rule already exist')
    except Exception as e:
        print (e)
        content = json.dumps('create failed')

    resp = Response(content)
    return resp


@app.route('/ruler/update', methods=['POST'])
def update():
    content = json.dumps('modify successful')
    data = request.get_data()
    info = json.loads(data)
    try:
        groups = info['groups']
        for group in groups:
            for rule in group['rules']:
                rs = update_rules_model(group['name'],rule['alert'],rule['expr'],rule['for'],rule['labels'],rule['annotations'])
                if rs == 1:
                    return json.dumps('error! please check the path is exist and Permission!')
    except Exception as e:
        print (e)
        content = json.dumps('error! please check the path is exist and Permission!')

    resp = Response(content)
    return resp


@app.route('/ruler/delete', methods=['POST'])
def delete():
    content = json.dumps('delete successful')
    data = request.get_data()
    info = json.loads(data)
    try:
        groups = info['groups']
        for group in groups:
            for rule in group['rules']:
                rs = delete_rules_model(group['name'],rule['alert'])
                if rs == 2:
                    return json.dumps('error! this rule not exist')
    except Exception as e:
        print (e)
        content = json.dumps('error! please check the path is exist and Permission!')
    resp = Response(content)
    return resp


@app.route('/ruler/query', methods=['GET'])
def get_rules_list():
    try:
        content = get_rule()
    except Exception as e:
        print (e)
        content = []

    resp = Response(json.dumps(content))
    return resp

if __name__ == '__main__':
    app.run("0.0.0.0", 8888)

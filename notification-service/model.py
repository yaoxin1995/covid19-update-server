from flask import Flask, g
from json import JSONEncoder
from datetime import datetime
from app import db

Column = db.Column
Integer = db.Integer
DateTime = db.DateTime
Text = db.Text


class ModelJSONEncoder(JSONEncoder):
    def default(self, o):
        _dict = {}
        if isinstance(o, Notification):
            for var_name in Notification.__dict__:
                if str(var_name).startswith('_'):
                    continue
                _dict[var_name] = getattr(o, var_name)
            return _dict
        elif isinstance(o, datetime):
            return f'{o}'
        return JSONEncoder.default(self, o)


class Notification(db.Model):
    id = Column(Integer, primary_key=True)
    creation_date = Column(DateTime, nullable=False)
    channel = Column(Text, nullable=False)
    recipient = Column(Text, nullable=False)
    msg = Column(Text, nullable=False)
    error_msg = Column(Text, nullable=True)

    def __str__(self):
        return f'Notification {self.id}: {self.creation_date} - {self.channel} - {self.recipient} - {self.msg}'

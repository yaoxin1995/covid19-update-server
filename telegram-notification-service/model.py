from sqlalchemy import inspect
from json import JSONEncoder
from datetime import datetime
from app import db

Column = db.Column
Integer = db.Integer
DateTime = db.DateTime
Text = db.Text


class ModelDictJSONEncoder(JSONEncoder):
    """
    This JSON encoder is capable to convert all data types
    to json which are used in notification model.
    """
    def default(self, o):
        _dict = {}
        if isinstance(o, Notification):
            for column_name in o.column_names:
                _dict[column_name] = getattr(o, column_name)
            return _dict
        elif isinstance(o, datetime):
            return f'{o}'
        return JSONEncoder.default(self, o)


class Notification(db.Model):
    id = Column(Integer, primary_key=True)
    creation_date = Column(DateTime, nullable=False)
    recipient = Column(Text, nullable=False)
    msg = Column(Text, nullable=False)
    error_msg = Column(Text, nullable=True)
    error_msg_human_readable = Column(Text, nullable=True)

    def fetch(self, **filter_attr):
        query_result = self.query.filter_by(**filter_attr).all()
        notifications = []
        for query in query_result:
            notifications.append(query)
        return notifications

    def save(self):
        db.session.add(self)
        db.session.commit()

    @property
    def as_dict(self):
        return {c.key: getattr(self, c.key)
                for c in inspect(self).mapper.column_attrs}

    @property
    def column_names(self):
        return self.__table__.columns.keys()

    def __str__(self):
        return f'Notification {self.id}: {self.creation_date} - {self.recipient} - {self.msg}'

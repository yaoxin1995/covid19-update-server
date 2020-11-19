from app import app, db
from model import Notification, ModelJSONEncoder
from flask import Response, request
import datetime
import json


@app.route('/notification')
def get_all():
    notifications = Notification.query.all()
    json_string = json.dumps(notifications, cls=ModelJSONEncoder)
    return Response(json_string, mimetype='application/json')


@app.route('/notification/<notification_id>')
def get(notification_id=None):
    n = Notification.query.filter_by(id=notification_id).first()
    json_string = json.dumps(n, cls=ModelJSONEncoder)
    return Response(json_string, mimetype='application/json')


@app.route('/notification', methods=['POST'])
def add():
    """
    TODO: Does caller need to get immediate response about success of sending message?
    TODO: Better implement sender queue if e.g., Telegram/mail server access currently not available.
    """
    try:
        channel = str(request.form['channel'])
        recipient = str(request.form['recipient'])
        msg = str(request.form['msg'])
    except KeyError:
        return Response('Required arguments: channel, recipient, msg!', status=400)

    n = Notification(
        creation_date=datetime.datetime.utcnow(),
        channel=channel,
        recipient=recipient,
        msg=msg,
        error_msg=None
    )
    db.session.add(n)
    db.session.commit()

    return Response(json.dumps(n, cls=ModelJSONEncoder), mimetype='application/json')


@app.route('/notification', methods=['POST'])
def notification():
    pass
from app import app, db
from model import Notification, ModelJSONEncoder
from flask import Response, request
from notification_providers.channel_telegram import ChannelTelegram
import datetime
import json


@app.route('/notification', methods=['GET', 'POST'])
def notification():
    """
    Return all notifications using JSON format on a GET request.
    Create a new notification object and send the message. Return this object using JSON format on a POST request.
    TODO: Get requests by recipient
    """
    if request.method == 'GET':
        notifications = Notification.query.all()
    else:
        try:
            channel = str(request.form['channel'])
            recipient = str(request.form['recipient'])
            msg = str(request.form['msg'])
        except KeyError:
            return Response('Required arguments: channel, recipient, msg!', status=400)

        telegram = ChannelTelegram()
        telegram_response, human_readable_error_msg = telegram.send_message(chat_id=recipient, msg=msg)

        n = Notification(
            creation_date=datetime.datetime.utcnow(),
            channel=channel,
            recipient=recipient,
            msg=msg,
            error_msg=json.dumps(telegram_response),
            error_msg_human_readable=human_readable_error_msg
        )

        db.session.add(n)
        db.session.commit()
        notifications = [n]
    return Response(json.dumps(notifications, cls=ModelJSONEncoder), mimetype='application/json')


@app.route('/notification/<notification_id>')
def get_notification(notification_id=None):
    """
    Return a notification by id using JSON format.
    """
    n = Notification.query.filter_by(id=notification_id).first()
    json_string = json.dumps(n, cls=ModelJSONEncoder)
    return Response(json_string, mimetype='application/json')

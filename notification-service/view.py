from app import app, db
from model import Notification, ModelJSONEncoder
from flask import Response, request
from notification_providers.telegram_notification import TelegramNotifier
import datetime
import json


@app.route('/notification', methods=['GET', 'POST'])
def notification():
    """
    Return all notifications using JSON format on a GET request.
    Create a new notification object and send the message. Return this object using JSON format on a POST request.
    """
    if request.method == 'GET':
        notifications = Notification.query.all()
    else:
        # TODO: Does caller need to get immediate response about success of sending message?
        # TODO: Better implement sender queue if e.g., Telegram/mail server access currently not available.
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

        tn = TelegramNotifier()
        tn.send_message(chat_id=recipient, msg=msg)
        # TODO: Do some error checking here.
        # n.error_msg =

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

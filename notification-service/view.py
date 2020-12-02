from app import app, db
from model import Notification, ModelJSONEncoder
from flask import Response, request
from notification_providers.channel_telegram import ChannelTelegram
import datetime
import json


@app.route('/notification', methods=['GET', 'POST'])
def notification():
    """
    On a GET request:
      Return notifications using JSON format.
      If recipient is set return all notifications for this recipient.
    On a POST request:
      Create a new notification object and send the message via the Telegram bot.
    """
    status_code = 404
    if request.method == 'GET':
        recipient = None
        try:
            recipient = request.args.get('recipient')
        except KeyError:
            pass

        if recipient:
            # Get messages from a specific recipient
            notifications = Notification.query.filter_by(recipient=recipient).all()
            if len(notifications) > 0:  # There might be no notification available by this user
                status_code = 200
        else:
            # Get all messages
            notifications = Notification.query.all()
            # Use always code 200 because this is the general endpoint for getting notifications
            status_code = 200
    else:
        try:
            channel = str(request.form['channel'])
            recipient = str(request.form['recipient'])
            msg = str(request.form['msg'])
        except KeyError:
            return Response('Required arguments: channel, recipient, msg!', mimetype='application/json', status=400)

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
        notifications = n
        status_code = 201
    return Response(json.dumps(notifications, cls=ModelJSONEncoder), mimetype='application/json', status=status_code)


@app.route('/notification/<notification_id>')
def get_notification(notification_id=None):
    """
    Return a notification by id using JSON format.
    """
    status_code = 404
    n = Notification.query.filter_by(id=notification_id).first()
    if n:
        status_code = 200
    json_string = json.dumps(n, cls=ModelJSONEncoder)
    return Response(json_string, mimetype='application/json', status=status_code)

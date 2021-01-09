from app import app, Config
from flask import request
from resources import NotificationResource
from helpers import get_params_if_in_query, get_params_if_in_form


@app.route(Config.NOTIFICATION_BASE_ROUTE, methods=['GET', 'POST'])
def notification():
    """
    On a GET request:
      Return notifications using JSON format.
      If recipient is set return all notifications for this recipient.
    On a POST request:
      Create a new notification object and send the message via the Telegram bot.
    """
    resource = NotificationResource(request)
    if request.method == 'GET':
        filter_param = get_params_if_in_query(request, ['recipient'])
        resource.resource_fetch(filter_param)
    elif request.method == 'POST':
        params = get_params_if_in_form(request, ['recipient', 'msg'])
        resource.resource_create(params)
    return resource.response


@app.route(f"{Config.NOTIFICATION_BASE_ROUTE}/<notification_id>", methods=['GET'])
def get_notification(notification_id):
    """
    Return a notification by id using JSON format.
    TODO: Add support for deleting notifications (by recipient).
    """
    resource = NotificationResource(request)
    resource.resource_fetch({'id': notification_id}, single_resource=True)
    return resource.response

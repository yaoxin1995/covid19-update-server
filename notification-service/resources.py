import json
import datetime
from app import Config
from flask import Response
from channel_telegram import ChannelTelegram
from model import Notification, ModelDictJSONEncoder
from json_hal import NotificationDocument, NotificationListDocument


class NotificationResource:
    def __init__(self, request=None):
        self.__request = request
        # response data attributes
        self.__response = Response()
        self.__status_code = 500
        self.__response_data = None
        self.__response.content_type = 'text/plain'

    def resource_create(self, params):
        # Prepare for case that resource will be successfully created
        self.__status_code = 201
        try:
            recipient = params['recipient']
            msg = params['msg']
        except KeyError:
            self.__response_data = {
                "code": "Required parameters: recipient, msg!",
                "description": "It is required to send both parameters recipient and msg via"
                               " application/x-www-form-urlencoded!"
            }
            self.__status_code = 400
            self.__make_response_by_content_type()
            return
        telegram = ChannelTelegram()
        telegram_response, human_readable_error_msg = telegram.send_message(chat_id=recipient, msg=msg)
        n = Notification(
            creation_date=datetime.datetime.utcnow(),
            recipient=recipient,
            msg=msg,
            error_msg=json.dumps(telegram_response),
            error_msg_human_readable=human_readable_error_msg
        )
        n.save()
        self.__response_data = n
        self.__make_response_by_content_type()

    def resource_fetch(self, filter_params=None, single_resource=False):
        if not filter_params:
            filter_params = {}
        # Prepare for case that a list of resources is requested
        self.__status_code = 200
        n = Notification()
        # print("GET", filter_params)
        notifications = n.fetch(**filter_params)
        if single_resource:
            if len(notifications) == 0:
                # We don't want to return an empty list of resources
                notifications = None
                # If we want to access single resource which doesn't exist.
                self.__status_code = 404
            elif len(notifications) == 1:
                # We don't want to return a list of resources, because we just have one.
                notifications = notifications[0]
        # print("got", notifications)
        self.__response_data = notifications
        self.__make_response_by_content_type()

    def __make_response_by_content_type(self):
        accept_header_list = str(self.__request.headers.get('accept')).split(',')
        # print(accept_header_list)
        if Config.JSON_HAL_MIME_TYPE in accept_header_list:
            self.__make_hal_response()
        elif Config.JSON_MIME_TYPE in accept_header_list:
            self.__make_json_response()
        else:
            self.__response_data = 'No supported content type was found!'
            self.__status_code = 406

    def __make_hal_response(self):
        data = self.__response_data
        if type(data) is list:
            data = NotificationListDocument(data).to_dict()
        else:
            data = NotificationDocument(data).to_dict()
        self.__response_data = json.dumps(data, cls=ModelDictJSONEncoder)
        self.__response.content_type = Config.JSON_HAL_MIME_TYPE

    def __make_json_response(self):
        self.__response_data = json.dumps(self.__response_data, cls=ModelDictJSONEncoder)
        self.__response.content_type = Config.JSON_MIME_TYPE

    @property
    def response(self):
        self.__response.data = self.__response_data
        self.__response.status_code = self.__status_code
        return self.__response

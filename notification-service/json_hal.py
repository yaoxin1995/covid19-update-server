from app import Config
from flask_hal.document import BaseDocument, Embedded, link
from model import Notification


class NotificationDocument(BaseDocument):
    def __init__(self, notification=None):
        data = None
        links = link.Collection()

        if notification and isinstance(notification, Notification):
            # add self link
            links.append(link.Link('self', f"{Config.NOTIFICATION_BASE_ROUTE}/{notification.id}"))
            data = notification.as_dict
        super().__init__(data=data, links=links)


class NotificationListDocument(BaseDocument):
    # TODO: Implement next to prevent json documents from being too long
    def __init__(self, notifications_list):
        links = link.Collection()

        # add self link
        links.append(link.Link('self', f"{Config.NOTIFICATION_BASE_ROUTE}"))

        ndoc_list = []
        for n in notifications_list:
            ndoc_list.append(NotificationDocument(n))
        embedded = Embedded(data=ndoc_list)
        super().__init__(embedded={'notification': embedded}, links=links)

from flask_hal.document import BaseDocument, Document, Embedded, link

# TODO: WIP
class NotificationDocument(BaseDocument):
    def __init__(self, n):
        data = {}
        links = link.Collection()

        # add self link
        links.append(link.Link('self', f"/notification/{n.id}"))

        if isinstance(n, Notification):
            for var_name in Notification.__dict__:
                if str(var_name).startswith('_'):
                    continue
                attr = getattr(n, var_name)
                if isinstance(attr, datetime.datetime):
                    # convert datetime obj into string
                    attr = f'{attr}'
                data[var_name] = attr
        super().__init__(data=data, links=links)


class NotificationList(BaseDocument):
    # TODO: Implement next to prevent json documents from being too long
    def __init__(self, notifications_list):
        ndoc_list = []
        for n in notifications_list:
            ndoc_list.append(NotificationDocument(n))
        super().__init__(data={'total': len(notifications_list)}, embedded={'notifications': ndoc_list[0]})

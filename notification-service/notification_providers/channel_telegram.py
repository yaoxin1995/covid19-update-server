import requests
from app import Config


class ChannelTelegram:
    """
    Provide Telegram bot functionality.
    """
    def __init__(self):
        self._bot_token = Config.TELEGRAM_BOT_TOKEN

    def send_message(self, chat_id, msg):
        """
        Send messages via Telegram bot to specified recipient identified by chat id.
        """
        human_readable_error_desc = None

        response = self._perform_request('sendMessage', {'chat_id': chat_id, 'text': msg})

        t_response = response['telegram_response']
        if t_response:
            if not t_response['ok']:
                # Telegram server returned error
                human_readable_error_desc = t_response['description']
        elif response['request_error']:
            human_readable_error_desc = 'Unknown error.'
        return response, human_readable_error_desc

    def _build_url(self, method_name):
        return f"https://api.telegram.org/bot{self._bot_token}/{method_name}"

    def _perform_request(self, method_name, params):
        request_error = None
        status_code = None
        telegram_response = None
        try:
            r = requests.get(self._build_url(method_name), params=params, timeout=3)
            telegram_response = r.json()
            status_code = r.status_code
        except requests.exceptions.ConnectionError:
            request_error = 'Connection error'
        except requests.exceptions:
            request_error = 'General network error'

        response = {
            'request_error': request_error,
            'status_code': status_code,
            'telegram_response': telegram_response
        }
        return response

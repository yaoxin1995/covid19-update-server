import os
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
import logging
from jinja2 import Template


class TelegramResponder:
    def __init__(self):
        logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.INFO)

        bot_token = os.environ['TELEGRAM_BOT_TOKEN']
        welcome_message = str(os.environ['WELCOME_MESSAGE'])
        unknown_message = str(os.environ['UNKNOWN_MESSAGE'])
        self.updater = Updater(token=bot_token)
        self.dispatcher = self.updater.dispatcher

        # Define some callback functions
        def callback_start(update, context):
            self.send_message(context, message=welcome_message, user_attributes=self.get_user_attributes(update))

        def callback_unknown(update, context):
            # print(update.effective_chat.id)
            self.send_message(context, message=unknown_message, user_attributes=self.get_user_attributes(update))

        # Register handlers to react to user messages and commands
        start_handler = CommandHandler('start', callback_start)
        self.dispatcher.add_handler(start_handler)
        unknown_handler = MessageHandler(Filters.text, callback_unknown)
        self.dispatcher.add_handler(unknown_handler)

        # Print messages with example user vars
        demo_vars = {
            'USER_FIRST_NAME': 'Max',
            'USER_LAST_NAME': 'Mustermann',
            'USER_FULL_NAME': 'Max Mustermann',
            'USER_USERNAME': 'MMuster',
            'USER_CHAT_ID': '123456789',
        }
        message = render_message(welcome_message, demo_vars, verbose=True)
        print(f"--------- DEMO OF RENDERED WELCOME MESSAGE ---------\n'{message}'")
        message = render_message(unknown_message, demo_vars, verbose=False)
        print(f"----- DEMO OF RENDERED UNKNOWN REQUEST MESSAGE -----\n'{message}'")

    @staticmethod
    def get_user_attributes(update):
        user = update.effective_user
        chat_id = update.effective_chat.id
        return {
            'USER_FIRST_NAME': user.first_name,
            'USER_LAST_NAME': user.last_name,
            'USER_FULL_NAME': user.full_name,
            'USER_USERNAME': user.username,
            'USER_CHAT_ID': chat_id,
        }

    @staticmethod
    def send_message(context, message, user_attributes):
        # Render template
        message = render_message(message, predefined_template_vars=user_attributes)
        context.bot.send_message(chat_id=user_attributes['USER_CHAT_ID'], text=message)

    def start_polling(self):
        # Check for bot messages every x seconds from Telegram
        self.updater.start_polling(poll_interval=5)
        self.updater.idle()
        #self.updater.stop()


def render_message(message, predefined_template_vars, verbose=False):
    """
    This renders the message by using predefined_template_vars and environment variables which have the prefix
    `TEMPLATE_VAR_`.
    :return: rendered message
    """
    template_var_start_pattern = 'TEMPLATE_VAR_'
    # Some template vars are already predefined (e.g., chat_id and Telegram user attributes)
    template_vars = dict(predefined_template_vars)
    # Collect env vars starting with template_var_start_pattern
    for env_var_name in os.environ.keys():
        env_var_name = str(env_var_name)
        # If we have an env variable matching our start pattern add it to the dict
        if env_var_name.startswith(template_var_start_pattern):
            template_var_name = env_var_name[len(template_var_start_pattern):]
            template_var_value = os.environ[env_var_name]
            template_vars[template_var_name] = template_var_value
    if verbose:
        print(f"Using template vars: {template_vars}.")
    # Render template using Jinja2
    template = Template(message)
    return template.render(template_vars)


if __name__ == '__main__':
    tn = TelegramResponder()
    tn.start_polling()

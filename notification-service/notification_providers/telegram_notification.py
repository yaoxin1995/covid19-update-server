import os
from telegram import Update
from telegram.ext import Updater, CommandHandler, MessageHandler, CallbackContext, Filters
import logging


class TelegramNotifier:
    def __init__(self):
        logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.INFO)

        # TODO: How to handle bot token?!
        bot_token = os.environ['TELEGRAM_BOT_TOKEN']
        self.updater = Updater(token=bot_token)
        self.dispatcher = self.updater.dispatcher

        # Define some callback functions
        def callback_start(update, context):
            context.bot.send_message(chat_id=update.effective_chat.id, text="Welcome! You are now able to configure "
                                                                            "COVID-19 incidence update notifications via "
                                                                            "Telegram.")

        def callback_unknown(update, context):
            # print(update.effective_chat.id)
            context.bot.send_message(chat_id=update.effective_chat.id, text="Sorry, I didn't understand you.")

        # Register handlers to react to user messages and commands
        start_handler = CommandHandler('start', callback_start)
        self.dispatcher.add_handler(start_handler)
        unknown_handler = MessageHandler(Filters.text, callback_unknown)
        self.dispatcher.add_handler(unknown_handler)

    def send_message(self, chat_id, msg):
        def callback_send_msg(context):
            params = context.job.context
            context.bot.send_message(chat_id=params['chat_id'],
                                     text=params['msg'])

        queue = self.updater.job_queue
        # Send message to chat immediately
        queue.run_once(callback_send_msg, when=0, context={'chat_id': chat_id, 'msg': msg})

    def start_polling(self):
        # Check for bot messages every x seconds from Telegram
        self.updater.start_polling(poll_interval=5)
        self.updater.idle()

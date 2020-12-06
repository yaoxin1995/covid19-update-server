import os
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
import logging


class TelegramResponder:
    def __init__(self):
        logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.INFO)

        # TODO: How to handle bot token?!
        bot_token = os.environ['TELEGRAM_BOT_TOKEN']
        self.updater = Updater(token=bot_token)
        self.dispatcher = self.updater.dispatcher

        # Define some callback functions
        def callback_start(update, context):
            chat_id = update.effective_chat.id
            context.bot.send_message(chat_id=chat_id, text="Welcome! You are now able to configure "
                                                           "COVID-19 incidence update notifications "
                                                           "via Telegram.")
            context.bot.send_message(chat_id=chat_id,
                                     text=f"Go to your dashboard and add the Telegram notification provider by "
                                          f"entering this ID {chat_id}.")

        def callback_unknown(update, context):
            # print(update.effective_chat.id)
            context.bot.send_message(chat_id=update.effective_chat.id, text="Sorry, I didn't understand you.")

        # Register handlers to react to user messages and commands
        start_handler = CommandHandler('start', callback_start)
        self.dispatcher.add_handler(start_handler)
        unknown_handler = MessageHandler(Filters.text, callback_unknown)
        self.dispatcher.add_handler(unknown_handler)

    def start_polling(self):
        # Check for bot messages every x seconds from Telegram
        self.updater.start_polling(poll_interval=5)
        self.updater.idle()
        #self.updater.stop()


if __name__ == '__main__':
    tn = TelegramResponder()
    tn.start_polling()

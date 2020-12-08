import os
from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from telegram_responder import TelegramResponder


class Config:
    """
    Some additional Flask config parameters
    """
    DEBUG = False
    SQLALCHEMY_DATABASE_URI = 'sqlite:///db.sql'
    TELEGRAM_BOT_TOKEN = os.environ['TELEGRAM_BOT_TOKEN']


app = Flask(__name__)
# load some additional config parameters
app.config.from_object(Config)

db = SQLAlchemy(app)

from view import *

# create db and tables
db.create_all()

tn = TelegramResponder()
# Run telegram responder for being able to react to user messages
tn.start_polling()

# To show messages we produced so far
# https://stackoverflow.com/questions/44405708/flask-doesnt-print-to-console
print(flush=True)

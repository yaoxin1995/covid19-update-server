import os
from flask import Flask, Response
from flask_sqlalchemy import SQLAlchemy
from flask_cors import CORS
from flask_hal import HAL, HALResponse
from telegram_responder import TelegramResponder


class Config:
    """
    Some additional Flask config parameters
    """
    DEBUG = False
    SQLALCHEMY_DATABASE_URI = 'sqlite:///./db/db.sql'
    TELEGRAM_BOT_TOKEN = os.environ['TELEGRAM_BOT_TOKEN']
    NOTIFICATION_BASE_ROUTE = '/notification'
    JSON_HAL_MIME_TYPE = 'application/hal+json'
    JSON_MIME_TYPE = 'application/json'


app = Flask(__name__)
# enabling CORS support
# TODO: Make cors' allowed origins more restrictive instead of using '*'
CORS(app)
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

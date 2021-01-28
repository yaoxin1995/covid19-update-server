import os
from flask import Flask, Response
from flask_sqlalchemy import SQLAlchemy
from flask_cors import CORS
import requests
from telegram_responder import TelegramResponder


class Config:
    """
    Some additional Flask config parameters
    """
    DEBUG = False
    SQLALCHEMY_DATABASE_URI = f"sqlite:///{os.environ['DB_FILE_PATH']}"
    TELEGRAM_BOT_TOKEN = os.environ['TELEGRAM_BOT_TOKEN']
    NOTIFICATION_BASE_ROUTE = '/notification'
    JSON_HAL_MIME_TYPE = 'application/hal+json'
    JSON_MIME_TYPE = 'application/json'
    AUTH0_ISS = os.environ['AUTH0_ISS']
    AUTH0_AUD = os.environ['AUTH0_AUD']
    try:
        request = requests.get(f"https://{AUTH0_ISS}.well-known/jwks.json", timeout=5)
        JWKS = request.json()
    except requests.exceptions.ConnectionError:
        raise Exception("Can't fetch jwks token! Please check network connectivity.")



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

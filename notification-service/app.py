import os
from flask import Flask
from flask_sqlalchemy import SQLAlchemy


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


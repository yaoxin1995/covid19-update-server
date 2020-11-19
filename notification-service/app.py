from flask import Flask
from flask_sqlalchemy import SQLAlchemy


class Config:
    """
    Some additional Flask config parameters
    """
    DEBUG = True
    SQLALCHEMY_DATABASE_URI = 'sqlite:///db.sql'


app = Flask(__name__)
# load some additional config parameters
app.config.from_object(Config)

db = SQLAlchemy(app)

from view import *

from model import *

# create db and tables
db.create_all()

if __name__ == '__main__':
    app.run()

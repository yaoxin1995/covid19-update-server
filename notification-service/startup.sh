#!/bin/bash

python3 telegram_responder.py &

export FLASK_APP=app.py
python3 -m flask run --host=0.0.0.0
#!/bin/bash

export FLASK_APP=app.py
python3 ./notification_providers/telegram_notification.py &
python3 -m flask run --host=0.0.0.0

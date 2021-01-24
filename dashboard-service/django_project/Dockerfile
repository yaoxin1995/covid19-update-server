# Use the official Python image from the Docker Hub
FROM python:3.8.2

# These two environment variables prevent __pycache__/ files.
ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

# Change the working directory. 
# Every command after this will be run from the /code directory.
WORKDIR /code

# Copy the requirements.txt file.
COPY requirements.txt /code/

# Upgrade pip
RUN pip install --upgrade pip

# Install the requirements.
RUN pip install -r requirements.txt

# Copy the rest of the code. 
COPY . /code/

RUN python3 manage.py migrate && \
 export DJANGO_SUPERUSER_PASSWORD="admin" && \
 python3 manage.py createsuperuser --username admin --email admin@localhost --skip-checks --no-input
FROM python:3.6.4

RUN mkdir /rulermanager

copy rulermanager /rulermanager

WORKDIR /rulermanager

CMD ["python","main.py"]
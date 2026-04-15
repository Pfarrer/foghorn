FROM mcr.microsoft.com/playwright/python:v1.52.0-noble

RUN pip install --no-cache-dir playwright==1.52.0

WORKDIR /app

COPY check.py /app/check.py

CMD ["python", "/app/check.py"]

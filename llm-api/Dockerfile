FROM python:3.11-slim

ENV PROVIDER=""
ENV BASE_URL=""
ENV MODEL=""
ENV DB_HOST=""
ENV DB_PORT=""
ENV DB_USER=""
ENV DB_PASS=""
ENV DB_NAME=""

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY *.py .
EXPOSE 80
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]

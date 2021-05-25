FROM python:3.8-slim

WORKDIR /usr/src/app

COPY requirement.txt ./

RUN pip install --no-cache-dir -r requirement.txt

COPY . .

CMD ["python", "update_floating_ip.py"]
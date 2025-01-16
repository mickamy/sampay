FROM node:22-alpine3.19
WORKDIR /app
RUN npm install -g aws-ses-v2-local
ENTRYPOINT ["aws-ses-v2-local", "--host=0.0.0.0"]

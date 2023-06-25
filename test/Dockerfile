FROM node:lts-alpine3.17
RUN npm install -g newman
WORKDIR /app
CMD ["newman", "run", "postman.json"]
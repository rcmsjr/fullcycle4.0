FROM node:12.18.3-alpine3.12

COPY src ./
RUN npm install --production

ENV NODE_PORT=8080
EXPOSE 8080

CMD npm run start

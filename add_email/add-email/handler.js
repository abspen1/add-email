"use strict";
const fs = require("fs");
const Redis = require("ioredis");
const validator = require("email-validator");

module.exports = async (event, context) => {
  const pass = fs.readFileSync("/var/openfaas/secrets/redis-password", "utf8");

  if (!validator.validate(event.body.info.email)) {
    throw new Error("Email not valid");
  }

  const redis = new Redis({ post: 6379, host: "192.168.1.6", password: pass });
  var reply = await redis.hset(
    "emails",
    event.body.info.email,
    event.body.info.name
  );

  context.succeed(true);
};

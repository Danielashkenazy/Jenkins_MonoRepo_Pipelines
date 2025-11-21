const request = require("supertest");
const { app } = require("./index");

test("valid user", async () => {
  const res = await request(app)
    .post("/users")
    .send({ name: "Dan", email: "dan@example.com" });

  expect(res.status).toBe(201);
});

test("invalid user", async () => {
  const res = await request(app)
    .post("/users")
    .send({ name: "d", email: "bad" });

  expect(res.status).toBe(400);
});

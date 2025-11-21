const { validateUser } = require("./index");

test("valid user", () => {
  expect(validateUser({ name: "Dan", email: "dan@example.com" })).toBe(true);
});

test("invalid email", () => {
  expect(validateUser({ name: "Dan", email: "bad" })).toBe(false);
});

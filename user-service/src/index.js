const express = require("express");
const app = express();
app.use(express.json());

function validateUser(data) {
  if (!data.name || data.name.length < 2) return false;
  if (!data.email || !data.email.includes("@")) return false;
  return true;
}

app.post("/users", (req, res) => {
  const user = req.body;

  if (!validateUser(user)) {
    return res.status(400).json({ error: "invalid user" });
  }

  res.status(201).json({ status: "created", user });
});

if (require.main === module) {
  app.listen(3000, () => console.log("ussssssSssser-service running"));
}

module.exports = { app, validateUser };

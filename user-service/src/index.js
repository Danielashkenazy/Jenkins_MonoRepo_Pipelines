const express = require("express");
const app = express();
const port = process.env.PORT || 3000;

app.get("/health", (req, res) => {
  res.json({ status: "ok", service: "user-service" });
});

app.listen(port, () => {
  console.log(`user-service running on port ${port}`);
  console.log('detect changes test22211')
});

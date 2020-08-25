const express = require('express');
const app = express();
const port = process.env.NODE_PORT || 3000;

app.get('/', (req, res) => {
  res.send('Eu sou Full Cycle.');
});

app.listen(port, () => {
  console.log(`App listening at port[${port}]`);
});

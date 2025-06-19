const express = require('express');
const app = express();

app.get('/public', (req, res) => {
  res.send('public');
});

app.get('/protected', (req, res) => {
  res.send('protected');
});

app.listen(8080, () => {
  console.log('Frontend running on port 8080');
});
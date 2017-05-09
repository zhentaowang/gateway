var express = require('express');
var path = require('path');
var app = express();

app.use(express.static(path.join(__dirname,'dist')));

app.set('port', process.env.PORT || 3100);
var server = app.listen(app.get('port'), function() {
    console.log("start in " + app.get('port'));
});

module.exports = app;

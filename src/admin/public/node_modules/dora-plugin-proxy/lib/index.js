'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _doraAnyproxy = require('dora-anyproxy');

var _getRule = require('./getRule');

var _getRule2 = _interopRequireDefault(_getRule);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

if (!(0, _doraAnyproxy.isRootCAFileExists)()) (0, _doraAnyproxy.generateRootCA)();

exports.default = {
  name: 'proxy',
  'server.before': function serverBefore() {
    this.set('__server_listen_log', false);
  },
  'middleware.before': function middlewareBefore() {
    var _this = this;

    return new Promise(function (resolve) {
      var log = _this.log;
      var query = _this.query;

      log.debug('query: ' + JSON.stringify(query));
      var port = query && query.port || 8989;
      var proxyServer = new _doraAnyproxy.proxyServer({
        type: 'http',
        port: port,
        hostname: 'localhost',
        rule: (0, _getRule2.default)(_this),
        autoTrust: true
      });
      proxyServer.on('finish', function (err) {
        if (err) {
          log.error(err);
        } else {
          log.info('listened on ' + port);
        }
        resolve();
      });
    });
  }
};
module.exports = exports['default'];
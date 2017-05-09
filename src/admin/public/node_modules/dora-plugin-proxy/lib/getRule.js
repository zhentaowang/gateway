'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = getRule;

var _path = require('path');

var _fs = require('fs');

var _getProxyConfig = require('./getProxyConfig');

var _getProxyConfig2 = _interopRequireDefault(_getProxyConfig);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function getRule(args) {
  var cwd = args.cwd;
  var port = args.port;
  var log = args.log;

  var _ref = args.query || {};

  var config = _ref.config;


  var userRuleFile = (0, _path.join)(cwd, 'rule.js');
  if ((0, _fs.existsSync)(userRuleFile)) {
    if (log) log.info('load rule from rule.js');
    return require(userRuleFile);
  }

  var getProxyConfig = (0, _getProxyConfig2.default)(config || 'proxy.config.js', args);
  return require('./rule')({
    port: port,
    hostname: '127.0.0.1',
    getProxyConfig: getProxyConfig,
    cwd: cwd,
    log: log
  });
}
module.exports = exports['default'];
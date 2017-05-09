'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = cdeps;

var _crequire = require('crequire');

var _crequire2 = _interopRequireDefault(_crequire);

var _fs = require('fs');

var _path = require('path');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function endWith(f, str) {
  return f.slice(-str.length).toLowerCase() === str.toLowerCase();
}

function isRelative(filepath) {
  return filepath.charAt(0) === '.';
}

/**
 * Test and return the corrent file.
 * @param f {String}
 * @returns {String}
 */
function getFile(f) {
  // 1. end with .js or .css and exists, return file
  // 2. not end with .js
  // 2.1 add .js and exists, return file
  // 2.2 add /index.js and exists, return file
  // null

  if (endWith(f, '.js') && (0, _fs.existsSync)(f)) return f;
  if (!endWith(f, '.js')) {
    if ((0, _fs.existsSync)(f + '.js')) return f + '.js';
    if ((0, _fs.existsSync)(f + '/index.js')) return f + '/index.js';
  }

  return null;
}

/**
 * Get deps of a file.
 * @param f {String} file
 * @returns {Array}
 */
function parseDeps(f) {
  function getPath(o) {
    return o.path;
  }

  var content = (0, _fs.readFileSync)(f, 'utf-8');
  if (endWith(f, '.js')) {
    return (0, _crequire2.default)(content).map(getPath);
  }

  return [];
}

function parse(entry) {
  var f = getFile(entry);
  if (!f) return [];

  var deps = [entry];
  parseDeps(f).forEach(function (dep) {
    if (isRelative(dep)) {
      var nextDep = (0, _path.resolve)((0, _path.dirname)(f), dep);
      deps = deps.concat(parse(nextDep));
    }
  });
  return deps;
}

function cdeps(entry) {
  return parse(entry);
}
module.exports = exports['default'];
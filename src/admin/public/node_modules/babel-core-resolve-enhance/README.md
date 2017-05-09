# babel-core-resolve-enhance

[![NPM version](https://img.shields.io/npm/v/babel-core-resolve-enhance.svg?style=flat)](https://npmjs.org/package/babel-core-resolve-enhance)
[![NPM downloads](http://img.shields.io/npm/dm/babel-core-resolve-enhance.svg?style=flat)](https://npmjs.org/package/babel-core-resolve-enhance)

Since babel 6 only resolve `presets` and `plugins` from current directory, this project enhanced it, and user can add dirnames for resolve. Extracted from [atool-build](https://github.com/ant-tool/atool-build) .

----

## Usage

```javascript
require('babel-core-resolve-enhance')({
  dirname: __dirname
});
```

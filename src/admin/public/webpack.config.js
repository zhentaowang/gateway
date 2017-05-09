var webpack = require('atool-build/lib/webpack');

module.exports = function(webpackConfig, env) {

    webpackConfig.babel.plugins.push('transform-runtime');
    webpackConfig.babel.plugins.push(['import', {
        libraryName: 'antd',
        style: 'css',
    }]);

    if (env === 'development') {
        webpackConfig.devtool = '#eval';
    }
    else {
        console.log('begin build in ' + process.env.NODE_ENV);
    }

    webpackConfig.plugins.some(function(plugin, i) {
        if (plugin instanceof webpack.optimize.CommonsChunkPlugin) {
            webpackConfig.plugins.splice(i, 1, new webpack.optimize.CommonsChunkPlugin({
                names: ["common", "vendor"],
                minChunks: 2
            }));

            return true;
        }
    });

    return webpackConfig;
};

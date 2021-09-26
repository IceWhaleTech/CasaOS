/*
 * @Author: JerryK
 * @Date: 2021-09-22 10:10:10
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-22 15:26:47
 * @Description: 
 * @FilePath: /CasaOS-UI/vue.config.js
 */
const webpack = require('webpack')

module.exports = {
    publicPath: '/ui/',
    runtimeCompiler: true,
    lintOnSave: false,
    productionSourceMap: false,
    pluginOptions: {

    },
    chainWebpack: config => {
        config.plugin('ignore')
            .use(new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/));
    }
}

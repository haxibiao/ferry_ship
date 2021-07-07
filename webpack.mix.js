/*
 * @Author: Bin
 * @Date: 2021-07-07
 * @FilePath: /ferry_ship/webpack.mix.js
 */

let mix = require('laravel-mix');

// 后台登陆页面
mix.js('src/pages/index/index.js', 'static/js/index.js')
	.react()
	.sass('src/scss/index.scss', 'static/css/index.css')
	.version();

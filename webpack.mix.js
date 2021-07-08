/*
 * @Author: Bin
 * @Date: 2021-07-07
 * @FilePath: /ferry_ship/webpack.mix.js
 */

let mix = require('laravel-mix');

// 首页
mix.js('src/pages/index/index.js', 'static/js/index.js').react().sass('src/scss/index.scss', 'static/css/index.css');

// 登陆页面
mix.js('src/pages/login/index.js', 'static/js/login.js').react().sass('src/scss/login.scss', 'static/css/login.css');

// 控制台页面
mix.js('src/pages/admin/index.js', 'static/js/admin.js').react().sass('src/scss/admin.scss', 'static/css/admin.css');

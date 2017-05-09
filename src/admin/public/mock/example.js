'use strict';

module.exports = {

    'GET /api/example': function(req, res) {
        setTimeout(function() {
            res.json({
                success: true,
                data: ['foo', 'bar'],
            });
        }, 500);
    },

    'GET /api/list': function(req, res) {
        setTimeout(function() {
            res.json({
                "listData": [{
                    "title": "传统 Ajax 已死，Fetch 永生",
                    "url": "https://segmentfault.com/a/1190000003810652"
                }, {
                    "title": "React官网",
                    "url": "http://facebook.github.io/react/index.html"
                }, {
                    "title": "react-components",
                    "url": "http://react-components.com/"
                }, {
                    "title": "React Community",
                    "url": "https://github.com/reactjs"
                }, {
                    "title": "React Router 中文文档",
                    "url": "http://react-guide.github.io/react-router-cn/index.html"
                }, {
                    "title": "Redux 中文文档",
                    "url": "http://cn.redux.js.org/index.html"
                }, {
                    "title": "React Hot Loader",
                    "url": "http://gaearon.github.io/react-hot-loader/"
                }, {
                    "title": "深入理解 react-router 路由系统",
                    "url": "https://zhuanlan.zhihu.com/p/20381597"
                }, {
                    "title": "Ant Design of React",
                    "url": "http://ant.design/docs/react/introduce"
                }, {
                    "title": "Vue.js中文网",
                    "url": "http://cn.vuejs.org/"
                }, {
                    "title": "vue-router 中文文档",
                    "url": "http://router.vuejs.org/zh-cn/index.html"
                }, {
                    "title": "vuex 中文文档",
                    "url": "http://vuex.vuejs.org/zh-cn/index.html"
                }, {
                    "title": "awesome-vue",
                    "url": "https://github.com/vuejs/awesome-vue"
                }, {
                    "title": "vue-resource",
                    "url": "https://github.com/vuejs/vue-resource"
                }, {
                    "title": "VueStrap",
                    "url": "http://yuche.github.io/vue-strap/"
                }, {
                    "title": "浅谈Vue.js",
                    "url": "https://segmentfault.com/a/1190000004704498"
                }, {
                    "title": "基于Vue.js的表格分页组件",
                    "url": "https://segmentfault.com/a/1190000005174322"
                }]
            });
        }, 100);
    },
};

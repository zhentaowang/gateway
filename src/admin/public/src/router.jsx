import React from 'react';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';

// 引入主路由
import SiderBar from './routes/SiderBar';

// 引入单个页面（包括嵌套的子页面）
import Service from './routes/Service';
import Form from './routes/Form';
import Mofang from './routes/Mofang';
import Login from './routes/Login';
import Reg from './routes/Reg';

export default (
    <Router history={hashHistory}>
        <Route path="/" component={SiderBar}>
            <IndexRoute component={Service} />
            <Route path="/login" component={Login} />
            <Route path="/reg" component={Reg} />
            <Route path="/mofang" component={Mofang} />
            <Route path="/table" component={Service} />
            <Route path="/form" component={Form} />
        </Route>
    </Router>
)

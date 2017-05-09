import './index.html';
import 'animate.css/animate.min.css';
import './index.less';
import React from 'react';
import {render} from 'react-dom';

import routes from './router';

render(routes, document.getElementById('app'))

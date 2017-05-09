import React from 'react';
import { Link} from 'react-router'
import { Menu, Icon, Switch } from 'antd';
const SubMenu = Menu.SubMenu;

const ACTIVE = { color: 'red' };

class SiderBar extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            current: '',
            username: ''
        }
    }

    handleClick = (e) => {
        this.setState({
            current: e.key
        })
    }

    componentWillMount() {
        if(this.isLogin()){
            console.log('已登录');
        }else{
            console.log('未登录');
        }
    }

    isLogin(){
        return false
    }

    componentDidMount() {
        this.getUser();
    }

    getUser = () => {
        this.setState({
            username: 'renzhao'
        })
    }

    render() {
        return (
            <div>
                <div id="leftMenu">
                    <img src={require('../assets/images/logo.png')} width="50" id="logo"/>
                    <Menu theme="dark"
                        onClick={this.handleClick}
                        style={{ width: 185 }}
                        defaultOpenKeys={['sub1', 'sub2']}
                        defaultSelectedKeys={[this.state.current]}
                        mode="inline"
                    >
                        <SubMenu key="sub1" title={<span><Icon type="mail" /><span>导航一</span></span>}>
                            <Menu.Item key="1"><Link to="/table">服务</Link></Menu.Item>
                            <Menu.Item key="1"><Link to="/table">API</Link></Menu.Item>
                            <Menu.Item key="2"><Link to="/form">表单</Link></Menu.Item>
                        </SubMenu>
                        <SubMenu key="sub2" title={<span><Icon type="appstore" /><span>导航二</span></span>}>
                            <Menu.Item key="7"><Link to="/mofang">魔方</Link></Menu.Item>
                        </SubMenu>
                    </Menu>
                </div>
                <div id="rightWrap">
                    <Menu mode="horizontal">
                        <SubMenu title={<span><Icon type="user" />{ this.state.username }</span>}>
                            <Menu.Item key="setting:1">退出</Menu.Item>
                        </SubMenu>
                        <SubMenu title={<Link to="/reg">注册</Link>}>
                        </SubMenu>
                        <SubMenu title={<Link to="/login">登录</Link>}>
                        </SubMenu>
                    </Menu>
                    <div className="right-box">
                        { this.props.children }
                    </div>
                </div>
            </div>
        )
    }
}
export default SiderBar;

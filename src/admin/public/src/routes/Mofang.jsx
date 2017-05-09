import React from 'react';
import { Card } from 'antd';

export default class Login extends React.Component {
    constructor(props) {
        super(props)
    }
    render() {
        const style={
            width: 300,
            margin: 5
        }
        return (
            <div>
                <Card style={style} bodyStyle={{ padding: 0 }}>
                    <div className="custom-card">
                        <h3>第二层</h3>
                        <p>右手 R U' R' F R' U' R</p>
                        <p>左手 L' U L 归位 </p>
                    </div>
                </Card>
                <Card style={style} bodyStyle={{ padding: 0 }}>
                    <div className="custom-image">
                        <img width="100%" src={require('../assets/images/1.jpg')} />
                    </div>
                    <div className="custom-card">
                        <h3>十字</h3>
                    </div>
                </Card>
                <Card style={style} bodyStyle={{ padding: 0 }}>
                    <div className="custom-image">
                        <img width="100%" src={require('../assets/images/2.jpg')} />
                    </div>
                    <div className="custom-card">
                        <h3>顶层</h3>
                    </div>
                </Card>
                <Card style={style} bodyStyle={{ padding: 0 }}>
                    <div className="custom-image">
                        <img width="100%" src={require('../assets/images/3.jpg')} />
                    </div>
                    <div className="custom-card">
                        <h3>倒数第二步</h3>
                    </div>
                </Card>
                <Card style={style} bodyStyle={{ padding: 0 }}>
                    <div className="custom-image">
                        <img width="100%" src={require('../assets/images/4.jpg')} />
                    </div>
                    <div className="custom-card">
                        <h3>倒数第一步</h3>
                    </div>
                </Card>
            </div>
        )
    }
}

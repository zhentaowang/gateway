import React from 'react';
import {Table, Icon} from 'antd';
import {Server} from '../config';
import fetch from '../utils/request';

export default class Service extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            tDate: [],
            selectedRowKeys: []
        }
    }

    componentWillMount() {
        this.routeService()
    }

    routeService() {
        let url = Server + "services"
        let options = {
            method: 'GET'
        }
        fetch(url, options).then(data => {
            console.log(data)
        })
        // const data = []

        // for (let i = 0; i < 99; i++) {
        //     data.push({
        //         key: i,
        //         name: `Mr劳卜${i}`,
        //         age: 18,
        //         address: `西湖区湖底公园${i}号`,
        //         remark: 'http://www.cnblogs.com/luozhihao/',
        //         operate: '暂无'
        //     })
        // }

        // this.setState({
        //     tDate: data
        // })
    }

    // checkbox状态
    onSelectChange = (selectedRowKeys) => {
        console.log('selectedRowKeys changed: ', selectedRowKeys)
        this.setState({ selectedRowKeys })
    }

    render() {
        const columns = [{
            title: 'id',
            width: '5%',
            dataIndex: 'service_id'
        }, {
            title: 'namespace',
            width: '20%',
            dataIndex: 'namespace',
        }, {
            title: '服务名称',
            width: '20%',
            dataIndex: 'name'
        }, {
            title: '端口',
            width: '10%',
            dataIndex: 'port',
        }, {
            title: '协议',
            width: '20%',
            dataIndex: 'protocol',
            render(text) {
                return <a href={text} target="_blank">博客园</a>
            }
        }, {
            title: '操作',
            width: '20%',
            dataIndex: 'operate'
        }]

        const { selectedRowKeys } = this.state

        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange
        }

        const pagination = {
            total: this.state.tDate.length,
            showSizeChanger: true,
            onShowSizeChange(current, pageSize) {
                console.log('Current: ', current, '; PageSize: ', pageSize)
            },
            onChange(current) {
                console.log('Current: ', current)
            }
        }

        return (
            <Table rowSelection={rowSelection} columns={columns} dataSource={this.state.tDate} bordered pagination={pagination} className="animated zoomIn"/>
        )
    }
}

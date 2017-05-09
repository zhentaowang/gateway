import React from 'react';
import { Form, Input, Select, Checkbox, DatePicker, Col, Radio, Button, Modal, message } from 'antd';

const FormItem = Form.Item;
const Option = Select.Option;
const RadioGroup = Radio.Group;
const CheckboxGroup = Checkbox.Group;

class myForm extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            visible: false
        }
    }

    // 选择select
    handleSelectChange = (value) => {
        console.log(`selected ${value}`)
    }

    // 提交表单
    handleSubmit = (e) => {
        e.preventDefault();
        this.props.form.validateFields((err, values) => {
            if (!err) {
                message.success('操作成功!');
                console.log('收到表单值: ', values);
                this.props.form.resetFields();
            }
        });
    }


    // 显示弹框
    showModal = () => {
        this.setState({ visible: true })
    }


    // 隐藏弹框
    hideModal = () => {
        this.setState({ visible: false })
    }

    render() {
        const { getFieldDecorator } = this.props.form

        const formItemLayout = {
            labelCol: { span: 3 },
            wrapperCol: { span: 10 }
        }

        return (
            <Form horizontal onSubmit={this.handleSubmit} className="animated zoomIn">
                <FormItem label="用户名" {...formItemLayout} hasFeedback>
                    {getFieldDecorator('username', {
                        rules: [{ required: true, message: '请输入用户名!' }],
                    })(
                        <Input placeholder="Username" />
                    )}
                </FormItem>

                <FormItem label="密码" {...formItemLayout} hasFeedback>
                    {getFieldDecorator('password', {
                        rules: [{ required: true, message: '请输入密码!' }],
                    })(
                        <Input type="password" placeholder="password" />
                    )}
                </FormItem>

                <FormItem label="日期选择框" labelCol={{ span: 3 }} required>
                    <Col span="4">
                        <FormItem>
                            {getFieldDecorator('startDate', {
                                rules: [{ required: true, message: '请输入起始日期!' }],
                            })(
                                <DatePicker />
                            )}
                        </FormItem>
                    </Col>
                    <Col span="1">
                        <p className="ant-form-split">-</p>
                    </Col>
                    <Col span="4">
                        <FormItem>
                            {getFieldDecorator('endDate', {
                                rules: [{ required: true, message: '请输入结束日期!' }],
                            })(
                                <DatePicker />
                            )}
                        </FormItem>
                    </Col>
                </FormItem>

                <FormItem id="control-textarea" label="简介" {...formItemLayout}>
                    {getFieldDecorator('content', {
                        rules: [{ required: true, message: '请输入简介!' }],
                    })(
                        <Input type="textarea" id="control-textarea" rows="3" />
                    )}
                </FormItem>

                <FormItem label="Select 选择器" {...formItemLayout} required>
                    {getFieldDecorator('people', {
                        initialValue: 'lucy'
                    })(
                        <Select size="large" style={{ width: 200 }} onChange={this.handleSelectChange}>
                            <Option value="jack">jack</Option>
                            <Option value="lucy">lucy</Option>
                            <Option value="disabled" disabled>disabled</Option>
                            <Option value="yiminghe">yiminghe</Option>
                        </Select>
                    )}
                </FormItem>

                <FormItem label="Checkbox 多选框" {...formItemLayout}>
                    {getFieldDecorator('checkboxItem')(
                        <CheckboxGroup options={
                            [
                                {
                                    label: 'Apple',
                                    value: 'Apple'
                                }, {
                                    label: 'Pear',
                                    value: 'Pear'
                                }, {
                                    label: 'Orange',
                                    value: 'Orange'
                                }
                            ]
                            }
                        />
                    )}
                </FormItem>

                <FormItem id="select" label="Radio 单选框" {...formItemLayout}>
                    {getFieldDecorator('radioItem')(
                        <RadioGroup>
                            <Radio value="a">A</Radio>
                            <Radio value="b">B</Radio>
                            <Radio value="c">C</Radio>
                            <Radio value="d">D</Radio>
                        </RadioGroup>
                    )}
                </FormItem>

                <FormItem wrapperCol={{ span: 6, offset: 3 }} style={{ marginTop: 24 }}>
                    <Button type="primary" htmlType="submit">确定</Button>
                    &nbsp;&nbsp;&nbsp;
                    <Button type="ghost" onClick={this.showModal}>点击有惊喜</Button>
                </FormItem>

                <Modal title="登录" visible={this.state.visible} onOk={this.hideModal} onCancel={this.hideModal}>
                    这是一个modal
                </Modal>
            </Form>
        )
    }
}

myForm = Form.create()(myForm);

export default myForm;

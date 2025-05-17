import {PlusOutlined} from '@ant-design/icons';
import {
  ModalForm,
  PageContainer,
  ProFormDigit,
  ProFormText,
  ProFormTextArea,
  ProTable
} from '@ant-design/pro-components';
import {Button, Divider, Form, message, Popconfirm} from 'antd';
import React, {useRef, useState} from 'react';
import {requestGet, requestPost} from "@/utils/requestTool";

const Index: React.FC = () => {
  const [modalOpen, setmodalOpen] = useState(false)
  const [modelTitle, setModelTitle] = useState('新增授权码')
  const [form] = Form.useForm();
  const actionRef = useRef();
  const [loading, setLoading] = useState(false);
  const columns = [
    {
      title: 'id',
      dataIndex: 'id',
      search: false,
    },
    {
      title: '授权码',
      dataIndex: 'authKey',
    },
    {
      title: '描述',
      dataIndex: 'remark',
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      search: false,
    },
    {
      title: '操作',
      search: false,
      render: record => {
        return <>
          <a onClick={() => {
            setModelTitle('编辑授权码')
            setmodalOpen(true)
            let data = {...record}
            form.setFieldsValue(data)
          }}>编辑</a>
          <Divider type="vertical"/>
          <Popconfirm
            title='您确定要删除吗？'
            description='删除后，数据将无法恢复，请慎重！'
            onConfirm={e => {
              requestPost('/api/delAuthKey', {ids: [record.id]}).then(res => {
                if (res.code === 200) {
                  message.success(res.message)
                  actionRef.current.reload()
                }
              })
            }}
          >
            <a style={{color: 'red',}}>删除</a>
          </Popconfirm>
        </>
      }
    },
  ];
  const renderToolBar = () => {
    return [
      <Button key="add" type="primary" icon={<PlusOutlined/>} onClick={() => {
        form.resetFields()
        setModelTitle('新增授权码')
        setmodalOpen(true)
      }}>
        添加
      </Button>,
    ];
  };
  const getData = async (params: {}) => {
    let value = {
      data: [],
      success: true,
      total: 0,
    }
    params.page = params.current || 1
    params.pageSize = params.pageSize || 10
    await requestGet('/api/getAuthKeyList', params).then(res => {
      value.success = res.code === 200 ? true : false
      value.data = res.data.list || []
      value.data = value.data.map(item => {
        item.key = item.Id
        return item
      })
      value.total = res.data.total || 0
    })
    return Promise.resolve(value)
  }
  return (
    <PageContainer style={{minHeight: window.innerHeight - 150}}>
      <ProTable
        rowKey='id'
        columns={columns}
        search={{
          labelWidth: 120,
        }}
        request={getData}
        toolBarRender={renderToolBar}
        actionRef={actionRef}
      />
      <ModalForm
        width={500}
        title={modelTitle}
        open={modalOpen}
        form={form}
        loading={loading}
        modalProps={{
          onCancel: () => setmodalOpen(false)
        }}
        onFinish={values => {
          setLoading(true)
          // @ts-ignore
          requestPost('/api/saveAuthKey', values).then(res => {
            if (res.code === 200) {
              message.success(res.message)
              setmodalOpen(false)
              form.resetFields()
              actionRef.current && actionRef.current.reload()
            }
            setLoading(false)
          })
        }}
      >
        <ProFormDigit
          name="id"
          hidden={true}
        />
        <ProFormText
          name="authKey"
          label="授权码"
          placeholder="自动生成"
          disabled
        />
        <ProFormTextArea
          rules={[
            {
              required: true,
              message: "请输入",
            },
          ]}
          name="remark"
          label="描述"
          placeholder="请输入"
        />
      </ModalForm>
    </PageContainer>
  );
};

export default Index;

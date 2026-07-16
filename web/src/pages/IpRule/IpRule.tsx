import {DeleteOutlined, PlusOutlined} from '@ant-design/icons';
import {
  ModalForm,
  PageContainer,
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable
} from '@ant-design/pro-components';
import {Button, Form, message, Popconfirm, Switch, Tag} from 'antd';
import React, {useRef, useState} from 'react';
import {requestGet, requestPost} from "@/utils/requestTool";
import {ProFormDateTimePicker} from "@ant-design/pro-form";
import moment from "moment";

const Index: React.FC = () => {
  const [modalOpen, setmodalOpen] = useState(false)
  const [modelTitle, setModelTitle] = useState('新增封禁IP')
  const [isEdit, setIsEdit] = useState(false)
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [form] = Form.useForm();
  const actionRef = useRef();
  const [loading, setLoading] = useState(false);

  const columns = [
    {
      title: 'ID标识',
      dataIndex: 'id',
      search: false,
    },
    {
      title: 'IP',
      dataIndex: 'ip',
      render: (ip, record) => {
        return <div>
          <div>{ip}</div>
          <Tag>{record.ip_location}</Tag>
        </div>
      }
    },
    {
      title: '端口',
      dataIndex: 'port',
      search: false,
    },
    {
      title: '协议',
      dataIndex: 'protocol',
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'status',
      search: true,
      valueType: 'select',
      valueEnum: {
        1: {text: '启用'},
        0: {text: '禁用'}
      },
      render: (_, record) => {
        return (
          <Popconfirm
            title="确认修改状态？"
            onConfirm={() => {
              requestPost('/api/changeStatus', {
                id: record.id,
                status: record.status === 1 ? 0 : 1
              }).then(res => {
                if (res.code === 200) {
                  message.success(res.message);
                  actionRef.current?.reload();
                }
              });
            }}
          >
            <Switch
              checked={record.status === 1}
              checkedChildren="启用"
              unCheckedChildren="禁用"
            />
          </Popconfirm>
        );
      }
    },
    {
      title: '封禁原因',
      dataIndex: 'reason',
      search: true,
    },
    {
      title: '过期时间',
      dataIndex: 'expired_at',
      search: false,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      search: false,
    },
    {
      title: '操作',
      search: false,
      render: record => {
        return <>
          <a onClick={() => {
            form.resetFields()
            form.setFieldsValue({
              id: record.id,
              ips: record.ip,
              protocol: record.protocol,
              port: record.port,
              reason: record.reason,
              expiredAt: record.expired_at ? moment(record.expired_at, 'YYYY-MM-DD HH:mm:ss') : undefined,
            })
            setIsEdit(true)
            setModelTitle("编辑封禁IP")
            setmodalOpen(true)
          }}>编辑</a>
          <Popconfirm
            title='您确定要删除吗？'
            description='删除后，数据将无法恢复，请慎重！'
            onConfirm={e => {
              requestPost('/api/unBanIp', {ids: [record.id]}).then(res => {
                if (res.code === 200) {
                  message.success(res.message)
                  actionRef.current.reload()
                }
              })
            }}
          >
            <a style={{color: 'red', marginLeft: 16}}>删除</a>
          </Popconfirm>
        </>
      }
    },
  ];

  const handleBatchDelete = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的记录');
      return;
    }
    requestPost('/api/unBanIp', {ids: selectedRowKeys}).then(res => {
      if (res.code === 200) {
        message.success(res.message);
        setSelectedRowKeys([]);
        actionRef.current?.reload();
      }
    });
  };

  const renderToolBar = () => {
    return [
      <Button key="add" type="primary" icon={<PlusOutlined/>} onClick={() => {
        form.resetFields()
        form.setFieldsValue({
          currentId: 0,
          step: 100,
          minLength: 0,
        })
        setIsEdit(false)
        setModelTitle("新增封禁IP")
        setmodalOpen(true)
      }}>
        添加
      </Button>,
      <Popconfirm
        key="batchDelete"
        title="确认删除选中记录？"
        description="删除后，数据将无法恢复，请慎重！"
        onConfirm={handleBatchDelete}
      >
        <Button
          type="primary"
          danger
          icon={<DeleteOutlined/>}
          disabled={selectedRowKeys.length === 0}
        >
          批量删除
        </Button>
      </Popconfirm>
    ];
  };
  const getData = async (params: any) => {
    let value = {
      data: [],
      success: true,
      total: 0,
    }
    // 构建后端需要的参数，不直接修改ProTable传入的params（避免污染其受控状态）
    const {current, pageSize, ...rest} = params
    const query = {
      ...rest,
      page: current || 1,
      page_size: Math.min(pageSize || 10, 200),
    }
    await requestGet('/api/getBanIpList', query).then(res => {
      value.success = res.code === 200 ? true : false
      value.data = res.data.list || []
      value.total = res.data.total || 0
    })
    return Promise.resolve(value)
  }

  const rowSelection = {
    selectedRowKeys,
    onChange: (newSelectedRowKeys: React.Key[]) => {
      setSelectedRowKeys(newSelectedRowKeys);
    },
  };

  return (
    <PageContainer style={{minHeight: window.innerHeight - 150}}>
      <ProTable
        rowKey='id'
        columns={columns}
        rowSelection={rowSelection}
        search={{
          labelWidth: 120,
        }}
        request={getData}
        toolBarRender={renderToolBar}
        actionRef={actionRef}
        pagination={{
          defaultPageSize: 10,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100', '200'],
        }}
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
          const url = isEdit ? '/api/editIpRule' : '/api/banIp'
          // @ts-ignore
          requestPost(url, values).then(res => {
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
        <ProFormText
          name="id"
          hidden
        />
        {isEdit ? (
          <ProFormText
            name="ips"
            label="封禁IP"
            placeholder="封禁IP"
            disabled
            readonly
            transform={undefined}
          />
        ) : (
          <ProFormTextArea
            rules={[
              {
                required: true,
                message: "请输入IP地址",
              },
            ]}
            name="ips"
            label="封禁IP"
            placeholder="请输入IP地址，多个IP请用逗号或换行符分隔"
            fieldProps={{
              autoSize: {minRows: 2, maxRows: 6}
            }}
            transform={(value) => ({
              ips: value.split(/[,\n]/).map(ip => ip.trim()).filter(ip => ip).join(',')
            })}
          />
        )}
        <ProFormTextArea
          name="reason"
          label="封禁原因"
          placeholder="请输入封禁原因"
          fieldProps={{
            autoSize: {minRows: 1, maxRows: 4}
          }}
        />
        <ProFormSelect
          name="protocol"
          // mode="tags"
          label="封禁协议"
          placeholder="请输入"
          options={[
            {value: 'tcp', label: 'tcp'},
            {value: 'udp', label: 'udp'},
            {value: 'icmp', label: 'icmp'},
          ]}
          tooltip="不设置将封禁所有协议"
        />
        <ProFormDigit
          rules={[
            {
              type: 'number',
              min: 0,
            },
          ]}
          name="port"
          label="封禁端口"
          placeholder="请输入"
          tooltip="不设置将封禁所有端口"
        />
        <ProFormDateTimePicker
          name="expiredAt"
          label="过期时间"
          placeholder="请选择过期时间"
          tooltip="不设置则永久有效"
          fieldProps={{
            format: 'YYYY-MM-DD HH:mm:ss',
            showTime: {defaultValue: moment('00:00:00', 'HH:mm:ss')},
            presets: [
              {label: '5分钟', value: moment().add(5, 'minutes')},
              {label: '半小时', value: moment().add(30, 'minutes')},
              {label: '2小时', value: moment().add(2, 'hours')},
              {label: '1天', value: moment().add(1, 'days')},
              {label: '1周', value: moment().add(7, 'days')},
              {label: '1个月', value: moment().add(1, 'months')}
            ]
          }}
        />
      </ModalForm>
    </PageContainer>
  );
};

export default Index;

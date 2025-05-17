import {LockOutlined, LogoutOutlined, SettingOutlined, UserOutlined} from '@ant-design/icons';
import {useEmotionCss} from '@ant-design/use-emotion-css';
import {useModel} from '@umijs/max';
import {Form, message, Modal, Spin} from 'antd';
import type {MenuInfo} from 'rc-menu/lib/interface';
import React, {useCallback} from 'react';
import {flushSync} from 'react-dom';
import HeaderDropdown from '../HeaderDropdown';
import {loginOut} from "@/utils/authority";
import Password from "antd/es/input/Password";
import {requestPost} from "@/utils/requestTool";

export type GlobalHeaderRightProps = {
  menu?: boolean;
  children?: React.ReactNode;
};

export const AvatarName = () => {
  const {initialState} = useModel('@@initialState');
  const {currentUser} = initialState || {};
  return <span className="anticon">{currentUser?.name}</span>;
};

export const AvatarDropdown: React.FC<GlobalHeaderRightProps> = ({menu, children}) => {
  /**
   * 退出登录，并且将当前的 url 保存
   */
  const actionClassName = useEmotionCss(({token}) => {
    return {
      display: 'flex',
      height: '48px',
      marginLeft: 'auto',
      overflow: 'hidden',
      alignItems: 'center',
      padding: '0 8px',
      cursor: 'pointer',
      borderRadius: token.borderRadius,
      '&:hover': {
        backgroundColor: token.colorBgTextHover,
      },
    };
  });
  const {initialState, setInitialState} = useModel('@@initialState');
  let passwordData = {}

  const onMenuClick = useCallback(
    (event: MenuInfo) => {
      const {key} = event;
      switch (key) {
        case 'logout':
          flushSync(() => {
            setInitialState((s) => ({...s, currentUser: undefined}));
          });
          loginOut();
          return;
          break;
        case 'changPassword':
          passwordData = {}
          Modal.confirm({
            icon: null,
            title: "修改密码",
            content: <div>
              <Form
                labelCol={{span: 6}}
                wrapperCol={{span: 18}}
              >
                <Form.Item label="新密码" required>
                  <Password
                    onChange={event => {
                      passwordData.password = event.target.value
                    }}
                  />
                </Form.Item>
                <Form.Item label="确认密码" required>
                  <Password
                    onChange={event => {
                      passwordData.confirmPassword = event.target.value
                    }}
                  />
                </Form.Item>
              </Form>
            </div>,
            onOk: e => {
              requestPost('/api/resetPassword', passwordData).then(res => {
                if (res.code === 200) {
                  message.success(res.message)
                  e()
                  //退出登录
                  loginOut()
                }
              })
            },
            onCancel: () => {
              passwordData = {}
            }
          })
          break;
      }
      // history.push(`/account/${key}`);
    },
    [setInitialState],
  );

  const loading = (
    <span className={actionClassName}>
      <Spin
        size="small"
        style={{
          marginLeft: 8,
          marginRight: 8,
        }}
      />
    </span>
  );

  if (!initialState) {
    return loading;
  }

  const {currentUser} = initialState;

  if (!currentUser || !currentUser.name) {
    return loading;
  }

  const menuItems = [
    ...(menu
      ? [
        {
          key: 'center',
          icon: <UserOutlined/>,
          label: '个人中心',
        },
        {
          key: 'settings',
          icon: <SettingOutlined/>,
          label: '个人设置',
        },
        {
          type: 'divider' as const,
        },
      ]
      : []),
    {
      key: 'changPassword',
      icon: <LockOutlined/>,
      label: '修改密码',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined/>,
      label: '退出登录',
    },
  ];

  return (
    <HeaderDropdown
      menu={{
        selectedKeys: [],
        onClick: onMenuClick,
        items: menuItems,
      }}
    >
      {children}
    </HeaderDropdown>
  );
};

import Footer from '@/components/Footer';
import {LockOutlined, UserOutlined,} from '@ant-design/icons';
import {LoginForm, ProFormText,} from '@ant-design/pro-components';
import {useEmotionCss} from '@ant-design/use-emotion-css';
import {Helmet, SelectLang, useIntl} from '@umijs/max';
import {Alert, message} from 'antd';
import Settings from '../../../../config/defaultSettings';
import defaultSettings from '../../../../config/defaultSettings';
import React, {useState} from 'react';
import {requestPost} from "@/utils/requestTool";
import {setAccessToken} from "@/utils/authority";
import {getFullPath} from "@/utils/utils";

const Lang = () => {
  const langClassName = useEmotionCss(({token}) => {
    return {
      width: 42,
      height: 42,
      lineHeight: '42px',
      position: 'fixed',
      right: 16,
      borderRadius: token.borderRadius,
      ':hover': {
        backgroundColor: token.colorBgTextHover,
      },
    };
  });

  return (
    <div className={langClassName} data-lang>
      {SelectLang && <SelectLang/>}
    </div>
  );
};

const LoginMessage: React.FC<{
  content: string;
}> = ({content}) => {
  return (
    <Alert
      style={{
        marginBottom: 24,
      }}
      message={content}
      type="error"
      showIcon
    />
  );
};

const Login: React.FC = () => {
  const [userLoginState, setUserLoginState] = useState<API.LoginResult>({});
  const [type, setType] = useState<string>('account');

  const containerClassName = useEmotionCss(() => {
    return {
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      overflow: 'auto',
      backgroundImage:
        "url('https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/V-_oS6r-i7wAAAAAAAAAAAAAFl94AQBr')",
      backgroundSize: '100% 100%',
    };
  });

  const intl = useIntl();

  const {status, type: loginType} = userLoginState;

  return (
    <div className={containerClassName}>
      <Helmet><title>{intl.formatMessage({id: 'menu.login'})}-{Settings.title}</title></Helmet>
      <Lang/>
      <div
        style={{
          flex: '1',
          padding: '32px 0',
        }}
      >
        <LoginForm
          contentStyle={{
            minWidth: 280,
            maxWidth: '75vw',
          }}
          logo={<img alt="logo" src={getFullPath("/favicon.ico")}/>}
          title={defaultSettings.title}
          subTitle="可视化ip防火墙管理系统，在线管理你的服务器安全"
          initialValues={{
            autoLogin: true,
          }}
          onFinish={async (values) => {
            requestPost('/api/login', values).then(res => {
              if (res.code === 200) {
                message.success(res.message);
                setAccessToken(res.data.token)
                const urlParams = new URL(window.location.href).searchParams;
                window.location.href = urlParams.get('redirect') || '/'
              }
            })
          }}
        >

          {status === 'error' && loginType === 'account' && (
            <LoginMessage
              content={intl.formatMessage({
                id: 'pages.login.accountLogin.errorMessage',
                defaultMessage: '账户或密码错误(admin/ant.design)',
              })}
            />
          )}
          {type === 'account' && (
            <>
              <ProFormText
                name="name"
                fieldProps={{
                  size: 'large',
                  prefix: <UserOutlined/>,
                }}
                placeholder="请输入账号"
                rules={[
                  {
                    required: true,
                    message: "请输入账号",
                  },
                ]}
              />
              <ProFormText.Password
                name="password"
                fieldProps={{
                  size: 'large',
                  prefix: <LockOutlined/>,
                }}
                placeholder="请输入密码"
                rules={[
                  {
                    required: true,
                    message: "请输入密码",
                  },
                ]}
              />
            </>
          )}

        </LoginForm>
      </div>
      <Footer/>
    </div>
  );
};

export default Login;

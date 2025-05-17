// use localStorage to store the authority info, which might be sent from server in actual project.
import {getStorage, removeStorage, setStorage} from '@/utils/utils';
import {extraConfig} from "../../config/extraConfig";
import {requestGet} from "@/utils/requestTool";

export function getAuthority(str: any) {
  // return localStorage.getItem('antd-pro-authority') || ['admin', 'user'];
  const authorityString =
    typeof str === 'undefined' ? localStorage.getItem('antd-pro-authority') : str;
  // authorityString could be admin, "admin", ["admin"]
  let authority;
  try {
    authority = JSON.parse(authorityString);
  } catch (e) {
    authority = authorityString;
  }
  if (typeof authority === 'string') {
    return [authority];
  }
  // preview.pro.ant.design only do not use in your production ; preview.pro.ant.design 专用环境变量，请不要在你的项目中使用它。
  if (!authority && ANT_DESIGN_PRO_ONLY_DO_NOT_USE_IN_YOUR_PRODUCTION === 'site') {
    return ['admin'];
  }
  return authority;
}

export function setAuthority(authority: string[]) {
  const proAuthority = typeof authority === 'string' ? [authority] : authority;
  return localStorage.setItem('antd-pro-authority', JSON.stringify(proAuthority));
}

export function delAuthority() {
  return localStorage.removeItem('antd-pro-authority');
}

export function getAccessToken() {
  return localStorage.getItem('access_token');
}

export function setAccessToken(access_token: string) {
  return localStorage.setItem('access_token', access_token);
}

export function delAccessToken() {
  return localStorage.removeItem('access_token');
}

export function setUserInfo(info: any, timeout = 7200) {
  return setStorage('user_info', info, timeout);
}

export function getUserInfo() {
  return getStorage('user_info');
}

export function removeUserInfo() {
  return removeStorage('user_info');
}

// 退出登录
export function loginOut() {
  const token = getAccessToken()
  delAuthority();
  delAccessToken();
  removeUserInfo();
  // const redirect = getQueryString('redirect');
  // let url;
  // if (window.location.pathname !== extraConfig.url401 && !redirect) {
  //   url = `/api/logout?authorization=${token}&redirect=` + encodeURI(extraConfig.url401 + '?redirect=' + window.location.href);
  // } else {
  //   url = `/api/logout?authorization=${token}&redirect=` + encodeURI(extraConfig.url401 + '?redirect=' + redirect);
  // }
  window.location.href = extraConfig.url401;
}

/**
 * 网站登录时，初始化登录账号信息
 * @param expires_in
 */
export function initUserInfo(expires_in: any) {
  requestGet('/api/auth/getUserInfo').then(json => {
    setUserInfo(json.info, expires_in || 7200);
    // const {redirect} = getPageQuery();
    let url = '/';
    if (redirect) {
      url = redirect;
    }
    if (json.info.is_admin === 1) {
      setAuthority(['super_admin']);
    } else {
      window.location.href = `/exception/403?redirect=${window.location.origin}/`;
      return;
    }

    window.location.href = url;
  });
}



/**
 * 设置存储
 * @param key
 * @param value
 * @param expires //超时时间，单位秒
 * @returns {boolean}
 */
import {extraConfig} from "../../config/extraConfig";

export function setStorage(key: string, value: any, expires = null) {
  var temp = {
    value: value,
  };
  if (expires !== null) {
    temp.expires = expires + parseInt(new Date().getTime() / 1000);
  }

  localStorage.setItem(key, JSON.stringify(temp));
  return true;
}

/**
 * 获取存储
 * @param key
 * @returns {string | number|boolean}
 */
export function getStorage(key: string) {
  const time = parseInt(new Date().getTime() / 1000);
  var res = localStorage.getItem(key);
  res = JSON.parse(res);
  if (res == null || !res.hasOwnProperty('value')) {
    return false;
  }
  if (!res.hasOwnProperty('expires')) {
    return res.value;
  }
  if (res.expires >= time) {
    return res.value;
  }
  localStorage.removeItem(key);
  return false;
}

/**
 * 移除缓存
 * @param key
 */
export function removeStorage(key: string) {
  localStorage.removeItem(key);
}

/**
 * 获取url参数
 * @param name
 * @returns {string|null}
 */
export function getQueryString(name: string) {
  var reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)', 'i');
  var r = window.location.search.substr(1).match(reg);
  if (r != null) return unescape(r[2]);
  return undefined;
}

/**
 * 获取带前缀的绝对地址
 * @param path
 * @returns {*}
 */
export function getFullPath(path: string) {
  const prefix = extraConfig.routePrefix ? extraConfig.routePrefix.replace(/\/$/, '') : ''
  return prefix + path;
}



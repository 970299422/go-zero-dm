import request from '@/utils/request'
import type { ResponseData } from '@/utils/request'

// 登录请求参数
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应数据
export interface LoginResponse {
  accessToken: string
  accessExpire: number
}

// 用户信息（对接后端接口）
export interface UserInfo {
  id: number
  username: string
  email?: string
}

/**
 * 用户登录
 * @param data - 登录数据
 * @returns Promise
 */
export function login(data: LoginRequest): Promise<ResponseData<LoginResponse>> {
  return request({
    url: '/api/user/login',
    method: 'post',
    data
  })
}

/**
 * 用户注册
 * @param data - 注册数据
 * @returns Promise
 */
export function register(data: {
  username: string
  password: string
}): Promise<ResponseData<LoginResponse>> {
  return request({
    url: '/api/user/register',
    method: 'post',
    data
  })
}

/**
 * 获取当前用户信息
 * @returns Promise
 */
export function getUserInfo(): Promise<ResponseData<UserInfo>> {
  return request({
    url: '/api/user/userinfo',
    method: 'get'
  })
}


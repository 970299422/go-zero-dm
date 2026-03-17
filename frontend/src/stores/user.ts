import { defineStore } from 'pinia'
import { getToken, setToken, removeToken, getUserInfo, setUserInfo, removeUserInfo } from '@/utils/auth'
import { login, getUserInfo as getUserInfoAPI, type LoginRequest, type UserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken() || '',
    userInfo: getUserInfo() as UserInfo | null
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.userInfo?.username || '',
    email: (state) => state.userInfo?.email || ''
  },

  actions: {
    // 登录
    async login(loginForm: LoginRequest) {
      try {
        const res = await login(loginForm)
        const { accessToken } = res.data
        this.token = accessToken
        setToken(accessToken)

        // 登录成功后再拉取一次用户信息
        const infoRes = await getUserInfoAPI()
        this.userInfo = infoRes.data
        setUserInfo(infoRes.data)
        return Promise.resolve(res)
      } catch (error) {
        return Promise.reject(error)
      }
    },

    // 获取用户信息
    async getUserInfo() {
      try {
        const res = await getUserInfoAPI()
        this.userInfo = res.data
        setUserInfo(res.data)
        return Promise.resolve(res)
      } catch (error) {
        return Promise.reject(error)
      }
    },

    // 登出
    async logout() {
      try {
        this.token = ''
        this.userInfo = null
        removeToken()
        removeUserInfo()
        return Promise.resolve()
      } catch (error) {
        return Promise.reject(error)
      }
    }
  }
})


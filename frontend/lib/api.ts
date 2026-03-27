import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const apiClient = {
  get: <T>(url: string) => api.get<T>(url).then((res) => res.data),
  post: <T>(url: string, data?: unknown) => api.post<T>(url, data).then((res) => res.data),
  patch: <T>(url: string, data?: unknown) => api.patch<T>(url, data).then((res) => res.data),
  delete: <T>(url: string) => api.delete<T>(url).then((res) => res.data),
}
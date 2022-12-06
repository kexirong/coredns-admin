import axios from 'axios'
import { AxiosResponse } from 'axios'

export interface DomainsData {
    name: string
    subdomain?: DomainsData[]
    key?: string
}

export interface RecordData {
    type: string
    ttl: number
    priority: number
    name: string
    content: string
    key?: string
    action?: 'add' | 'edit'
    fingerprint?: string
}

export function getDomains(deep: number) {
    return axios.get<DomainsData>('/api/v1/redis/domains', { params: { deep } })
}

export function getRecords() {
    return axios.get<RecordData[]>('/api/v1/redis/records')
}

export function getRecord(key: string) {
    return axios.get<RecordData[]>(`/api/v1/redis/record/${key}`)
}

export function putRecord(key: string, fingerprint: string, data: RecordData) {
    return axios.put(`/api/v1/redis/record/${key}`, data, { headers: { fingerprint } })
}

export function postRecordSignature(data: RecordData) {
    return axios.post<{ fingerprint: string }>('/api/v1/redis/record/signature', data)
}

export function postRecord(data: RecordData) {
    return axios.post<{ key: string }>('/api/v1/redis/record', data)
}

export function deleteRecord(key: string, fingerprint: string) {
    return axios.delete(`/api/v1/redis/record/${key}`, { headers: { fingerprint } })
}
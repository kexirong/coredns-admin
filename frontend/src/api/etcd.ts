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
}

export function getDomains() {
    return axios.get<DomainsData>('/api/v1/etcd/domains')
}

export function getRecords() {
    return axios.get<RecordData[]>('/api/v1/etcd/records')
}

export function getRecord(path: string) {
    return axios.get<RecordData[]>(`/api/v1/etcd/record/${path}`)
}

export function putRecord(path: string, data: RecordData) {
    return axios.put(`/api/v1/etcd/record/${path}`, data)
}

export function postRecord(data: RecordData) {
    return axios.post<{ key: string }>('/api/v1/etcd/record', data)
}

export function deleteRecord(path: string) {
    return axios.delete(`/api/v1/etcd/record/${path}`)
}
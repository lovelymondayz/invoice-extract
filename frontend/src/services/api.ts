import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface Invoice {
  id: number
  filename: string
  vendor?: string
  invoice_number?: string
  invoice_date?: string
  total_amount?: number
  currency?: string
  status: string
  created_at: string
}

export interface LineItem {
  id: number
  description: string
  quantity: number
  unit_price: number
  total: number
}

export interface InvoiceDetail extends Invoice {
  line_items: LineItem[]
  raw_text?: string
}

export async function getInvoices(): Promise<Invoice[]> {
  const res = await api.get('/invoices')
  return res.data
}

export async function getInvoice(id: number): Promise<InvoiceDetail> {
  const res = await api.get(`/invoices/${id}`)
  return res.data
}

export async function extractInvoice(file: File): Promise<any> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await api.post('/extract', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data
}

export async function exportInvoices(): Promise<Blob> {
  const res = await api.get('/export', { responseType: 'blob' })
  return res.data
}

export default api
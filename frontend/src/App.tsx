import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Invoices from './pages/Invoices'
import InvoiceDetail from './pages/InvoiceDetail'
import Upload from './pages/Upload'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout><Navigate to="/invoices" replace /></Layout>} />
        <Route path="/invoices" element={<Layout><Invoices /></Layout>} />
        <Route path="/invoices/:id" element={<Layout><InvoiceDetail /></Layout>} />
        <Route path="/upload" element={<Layout><Upload /></Layout>} />
        <Route path="*" element={<Navigate to="/invoices" replace />} />
      </Routes>
    </BrowserRouter>
  )
}
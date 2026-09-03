# Invoice Extract — Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Cloudflare Edge                          │
│                invoiceextract.arjism.com (HTTPS)                │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Cloudflare Tunnel (cf-tunnel)                │
│              http://192.168.88.101:8105 (plain HTTP)            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Nginx Reverse Proxy                      │
│                    :8105 → :8000 (backend)                      │
│                    :8106 → :80 (dashboard)                      │
└────────────────────────────┬────────────────────────────────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                              ▼
┌──────────────────────┐        ┌──────────────────────┐
│   Python + FastAPI   │        │  React Dashboard     │
│   :8000 (internal)   │        │  :80 (internal)      │
│                      │        │                      │
│  - OCR Processing    │        │  - Tailwind CSS      │
│  - AI Extraction     │        │  - Upload Interface  │
│  - PDF Parsing       │        │  - Data Viewer       │
│  - OpenAI Vision     │        │  - Export            │
└──────────┬───────────┘        └──────────────────────┘
           │
           ▼
┌──────────────────────┐
│   Local Storage      │
│   /app/data/         │
└──────────────────────┘
```

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Language | Python | 3.11+ |
| Web Framework | FastAPI | 0.115+ |
| AI Integration | OpenAI Vision / 9Router | - |
| OCR | Tesseract / AI Vision | - |
| PDF Processing | PyPDF2 / pdfplumber | - |
| Frontend | React + Vite + TypeScript | Vite 5, React 18 |
| Styling | Tailwind CSS | v3 |
| Deployment | Docker Compose | v3.8 |
| Reverse Proxy | Nginx | - |
| Tunnel | Cloudflare Tunnel | - |

## Key Design Decisions

### 1. AI-Powered Extraction
- OpenAI Vision API for accurate field extraction
- Handles structured and semi-structured invoices
- Confidence scoring for extracted fields

### 2. Multi-Format Support
- PDF, PNG, JPG, TIFF support
- OCR fallback for scanned documents
- Batch processing capability

### 3. Structured Output
- JSON extraction with standardized fields
- Line items, totals, vendor info
- Export to CSV, Excel

### 4. Upload Pipeline
- Upload → Validate → Extract → Review → Export
- Progress tracking per file
- Error handling for failed extractions

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check |

### Authenticated
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/extract` | Extract invoice data |
| GET | `/api/invoices` | List invoices |
| GET | `/api/invoices/:id` | Get invoice details |
| DELETE | `/api/invoices/:id` | Delete invoice |
| GET | `/api/export` | Export data |

## Ports

| Service | External | Internal |
|---------|----------|----------|
| Backend | `:8105` | `:8000` |
| Dashboard | `:8106` | `:80` |

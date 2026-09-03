# Invoice Extract — AI Invoice Data Extraction

An AI-powered invoice data extraction tool that pulls structured data from PDFs and images.

## Quick Start

```bash
# Clone
git clone https://github.com/lovelymondayz/invoice-extract.git
cd invoice-extract

# Start all services
docker compose up -d --build

# Dashboard: http://localhost:8106
# API: http://localhost:8105
```

## Features

- **AI Extraction**: Extract vendor, line items, totals from invoices
- **Multi-Format**: PDF, PNG, JPG, TIFF support
- **OCR Fallback**: AI Vision + OCR for scanned documents
- **Batch Processing**: Process multiple invoices at once
- **Export**: CSV, Excel, JSON output
- **Confidence Scoring**: Review low-confidence extractions

## API Endpoints

### Public
- `GET /api/health` — Health check

### Authenticated
- `POST /api/extract` — Extract invoice data
- `GET /api/invoices` — List invoices
- `GET /api/invoices/:id` — Get invoice details
- `DELETE /api/invoices/:id` — Delete invoice
- `GET /api/export` — Export data

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| OPENAI_API_KEY | - | OpenAI API key |
| DATABASE_PATH | /app/data | Data directory |
| UPLOAD_DIR | - | Upload directory |
| OUTPUT_DIR | - | Output directory |
| MAX_FILE_SIZE_MB | 10 | Max file size (MB) |

## Development

```bash
# Backend only
cd backend
pip install -r requirements.txt
uvicorn src.api:app --reload

# Frontend only
cd frontend
npm install
npm run dev
```

## Deployment

1. Push to `main` → GitHub Action auto-deploys
2. Or manually: `ssh vps && cd /root/invoice-extract && ./update.sh`

## License

MIT

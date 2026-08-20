
from fastapi import FastAPI, HTTPException, UploadFile, File
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse
from pydantic import BaseModel
from typing import Optional
import sqlite3, os, uuid, csv, json
from datetime import datetime
from pathlib import Path

app = FastAPI(title="Invoice Extract API", version="1.0.0")
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])

DATABASE_PATH = os.getenv("DATABASE_PATH", "/app/data/invoices.db")
OUTPUT_DIR = os.getenv("OUTPUT_DIR", "/app/output/exports")
UPLOAD_DIR = "/app/data/uploads"
Path(DATABASE_PATH).parent.mkdir(parents=True, exist_ok=True)
Path(OUTPUT_DIR).mkdir(parents=True, exist_ok=True)
Path(UPLOAD_DIR).mkdir(parents=True, exist_ok=True)

def init_db():
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()
    c.execute("CREATE TABLE IF NOT EXISTS invoices (id TEXT PRIMARY KEY, filename TEXT, vendor_name TEXT, vendor_address TEXT, invoice_number TEXT, invoice_date TEXT, due_date TEXT, subtotal TEXT, tax TEXT, total TEXT, raw_text TEXT, created_at TEXT)")
    c.execute("CREATE TABLE IF NOT EXISTS line_items (id TEXT PRIMARY KEY, invoice_id TEXT, description TEXT, quantity TEXT, unit_price TEXT, total TEXT)")
    conn.commit(); conn.close()

def get_db():
    conn = sqlite3.connect(DATABASE_PATH); conn.row_factory = sqlite3.Row; return conn

@app.get("/health")
async def health(): return {"status": "healthy"}

@app.get("/")
async def root(): return {"service": "Invoice Extract API", "version": "1.0.0"}

@app.post("/extract")
async def extract(file: UploadFile = File(...)):
    iid = str(uuid.uuid4())
    filepath = Path(UPLOAD_DIR) / f"{iid}_{file.filename}"
    with open(filepath, "wb") as f:
        f.write(await file.read())
    # Basic extraction simulation
    conn = get_db()
    c = conn.cursor()
    c.execute("INSERT INTO invoices (id,filename,vendor_name,invoice_number,invoice_date,total,raw_text,created_at) VALUES (?,?,?,?,?,?,?,?)",
        (iid, file.filename, "Unknown Vendor", "N/A", "N/A", "N/A", f"File: {file.filename}", datetime.utcnow().isoformat()))
    conn.commit(); conn.close()
    return {"invoice_id": iid, "message": "Invoice uploaded, extraction pending"}

@app.get("/invoices")
async def list_invoices():
    conn = get_db()
    invoices = [dict(r) for r in conn.execute("SELECT * FROM invoices ORDER BY created_at DESC").fetchall()]
    conn.close(); return {"invoices": invoices}

@app.get("/invoices/{iid}")
async def get_invoice(iid: str):
    conn = get_db()
    c = conn.cursor()
    c.execute("SELECT * FROM invoices WHERE id=?", (iid,))
    inv = c.fetchone()
    if not inv: conn.close(); raise HTTPException(404, "Invoice not found")
    c.execute("SELECT * FROM line_items WHERE invoice_id=?", (iid,))
    items = [dict(r) for r in c.fetchall()]
    conn.close()
    result = dict(inv); result["line_items"] = items
    return result

@app.get("/export")
async def export_all():
    conn = get_db()
    invoices = [dict(r) for r in conn.execute("SELECT * FROM invoices ORDER BY created_at DESC").fetchall()]
    conn.close()
    filename = f"invoices_export.csv"
    filepath = Path(OUTPUT_DIR) / filename
    with open(filepath, "w", newline="") as f:
        w = csv.DictWriter(f, fieldnames=["id","filename","vendor_name","invoice_number","invoice_date","due_date","total"], extrasaction="ignore")
        w.writeheader(); w.writerows(invoices)
    return FileResponse(filepath, filename=filename)

@app.on_event("startup")
async def startup(): init_db()

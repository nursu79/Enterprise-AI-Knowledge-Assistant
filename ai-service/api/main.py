from fastapi import FastAPI, UploadFile, File, HTTPException
from typing import List

from .models import QueryRequest, QueryResponse, IngestResponse

app = FastAPI(
    title="Enterprise AI Knowledge Assistant",
    description="AI Service for RAG Pipeline",
    version="1.0.0"
)

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "ai-service"}

@app.post("/ingest", response_model=IngestResponse)
async def ingest_document(file: UploadFile = File(...)):
    # Placeholder for document ingestion logic
    if not file.filename.endswith(('.pdf', '.txt', '.csv')):
        raise HTTPException(status_code=400, detail="Unsupported file format")
    
    return IngestResponse(
        status="success",
        message=f"Document {file.filename} ingested successfully",
        document_id="dummy-id"
    )

@app.post("/query", response_model=QueryResponse)
async def query_assistant(request: QueryRequest):
    # Placeholder for RAG pipeline logic
    return QueryResponse(
        answer="This is a placeholder answer.",
        context=["Placeholder context chunk 1", "Placeholder context chunk 2"]
    )

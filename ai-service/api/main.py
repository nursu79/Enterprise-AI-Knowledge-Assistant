from fastapi import FastAPI, UploadFile, File, HTTPException
from typing import List
import uuid

from .models import QueryRequest, QueryResponse, IngestResponse
from rag.document_processor import DocumentProcessor
from embeddings.generator import EmbeddingGenerator

processor = DocumentProcessor(chunk_size=1000, chunk_overlap=200)
embedder = EmbeddingGenerator()

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
    if not file.filename.endswith(('.pdf', '.txt', '.csv')):
        raise HTTPException(status_code=400, detail="Unsupported file format")
    
    content = await file.read()
    try:
        chunks = processor.process_file(content, file.filename)
        embeddings = embedder.generate_embeddings(chunks)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error processing document: {str(e)}")
        
    doc_id = str(uuid.uuid4())
    # TODO: store chunks and embeddings into FAISS
    
    return IngestResponse(
        status="success",
        message=f"Document {file.filename} ingested, split into {len(chunks)} chunks, and {len(embeddings)} embeddings generated",
        document_id=doc_id
    )

@app.post("/query", response_model=QueryResponse)
async def query_assistant(request: QueryRequest):
    # Placeholder for RAG pipeline logic
    return QueryResponse(
        answer="This is a placeholder answer.",
        context=["Placeholder context chunk 1", "Placeholder context chunk 2"]
    )

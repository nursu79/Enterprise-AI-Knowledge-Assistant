from fastapi import FastAPI, UploadFile, File, HTTPException
from typing import List
import uuid

from .models import QueryRequest, QueryResponse, IngestResponse
from rag.document_processor import DocumentProcessor
from embeddings.generator import EmbeddingGenerator
from retrieval.faiss_store import FAISSStore

processor = DocumentProcessor(chunk_size=1000, chunk_overlap=200)
embedder = EmbeddingGenerator()
vector_store = FAISSStore()

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
        vector_store.add_embeddings(embeddings, chunks)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error processing document: {str(e)}")
        
    doc_id = str(uuid.uuid4())

    
    return IngestResponse(
        status="success",
        message=f"Document {file.filename} ingested, split into {len(chunks)} chunks, and {len(embeddings)} embeddings generated",
        document_id=doc_id
    )

@app.post("/query", response_model=QueryResponse)
async def query_assistant(request: QueryRequest):
    try:
        query_embedding = embedder.generate_embeddings([request.query])[0]
        results = vector_store.search(query_embedding, request.top_k)
        
        context = [chunk for chunk, _score in results]
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error during retrieval: {str(e)}")
        
    # Placeholder for LLM generation
    return QueryResponse(
        answer="This is a placeholder answer. Ensure RAG LLM pipeline processes context.",
        context=context
    )

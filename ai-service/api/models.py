from pydantic import BaseModel
from typing import List, Optional

class QueryRequest(BaseModel):
    query: str
    top_k: int = 5

class QueryResponse(BaseModel):
    answer: str
    context: List[str]

class IngestResponse(BaseModel):
    status: str
    message: str
    document_id: str

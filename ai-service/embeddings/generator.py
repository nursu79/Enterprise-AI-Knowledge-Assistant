from sentence_transformers import SentenceTransformer
from typing import List

class EmbeddingGenerator:
    def __init__(self, model_name: str = "all-MiniLM-L6-v2"):
        # Load a lightweight model suitable for semantic search
        self.model = SentenceTransformer(model_name)
    
    def generate_embeddings(self, chunks: List[str]) -> List[List[float]]:
        """
        Generate embeddings for document chunks.
        sentence-transformers natively handles batching when passed a list of strings.
        """
        if not chunks:
            return []
            
        embeddings = self.model.encode(chunks, batch_size=32, show_progress_bar=False)
        return embeddings.tolist()

import faiss
import numpy as np
from typing import List, Tuple, Dict

class FAISSStore:
    def __init__(self, dimension: int = 384):
        """
        Initialize FAISS index. dimension=384 for all-MiniLM-L6-v2.
        Using IndexFlatL2 for exact, fast similarity search for typical RAG workloads.
        """
        self.dimension = dimension
        self.index = faiss.IndexFlatL2(self.dimension)
        self.chunk_metadata: Dict[int, str] = {}
        self.current_id = 0
        
    def add_embeddings(self, embeddings: List[List[float]], chunks: List[str]):
        """Store chunk vectors in the FAISS index and save chunk text in memory."""
        if not embeddings or not chunks:
            return
            
        vectors = np.array(embeddings).astype('float32')
        self.index.add(vectors)
        
        for chunk in chunks:
            self.chunk_metadata[self.current_id] = chunk
            self.current_id += 1
            
    def search(self, query_embedding: List[float], top_k: int = 5) -> List[Tuple[str, float]]:
        """Search for top_k most similar chunks."""
        if self.index.ntotal == 0:
            return []
            
        vector = np.array([query_embedding]).astype('float32')
        distances, indices = self.index.search(vector, top_k)
        
        results = []
        for dist, idx in zip(distances[0], indices[0]):
            if idx != -1 and idx in self.chunk_metadata:
                results.append((self.chunk_metadata[idx], float(dist)))
                
        return results

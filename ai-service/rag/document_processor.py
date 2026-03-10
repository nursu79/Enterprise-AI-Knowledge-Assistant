import io
import csv
from pypdf import PdfReader
from typing import List

class DocumentProcessor:
    def __init__(self, chunk_size: int = 1000, chunk_overlap: int = 200):
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap

    def process_file(self, file_content: bytes, filename: str) -> List[str]:
        text = self._extract_text(file_content, filename)
        return self._chunk_text(text)

    def _extract_text(self, file_content: bytes, filename: str) -> str:
        if filename.endswith(".pdf"):
            reader = PdfReader(io.BytesIO(file_content))
            return "\n".join([page.extract_text() for page in reader.pages if page.extract_text()])
        elif filename.endswith(".csv"):
            decoded_content = file_content.decode("utf-8")
            reader = csv.reader(io.StringIO(decoded_content))
            return "\n".join([", ".join(row) for row in reader])
        elif filename.endswith(".txt"):
            return file_content.decode("utf-8")
        else:
            raise ValueError(f"Unsupported file type: {filename}")

    def _chunk_text(self, text: str) -> List[str]:
        chunks = []
        start = 0
        text_len = len(text)
        
        if text_len == 0:
            return chunks
            
        while start < text_len:
            end = min(start + self.chunk_size, text_len)
            
            # Try to snap to the nearest space character to avoid cutting words
            if end < text_len:
                last_space = text.rfind(' ', start, end)
                if last_space != -1 and last_space > start:
                    end = last_space
                    
            chunks.append(text[start:end].strip())
            
            if end == text_len:
                break
            start = end - self.chunk_overlap
            
            # Ensure we always make progress to combat edge cases
            if start <= 0 or start >= end:
                start = end
                
        return [c for c in chunks if c]

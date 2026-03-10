from transformers import pipeline

class LLMGenerator:
    def __init__(self, model_name: str = "google/flan-t5-small"):
        """
        Initialize the HuggingFace LLM generating pipeline.
        We use a lightweight model suitable for edge or quick inference.
        """
        self.generator = pipeline("text2text-generation", model=model_name)
    
    def generate_answer(self, query: str, context: list[str]) -> str:
        """
        Construct a grounded prompt with retrieved chunks and generate an answer.
        """
        if not context:
            return "I don't have enough context to answer that."
            
        context_text = "\n\n".join(context)
        prompt = (
            "Use the following context to answer the question. "
            "If the answer is not in the context, say 'I cannot answer based on the context.'\n\n"
            f"Context:\n{context_text}\n\n"
            f"Question: {query}\n\n"
            "Answer:"
        )
        
        # Generation with reasonable defaults
        result = self.generator(
            prompt, 
            max_new_tokens=150, 
            do_sample=False,  # deterministic generation for facts
            truncation=True
        )
        return result[0]['generated_text'].strip()

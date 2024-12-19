from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List

app = FastAPI()

class QueryRequest(BaseModel):
    query: str


@app.post("/embed/query")
def get_embedding(request: QueryRequest):
    query_text = request.query

    if not query_text:
        raise HTTPException(status_code=400, detail="Query is empty.")

    # query를 벡터로 변환
    vector = embeddings_model.embed_query(query_text)

    # 결과 반환 (JSON)
    return {"vector": vector}

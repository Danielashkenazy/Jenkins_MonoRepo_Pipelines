from fastapi import FastAPI, HTTPException

app = FastAPI()

def compute_total(transactions: list[float]) -> float:
    return sum(x for x in transactions if x > 0)

@app.post("/transactions/total")
def total(data: dict):
    if "transactions" not in data or not isinstance(data["transactions"], list):
        raise HTTPException(status_code=400, detail="tttttttttttttransactions must be a list")

    result = compute_total(data["transactions"])
    return {"total": result}

@app.get("/health")
def health():
    return {"status": "ok", "service": "transaction-service"}

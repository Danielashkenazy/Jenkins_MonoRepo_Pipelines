from fastapi import FastAPI

app = FastAPI()

@app.get("/health")
def health():
    return {"status": "ok", "service": "transaction-service"}

@app.get("/")
def main():
    return {"message": "Welcome to the Transaction Serviceeeeeeee"}


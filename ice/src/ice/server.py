import argparse
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI

from ice.classifier import Classifier

_classifier: Classifier | None = None


def _get_classifier() -> Classifier:
    global _classifier
    if _classifier is None:
        _classifier = Classifier()
    return _classifier


@asynccontextmanager
async def lifespan(app: FastAPI):
    _get_classifier()
    yield


app = FastAPI(title="ice", docs_url=None, redoc_url=None, lifespan=lifespan)


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/classify")
def classify(payload: dict):
    text = payload.get("text", "")
    if not text:
        return {"error": "missing 'text' field"}
    return _get_classifier().classify(text)


def main():
    parser = argparse.ArgumentParser(description="ice — prompt injection classifier")
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=9119)
    args = parser.parse_args()
    uvicorn.run(app, host=args.host, port=args.port, log_level="warning")


if __name__ == "__main__":
    main()

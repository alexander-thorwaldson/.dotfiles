from pathlib import Path

import numpy as np
from transformers import AutoTokenizer
import onnxruntime as ort

MODEL_ID = "protectai/deberta-v3-base-prompt-injection-v2"
CACHE_DIR = Path.home() / ".cache" / "ice"
LABELS = ["SAFE", "INJECTION"]


class Classifier:
    def __init__(self) -> None:
        self._tokenizer = AutoTokenizer.from_pretrained(MODEL_ID, cache_dir=CACHE_DIR)
        onnx_path = self._ensure_onnx()
        self._session = ort.InferenceSession(
            str(onnx_path), providers=["CPUExecutionProvider"]
        )

    def _ensure_onnx(self) -> Path:
        """Download and convert model to ONNX if not already cached."""
        onnx_dir = CACHE_DIR / "onnx"
        onnx_path = onnx_dir / "model.onnx"
        if onnx_path.exists():
            return onnx_path

        from optimum.onnxruntime import ORTModelForSequenceClassification

        model = ORTModelForSequenceClassification.from_pretrained(
            MODEL_ID, export=True, cache_dir=CACHE_DIR
        )
        onnx_dir.mkdir(parents=True, exist_ok=True)
        model.save_pretrained(onnx_dir)
        return onnx_path

    def classify(self, text: str) -> dict:
        tokens = self._tokenizer(
            text, return_tensors="np", truncation=True, max_length=512
        )
        inputs = {k: v for k, v in tokens.items() if k in ("input_ids", "attention_mask")}
        (logits,) = self._session.run(None, inputs)
        probs = _softmax(logits[0])
        label_idx = int(np.argmax(probs))
        return {
            "label": LABELS[label_idx],
            "score": round(float(probs[label_idx]), 4),
            "scores": {LABELS[i]: round(float(probs[i]), 4) for i in range(len(LABELS))},
        }


def _softmax(x: np.ndarray) -> np.ndarray:
    e = np.exp(x - np.max(x))
    return e / e.sum()

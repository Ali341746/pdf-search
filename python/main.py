# pdf_service_path.py
from fastapi import FastAPI
from pydantic import BaseModel
import PyPDF2
from fastapi.responses import JSONResponse

app = FastAPI(title="PDF Text Extraction Service")

class PDFPath(BaseModel):
    path: str

@app.post("/extract")
async def extract_text(pdf: PDFPath):
    try:
        with open(pdf.path, "rb") as f:
            reader = PyPDF2.PdfReader(f)
            text = ""
            for page in reader.pages:
                page_text = page.extract_text()
                if page_text:
                    text += page_text + "\n"
    except FileNotFoundError:
        return JSONResponse(status_code=404, content={"error": "File not found"})
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})
    
    return {"filename": pdf.path, "text": text}

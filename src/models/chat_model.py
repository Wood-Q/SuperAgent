import os

from langchain_openai import ChatOpenAI
from pydantic import SecretStr

chat_model=ChatOpenAI(
    model="doubao-1-5-pro-32k-250115",
    base_url="https://ark-cn-beijing.bytedance.net/api/v3",
    api_key=SecretStr(os.getenv("DOUBAO_API_KEY")),
    temperature=0
)
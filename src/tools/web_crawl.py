
import os
from firecrawl import FirecrawlApp
from langchain.tools import tool

@tool
def web_crawl(url:str)->str:
    """
    Crawl a website and return the markdown content.

    Args:
        url: The URL of the website to crawl.

    Returns:
        The markdown content of the website.
    """
    firecrawl = FirecrawlApp(api_key=os.getenv("FIRECRAWL_API_KEY"))
    response = firecrawl.scrape_url(
        url=url,
        formats=["markdown"],
        only_main_content=True
    )
    return response.markdown

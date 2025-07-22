import uuid

from langgraph.graph.state import CompiledStateGraph

def run_agent(agent: CompiledStateGraph, message: str):
    result = agent.stream(
        {"messages": [{"role": "user", "content": message}]},
        stream_mode="values",
        config={"thread_id": uuid.uuid4()},
    )
    for chunk in result:
        messages = chunk["messages"]
        last_message = messages[-1]
        last_message.pretty_print()
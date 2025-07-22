from rich.traceback import install

# 运行这一行，rich 就会自动接管后续所有的报错信息
install()

from dotenv import load_dotenv

load_dotenv()

from src.agents.agent_runner import run_agent
from src.agents.researcher import research


def main():
    run_agent(agent=research, message="明日方舟缪尔赛思的种族是什么")
    # 获取对应的 Mermaid 代码

    research.get_graph().draw_mermaid()

    # 保存为 PNG 文件
    research.get_graph().draw_mermaid_png(output_file_path="react_agent.png")


if __name__ == "__main__":
    main()

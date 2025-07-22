# -*- coding: utf-8 -*-
from langchain.schema import SystemMessage, HumanMessage

import json
from src.models.chat_model import chat_model

res = chat_model.invoke(
    [
        SystemMessage(
            content="""Provide answers to the user's questions based on the context they provide.

# Steps

1. Carefully read and understand the context provided by the user.
2. Analyze the question in relation to the given context.
3. Use logical reasoning to derive the answer based on the context.
4. If the context is insufficient to answer the question, politely ask for clarification or additional details.

# Output Format

- Provide a clear and concise answer in complete sentences.
- If reasoning is required, include the reasoning process before presenting the conclusion.

# Notes

- Ensure the response is directly relevant to the user's question and context.
- Avoid making assumptions beyond the provided context unless explicitly requested."""
        ),
        HumanMessage(
            content=json.dumps(
                [
                    {
                        "title": "欧国联决赛:C罗破门，带队夺冠，葡萄牙点球6-5击败西班牙!|c罗|门迪|尼科|葡萄牙队|豪尔赫·门德斯_网易订阅",
                        "url": "https://www.163.com/dy/article/K1JCRTCN05567SBP.html",
                        "content": "北京时间6月9日，欧国联决赛，葡萄牙2-2战平西班牙，点球大战总比分6-5取胜夺冠!这场比赛踢得非常精彩，西班牙的传控确实厉害，他们整体的技术和脚法是世界顶级的，西班牙队无论面对任何对手，都能在控球率方面占据优势，而葡萄牙同样很强，在推进到对方的禁区前沿时，葡萄牙有很多种",
                    },
                    {
                        "title": "欧国联决赛：葡萄牙5-3点球胜西班牙，C罗闪耀，金球奖归属引关注 - 知乎",
                        "url": "https://zhuanlan.zhihu.com/p/1915618913749795333",
                        "content": "欧国联 最后一轮踢完了。 葡萄牙在决赛5-3点球战胜西班牙，第二次拿冠军。法国打德国2-0赢了，拿了第三。c罗这场比赛进了球，成了国家队历史第一射手， 金球奖 基本没悬念了。 法国和德国那场球， 姆巴佩 表现很猛。 上半场他带球进禁区自己射门得分，下半场又给队友传球让对方空门进球。",
                    },
                    {
                        "title": "欧国联决赛葡萄牙点球5-3胜西班牙，C罗门德斯建功",
                        "url": "https://sports.sina.com.cn/g/2025-06-09/doc-inezmvxn4600761.shtml",
                        "content": "🆚欧国联a级决赛|葡萄牙2-2西班牙（点球5-3）🎯努诺·门德斯传射建功，巴黎四人随队夺冠🎞第21分钟，苏比门迪补射破门，西班牙1-0葡萄牙。🎞第",
                    },
                    {
                        "title": "欧国联：葡萄牙点球大战7-5西班牙第二次夺冠 C罗破门+伤退",
                        "url": "https://news.qq.com/rain/a/20250609A01I7S00",
                        "content": "欧国联-葡萄牙2-2点球5-3西班牙第2次欧国联夺冠!c罗破门扳平+伤退，莫拉塔失点. 北京时间6月9日凌晨，欧国联迎来最终的决赛，葡萄牙大战西班牙。",
                    },
                    {
                        "title": "40岁C罗收获生涯第36冠!葡萄牙点球7-5战胜西班牙，欧国联夺冠_澎湃号·媒体_澎湃新闻-The Paper",
                        "url": "https://www.thepaper.cn/newsDetail_forward_30952562",
                        "content": "40岁C罗收获生涯第36冠！葡萄牙点球7-5战胜西班牙，欧国联夺冠_澎湃号·媒体_澎湃新闻-The Paper Image 1: 澎湃Logo *   _要闻_ *   _深度_ *   _直播_ *   _视频_ *   _更多_ **下载客户端** 登录 无障碍) 40岁C罗收获生涯第36冠！葡萄牙点球7-5战胜西班牙，欧国联夺冠 Image 2 上游新闻 关注 重庆 来源：澎湃新闻·澎湃号·媒体 字号 北京时间6月9日，在欧国联决赛中，葡萄牙点球7-5战胜西班牙，队史第二次夺得欧国联冠军！40岁C罗连场破门，门德斯建功，祖比门迪、奥亚萨瓦尔破门。 Image 3Image 4 葡萄牙国家队官方发文庆祝夺冠： 我们又做到了！第二次捧起欧国联的奖杯，属于葡萄牙的荣耀，再次刻进这段绿茵史。不是奇迹，是坚持；不是偶然，是一代又一代的传承与信念。我们，是欧国联历史上夺冠最多的那支队伍！ Image 5 本场比赛是C罗国家队生涯至今参加的第4场决赛，此前分别是2004年欧洲杯决赛、2016年欧洲杯决赛（因伤离场）、2019年欧国联决赛，仅2004年欧洲杯未能夺冠。至此，40岁C罗职业生涯已夺36冠：国家队3冠，俱乐部33冠。 Image 6 赛后在接受采访时，C罗谈到了自己的伤势，他表示，热身时自己就感觉有些不适。 C罗这样谈道：“我在热身时就已经感觉到了，已经有一段时间了。但为了国家队，即使要断腿我也愿意。这是一个冠军，我必须上场，我尽了全力，坚持到了最后一刻，还打进了一球。” “我非常高兴。首先是为了这一代球员，他们值得一个冠军。为了我们的家人，我的家人都在这里……为葡萄牙赢得胜利是特别的。我有很多冠军头衔，但没有什么比为葡萄牙赢得胜利更美好的了。泪水和完成使命的感觉……这是美妙的。这是我们的国家。我们是一个小国，但有着远大的抱负。” “我曾在很多国家和俱乐部踢过球，但当人们提到葡萄牙时，那种感觉是特别的。作为这支球队的队长，我感到非常自豪，赢得冠军总是国家队的最高荣誉。未来我会考虑短期目标，我受了伤，而且伤势加重了……但我还是坚持了下来，因为为了国家队你必须全力以赴。” 关于主教练马丁内斯，C罗表示：“我为他感到非常高兴，他是一个西班牙人，但他为我们的国家付出了最大的努力。我们已经赢得了这个冠军，但这只是我们的动力，是我们渴望更多的开始。” Image 7 知名体育评论员詹俊赛后点评道： 逆转之夜！葡萄牙队两次落后两次扳平比分，最终在点球决战力克西班牙再次夺得欧国联冠军。C罗vs亚马尔，2003年18岁的C罗迎来在国家队的首秀，二十二年后40岁的C罗仍然能打进关键入球拿到个人第三项国际大赛的冠军。而亚马尔略显平淡因为碰到克星——努诺-门德斯，这位大巴黎的左后卫是个人心目中本场的MVP，“世界第一左后卫”的封号当之无愧。西班牙队的右路防守是软肋，明年世界杯他们还是需要卡瓦哈尔。 Image 8 本文为澎湃号作者或机构在澎湃新闻上传并发布，仅代表该作者或机构观点，不代表澎湃新闻的观点或立场，澎湃新闻仅提供信息发布平台。申请澎湃号请用电脑访问http://renzheng.thepaper.cn。 *   Image 9### C罗进球封王！这支葡萄牙让人想起了2022年的阿根廷  *   Image 10### C罗语出惊人：金球奖失去公信力  *   Image 11### 体坛联播｜C罗进球助葡萄牙夺冠，阿尔拉卡斯卫冕法网  *   Image 12_01:00_### 王大雷国脚生涯谢幕倒计时：希望给年轻球员做好榜样  *   Image 13### C罗进球封王！这支葡萄牙让人想起了2022年的阿根廷  *   Image 14### C罗语出惊人：金球奖失去公信力  *   Image 15### 体坛联播｜C罗进球助葡萄牙夺冠，阿尔拉卡斯卫冕法网  *   Image 16_01:00_### 王大雷国脚生涯谢幕倒计时：希望给年轻球员做好榜样  *   Image 17### C罗进球封王！这支葡萄牙让人想起了2022年的阿根廷  *   Image 18### C罗语出惊人：金球奖失去公信力  *   Image 19### 体坛联播｜C罗进球助葡萄牙夺冠，阿尔拉卡斯卫冕法网  *   Image 20_01:00_### 王大雷国脚生涯谢幕倒计时：希望给年轻球员做好榜样  Image 21 Image 22 Image 23 Image 24 Image 25 Image 26 *   报料邮箱: news@thepaper.cn Image 28 Image 29",
                    },
                ]
            )
        ),
        HumanMessage(content="欧足联决赛有什么看点？"),
    ]
)
print(res.content)
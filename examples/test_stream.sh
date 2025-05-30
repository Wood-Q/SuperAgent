#!/bin/bash

# 测试流式聊天API
echo "测试流式聊天API..."

curl -N -X POST http://localhost:8888/api/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "userInput": "你好，请介绍一下你自己"
  }'

echo -e "\n\n测试完成" 
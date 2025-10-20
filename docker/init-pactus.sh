#!/bin/bash
set -e

PACTUS_IMAGE=pactus/pactus:latest
PACTUS_WALLET_PASSWORD=${PACTUS_WALLET_PASSWORD:-onepactus}

# 获取项目名称（docker-compose 会用目录名作为前缀）
PROJECT_NAME=$(basename "$(pwd)")
VOLUME_NAME="${PROJECT_NAME}_pactus-data"

echo "检查 Pactus 钱包初始化状态..."
echo "使用 Docker 卷: $VOLUME_NAME"

# 检查卷是否存在，如果不存在则创建
if ! docker volume inspect "$VOLUME_NAME" >/dev/null 2>&1; then
    echo "创建 Docker 卷: $VOLUME_NAME"
    docker volume create "$VOLUME_NAME"
fi

# 检查卷是否已经初始化（检查是否存在 config.toml 文件）
if docker run --rm -v "$VOLUME_NAME:/root/pactus" alpine test -f /root/pactus/config.toml 2>/dev/null; then
    echo "✓ Pactus 节点已经初始化，跳过初始化步骤。"
    exit 0
fi

echo "开始初始化 Pactus 节点..."
echo "钱包密码: ${PACTUS_WALLET_PASSWORD:0:3}***"

# 运行初始化命令，将数据直接写入 Docker 卷
yes | docker run -i --rm \
    -v "$VOLUME_NAME:/root/pactus" \
    "$PACTUS_IMAGE" \
    pactus-daemon init \
    --password "$PACTUS_WALLET_PASSWORD" \
    --val-num 1

echo ""
echo "✓ Pactus 节点初始化完成！"
echo "钱包数据已保存到卷: $VOLUME_NAME"
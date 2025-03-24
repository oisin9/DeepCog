# DeepCog

## 项目简介

DeepCog是一个开源的AI API服务，将Deepseek和Claude、Gemini、豆包、Qwen等模型结合，通过把Deepseek的思考过程输出给其他模型，达到1+1>2的效果。

## 特性

*   **OpenAI API兼容**：支持标准的OpenAI API接口，方便与现有工具和库集成。
*   **DeepSeek 整合**: 无缝集成 DeepSeek 模型，获取其强大的推理能力。
*   **多模型支持**: 支持多个兼容 OpenAI 接口的模型，可以将 DeepSeek 的思考过程传递给其他模型。
*   **中间过程输出**: 获取 DeepSeek 模型的中间思考过程，用于增强其他模型的推理结果。
*   **流式响应**: 支持 SSE (Server-Sent Events) 流式响应。
*   **基于Go和Echo框架**: 高性能,易于开发和维护。
*   **易于部署**：使用go编写，开箱即用无复杂依赖。

## 快速开始

### 安装

需要提前配置golang的环境，如果没有可以从[Golang官网](https://go.dev/dl/)下载安装。

Linux执行下面的命令安装：

```shell
git clone https://github.com/oisin9/DeepCog
cd DeepCog
make
make install
```

### 配置

修改`/etc/deepcog/config.toml`配置文件。

下面是示例配置文件

```toml
[server]
port = "8080" # 监听端口

# 基础模型，这里的模型不提供对外服务，仅供models来使用
[[base_models]]
id="THINKING-MODEL" # 这里的id可以自己定义，在后面的models里使用
model_name="deepseek-r1" # 需要调用api平台的model name
base_url="https://api.deepseek.com" # 调用api平台的base url
api_key="" # 调用api平台的api_key

[[base_models]]
id="CONTENT_MODEL"
model_name="claude-3-7-sonnet-latest"
base_url=""
api_key=""

[[models]]
id="" # 对外服务的模型名
owned_by="" # 对外服务显示的服务提供商，可不填
thinking_model="THINKING-MODEL" # base_models里的id，如果不填则不使用思考模型的思考过程
generate_model="CONTENT_MODEL" # base_models里的id
api_key="" # 对外提供服务的api_key，如果留空则不校验
```

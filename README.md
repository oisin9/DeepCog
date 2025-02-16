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


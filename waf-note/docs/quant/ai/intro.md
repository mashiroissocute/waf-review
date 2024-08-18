## 背景

FreqAI旨在提供一个通用且可扩展的开源框架，用于实时部署用于市场预测的自适应模型

## 示例

运行机器学习模型：

freqtrade trade --config config_examples/config_freqai.example.json --strategy FreqaiExampleStrategy --freqaimodel LightGBMRegressor --strategy-path freqtrade/templates



运行深度学习模型：

freqtrade trade --config config_examples/config_freqai.example.json --strategy FreqaiExampleStrategy --freqaimodel PyTorchMLPRegressor --strategy-path freqtrade/templates
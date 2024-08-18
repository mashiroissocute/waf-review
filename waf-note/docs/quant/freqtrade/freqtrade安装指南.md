# freqtrade安装指南
- 开通境外云主机
- 下载项目
```
git clone https://github.com/freqtrade/freqtrade.git
```
- 下载安装miniconda 
```
wget https://mirrors.aliyun.com/anaconda/miniconda/Miniconda3-py311_23.5.2-0-Linux-x86_64.sh?spm=a2c6h.25603864.0.0.5c152eb1khtps2
```
- 创建环境 
```
conda create --name freqtrade python=3.11
```
- 安装Ta-lib 
```
cd build_helpers
bash install_ta-lib.sh ${CONDA_PREFIX} nosudo
```

有些freqtrade版本存在安装失败的问题，基于经验可以替换build_helpers.
使用旧版本的build_helpers.


- 安装python依赖
```
python3 -m pip install --upgrade pip
python3 -m pip install -r requirements.txt
python3 -m pip install -e .
```

有些freqtrade版本存在安装失败的问题，基于经验可以替换requirements.txt
使用旧版本的requirements.txt

- 开始使用
```
# Step 1 - Initialize user folder
freqtrade create-userdir --userdir user_data

# Step 2 - Create a new configuration file
freqtrade new-config --config user_data/config.json
```

- 安装UI
freqtrade install-ui


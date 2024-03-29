# FOFAPK
fofapk 用于对FOFA工程师的线索梳理能力进行PK，来验证能力的强弱。

针对同一个任务（比如企业的互联网暴露面梳理或者某个厂商的所有设备），两个参赛者提交FOFA的查询语句，看哪个参赛者提交的语句更好（全面且准确）。

## 实时开发
- 安装golang环境和wails（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`）。
- 安装nodejs

根目录下运行：`wails dev`

注意：确保环境变量配置了`FOFA_KEY`。

## 编译

根目录下运行：`wails build`.

## 效果
<img src="./docs/fofapk.jpg" width="800" />

## 可增加特性
- [x] 显示两边的IP总量
- [ ] 记录操作的历史，便于检查
- [ ] 聚类相同项目，比如fid/根域名之类的
- [ ] 扫码登陆，在各自的电脑上点击提交
- [ ] 统计具体的输赢情况数据

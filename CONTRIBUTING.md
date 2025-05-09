# 贡献指南

首先，感谢您考虑为 Bilibili Watcher 项目做出贡献！我们非常欢迎各种形式的贡献，无论是代码、文档、错误报告还是功能建议。

## 如何贡献

我们鼓励通过以下方式参与项目：

*   **报告 Bug**：如果您在使用过程中发现任何错误或问题，请通过 [GitHub Issues](https://github.com/krisxia0506/bilibili-watcher/issues) 提交详细的 Bug 报告。
*   **提出功能建议**：如果您对项目有新的想法或功能需求，也欢迎通过 [GitHub Issues](https://github.com/krisxia0506/bilibili-watcher/issues) 提出。
*   **提交代码 (Pull Requests)**：我们欢迎代码贡献，无论是修复 Bug 还是实现新功能。
*   **改进文档**：帮助我们改进项目的文档，使其更清晰、更易懂。

## 贡献流程

### 报告 Bug

在提交 Bug 报告之前，请先搜索现有的 Issues，确保您的问题尚未被报告。提交时，请尽可能提供以下信息：

*   清晰简洁的标题。
*   详细的复现步骤。
*   您期望的结果和实际发生的结果。
*   相关的环境信息（例如操作系统、浏览器版本、Go 版本、Node.js 版本等）。
*   相关的日志或截图（如果有）。

### 提交 Pull Request (PR)

1.  **Fork 本仓库**：点击仓库右上角的 "Fork" 按钮。
2.  **Clone 您的 Fork**：`git clone https://github.com/YOUR_USERNAME/bilibili-watcher.git`
3.  **创建新分支**：`git checkout -b feature/your-feature-name` 或 `fix/issue-number-description`。
4.  **进行修改**：按照项目的编码规范进行代码修改。
5.  **代码格式化与检查**：确保您的代码通过了项目配置的 linter 和 formatter（如果项目中有配置）。
6.  **提交更改**：编写清晰的 Commit Message (遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范)。
    ```bash
    git add .
    git commit -m "feat: 添加了某个很棒的功能"
    ```
7.  **推送到您的 Fork**：`git push origin feature/your-feature-name`
8.  **创建 Pull Request**：回到原始仓库的 GitHub 页面，点击 "New pull request" 按钮，选择您的分支与主仓库的 `main` (或目标) 分支进行比较，并提交 PR。
9.  **描述您的 PR**：在 PR 描述中清晰地说明您所做的更改、解决的问题以及任何相关的背景信息。

### 开发环境设置

请参考项目根目录下的 `README.md` 文件中关于项目设置和本地开发的部分。

## 编码规范

*   请遵循项目中已有的编码风格和规范。
*   对于 Go 代码，请遵循社区通用的 Go 编码规范，并确保代码通过 `gofmt` 或 `goimports` 格式化。
*   对于 TypeScript/JavaScript 代码，请遵循项目配置的 ESLint 和 Prettier 规范。
*   代码注释请使用中文。
*   日志输出内容请使用英文。
*   Commit Message 请遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

## 行为准则

我们期望所有贡献者都能遵守项目的 [行为准则 (CODE_OF_CONDUCT.md)](./CODE_OF_CONDUCT.md)。

感谢您的贡献！
